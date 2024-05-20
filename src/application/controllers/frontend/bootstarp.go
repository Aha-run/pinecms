package frontend

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xiusin/pine"
	"github.com/xiusin/pinecms/src/application/models"
	"github.com/xiusin/pinecms/src/config"
)

const IndexTpl = "editor.tpl"

func (c *IndexController) Bootstrap() {
	begin := time.Now()
	defer func() {
		pine.Logger().Debug("请求+渲染总耗时", time.Since(begin))
		if err := recover(); err != nil {
			c.Ctx().Abort(http.StatusInternalServerError, err.(error).Error())
		}
	}()
	// todo 开启前端资源缓存 304
	// todo 拦截存在静态文件的问题, 不过最好交给nginx等服务器转发
	pageName := c.Ctx().Params().Get("pagename") // 必须包含.html, 在nginx要注意如果以/结尾的path需要追加index.html
	if config.GetSiteConfigByKey("SITE_DEBUG", "关闭") == "关闭" {
		if pageName == "" {
			pageName = "/"
		}
		if strings.HasSuffix(pageName, "/") {
			pageName += IndexTpl
		}
		absFilePath := filepath.Join(config.GetSiteConfigByKey("SITE_STATIC_PAGE_DIR"), pageName)
		if byts, err := os.ReadFile(absFilePath); err == nil { // 如果已经存在缓存页面则直接并发执行
			c.Ctx().Render().ContentType(pine.ContentTypeHTML)
			pine.Logger().Print("render for file", absFilePath)
			_ = c.Ctx().Render().Bytes(byts)
			return
		}
	}

	pageName = strings.Trim(strings.ReplaceAll(c.Ctx().Params().Get("pagename"), "//", "/"), "/") // 必须包含.html

	switch pageName {
	case IndexTpl, "", "/":
		c.Index()
	default:
		urlPartials := strings.Split(pageName, "/")
		var last string
		var fileName string
		var isDetail bool
		// 如果地址内包含 .html 认为是需要请求静态页面
		if strings.HasSuffix(pageName, ".html") {
			last = urlPartials[len(urlPartials)-2]     // 获取目录名
			fileName = urlPartials[len(urlPartials)-1] // 获取文件名
			if strings.HasPrefix(fileName, "index_") { // 分析页码
				fileInfo := strings.Split(fileName, "_") // index_2.html => 某个分类的第二页
				c.Ctx().Params().Set("page", strings.TrimSuffix(fileInfo[1], ".html"))
			} else if fileName != IndexTpl {
				isDetail = true
				c.Ctx().Params().Set("aid", strings.TrimSuffix(fileName, ".html")) // 设置文档ID
			}
		} else {
			last = urlPartials[len(urlPartials)-1] // 目录名
			fileName = IndexTpl
			pageName = filepath.Join(pageName, fileName)
			c.Ctx().Params().Set("page", "1")
		}
		cat := models.NewCategoryModel().GetWithDirForBE(last)
		if cat == nil {
			// 拆分出来想要的数据 page_{tid}
			if strings.HasPrefix(last, "page_") {
				c.Ctx().Params().Set("tid", strings.TrimPrefix(last, "page_"))
				c.Page(pageName)
				return
			} else {
				infos := strings.Split(last, "_") // 根据模型{model_table}_{tid}拆分信息
				modelTable := infos[0]
				model := models.NewDocumentModel().GetWithTableNameForBE(modelTable)
				if model != nil {
					c.Ctx().Params().Set("tid", infos[1])
					c.List(pageName)
					return
				}
			}
			c.Ctx().Abort(http.StatusNotFound)
			c.Logger().Debug("无法匹配列表或单页", c.Ctx().Path())
			return
		}
		// 匹配所有内容
		prefix := models.NewCategoryModel().GetUrlPrefix(cat.Catid)
		if !strings.HasPrefix(pageName, prefix) {
			c.Logger().Debug("2:地址前缀无法匹配", c.Ctx().Path())
			c.Ctx().Abort(http.StatusNotFound)
			return
		}
		c.Ctx().Params().Set("tid", strconv.Itoa(int(cat.Catid)))
		if isDetail {
			c.Detail(pageName)
		} else if cat.Type == 0 {
			c.List(pageName)
		} else {
			c.Page(pageName)
		}
	}
}
