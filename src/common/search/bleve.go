package search

import (
	gse "github.com/vcaesar/gse-bleve"
	"github.com/xiusin/pine"
	"github.com/xiusin/pinecms/src/application/controllers"
	"github.com/xiusin/pinecms/src/common/helper"
)

func NewBleve() {
	indexPath := helper.GetRootPath("runtime/document.bleve")

	index, err := gse.New(gse.Option{
		Index: indexPath,
		Dicts: "embed, zh",
		Opt:   "search-hmm",
		Trim:  "trim",
	})
	helper.PanicErr(err)
	pine.RegisterOnInterrupt(func() { _ = index.Close() })
	pine.Logger().Debug("启动搜索引擎")
	helper.Inject(controllers.ServiceSearchName, index)
}
