package filemanager

import (
	"sync"

	"github.com/xiusin/pine/di"

	"github.com/xiusin/pinecms/src/config"

	"github.com/xiusin/pine"
	"github.com/xiusin/pinecms/src/application/controllers/backend/filemanager/tables"
	"github.com/xiusin/pinecms/src/common/helper"
	"github.com/xiusin/pinecms/src/common/storage"
)

type ResResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var once sync.Once

const Logined = "true"
const DownloadFlag = "download"

const serviceFtpStorage = "pinecms.service.fm.ftp"

func InitInstall(app *pine.Application, urlPrefix, dir string) {
	once.Do(func() {
		app.Use(func(ctx *pine.Context) {
			if str, _ := ctx.Input().GetString("fmq"); str == DownloadFlag {
				ctx.Response.Header.Set("Content-Disposition", "attachment")
			}
			ctx.Next()
		})
		app.Static(urlPrefix, dir, 1)
		orm := helper.GetORM()

		defer func() {
			if err := recover(); err != nil {
				pine.Logger().Warn("初始化安装失败", err)
			}
		}()

		if err := orm.Sync2(&tables.FileManagerAccount{}); err != nil {
			pine.Logger().Warn(err.Error())
		}
		if count, _ := orm.Count(&tables.FileManagerAccount{}); count == 0 {
			user := &tables.FileManagerAccount{Username: "admin", Nickname: "Administer", Engine: "本地存储"}
			user.Init()
			user.Password = user.GetMd5Pwd("admin888")
			if _, err := orm.InsertOne(user); err != nil {
				pine.Logger().Warn("新增用户失败", err)
			}
		}

		di.Set(serviceFtpStorage, func(builder di.AbstractBuilder) (any, error) {
			cfg, _ := config.SiteConfig()
			cfg["PROXY_SITE_URL"] = "/filemanager/proxy_content?path={path}"
			ftp := storage.NewFtpUploader(cfg)
			return ftp, nil
		}, true)

	})
}

func ResponseError(c *pine.Context, msg string) {
	c.Render().JSON(pine.H{"result": ResResult{Status: "danger", Message: msg}})
}

type EngineFn func(map[string]string) storage.Uploader

type Engine struct {
	Name   string
	Engine EngineFn
}

func EngineList() []Engine {
	return []Engine{
		{"本地存储", func(opt map[string]string) storage.Uploader {
			return storage.NewFileUploader(opt)
		}},
		{"Oss存储", func(opt map[string]string) storage.Uploader {
			return storage.NewOssUploader(opt)
		}},
		{"Cos存储", func(opt map[string]string) storage.Uploader {
			return storage.NewCosUploader(opt)
		}},
		{"FTP存储", func(opt map[string]string) storage.Uploader {
			return di.MustGet(serviceFtpStorage).(storage.Uploader) // 由于限制链接数, 全局提供单例模式
		}},
	}
}

func GetUserUploader(u *tables.FileManagerAccount) storage.Uploader {
	if u == nil {
		return nil
	}
	u.Engine = "Cos存储"
	for _, v := range EngineList() {
		if v.Name == u.Engine {
			cnf, _ := config.SiteConfig()
			return v.Engine(cnf)
		}
	}
	return nil
}

type FMFileProps struct {
	HasSubdirectories    bool `json:"hasSubdirectories"`
	SubdirectoriesLoaded bool `json:"subdirectoriesLoaded"`
	ShowSubdirectories   bool `json:"showSubdirectories"`
}

type FMFile struct {
	ID        any         `json:"id"`
	Basename  string      `json:"basename"`
	Filename  string      `json:"filename"`
	Dirname   string      `json:"dirname"`
	Path      string      `json:"path"`
	ParentID  string      `json:"parentId"`
	Timestamp int64       `json:"timestamp"`
	ACL       int         `json:"acl"`
	Size      int         `json:"size"`
	Type      string      `json:"type"`
	Extension string      `json:"extension"`
	Props     FMFileProps `json:"props"`
	Author    string      `json:"author"`
}

type DelItem struct {
	Path string `json:"path"`
	Type string `json:"type"`
}
