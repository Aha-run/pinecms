package server

import (
	"github.com/xiusin/pine/contracts"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/CloudyKit/jet"
	"github.com/allegro/bigcache/v3"

	"github.com/xiusin/pine"
	"github.com/xiusin/pine/cache/providers/pbigcache"
	"github.com/xiusin/pine/di"
	"github.com/xiusin/pine/middlewares/cache304"
	"github.com/xiusin/pine/render"
	"github.com/xiusin/pine/render/engine/pjet"
	"github.com/xiusin/pine/sessions"
	cacheProvider "github.com/xiusin/pine/sessions/providers/cache"
	"github.com/xiusin/pinecms/src/application/controllers"
	"github.com/xiusin/pinecms/src/application/controllers/taglibs"
	"github.com/xiusin/pinecms/src/application/controllers/tplfun"
	"github.com/xiusin/pinecms/src/common/helper"
	commonLogger "github.com/xiusin/pinecms/src/common/logger"
	"github.com/xiusin/pinecms/src/common/search"
	"github.com/xiusin/pinecms/src/config"
)

var (
	app          = pine.New()
	conf         = config.App()
	cacheHandler contracts.Cache
)

func InitApp() {
	InitDI()
	app.Use(cache304.Cache304(30000*time.Second, conf.StaticPrefixArr()...))
}

func InitCache() {
	cfg := bigcache.DefaultConfig(time.Hour)
	cfg.Shards = 512
	cfg.MaxEntriesInWindow = 1000
	cfg.MaxEntrySize = 100
	cacheHandler = pbigcache.New(cfg)
	if theme, _ := cacheHandler.Get(controllers.CacheTheme); len(theme) > 0 {
		conf.View.Theme = string(theme)
	}
	helper.Inject(controllers.ServiceICache, cacheHandler)
	sess := sessions.New(cacheProvider.NewStore(cacheHandler), &sessions.Config{CookieName: conf.Session.Name, Expires: conf.Session.Expires})
	di.Instance(sess)
}

func InitDI() {
	helper.Inject(controllers.ServiceApplication, app)
	helper.Inject(controllers.ServiceConfig, conf)
	helper.Inject(slog.Default(), initLoggerService())
	helper.Inject(controllers.ServiceJetEngine, initJetEngine())
	helper.Inject(controllers.ServiceSearchName, search.NewZincSearch())

	pine.RegisterViewEngine(di.MustGet(controllers.ServiceJetEngine).(render.AbstractRenderer))
}

func initLoggerService() di.BuildHandler {
	return func(_ di.AbstractBuilder) (i any, e error) {
		ormLogger := commonLogger.NewPineCmsLogger(config.Orm(), 100)
		cmsLogger, err := os.OpenFile(filepath.Join(conf.LogPath, "pinecms.log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		helper.PanicErr(err)

		pine.RegisterOnInterrupt(func() { ormLogger.Close() })
		pine.RegisterOnInterrupt(func() { _ = cmsLogger.Close() })

		var opt = slog.HandlerOptions{AddSource: true}
		_ = debug.SetCrashOutput(cmsLogger, debug.CrashOptions{})
		if config.IsDebug() {
			opt.Level = slog.LevelDebug
		} else {
			opt.Level = slog.LevelWarn
		}
		return slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, ormLogger, cmsLogger), &opt)), nil
	}
}

func initJetEngine() *pjet.PineJet {
	jetEngine := pjet.New(conf.View.FeDirname, ".jet", conf.View.Reload)
	jetEngine.AddPath("resources/taglibs/")
	jetEngine.SetDevelopmentMode(true)
	tags := map[string]jet.Func{
		"flink":          taglibs.Flink,
		"type":           taglibs.Type,
		"adlist":         taglibs.AdList,
		"myad":           taglibs.MyAd,
		"channel":        taglibs.Channel,
		"channelartlist": taglibs.ChannelArtList,
		"artlist":        taglibs.ArcList,
		"pagelist":       taglibs.PageList,
		"list":           taglibs.List,
		"query":          taglibs.Query,
		"likearticle":    taglibs.LikeArticle,
		"hotwords":       taglibs.HotWords,
		"tags":           taglibs.Tags,
		"position":       taglibs.Position,
		"toptype":        taglibs.TopType,
		"format_time":    tplfun.FormatTime,
		"cn_substr":      tplfun.CnSubstr,
		"GetDateTimeMK":  tplfun.GetDateTimeMK,
		"MyDate":         tplfun.MyDate,
	}

	for name, fn := range tags {
		jetEngine.AddGlobalFunc(name, fn)
	}

	return jetEngine
}
