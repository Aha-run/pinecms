package search

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/xiusin/pine"
	"github.com/xiusin/pine/di"
	"github.com/xiusin/pinecms/src/common/helper"
)

const ServiceSearchName = "search"

// https://vimsky.com/examples/detail/golang-ex-github.com.blevesearch.bleve-Index-Search-method.html
// https://studygolang.com/articles/36305

func NewBleve() {
	mapping := bleve.NewIndexMapping()
	mapping.DefaultAnalyzer = "cn"

	index, err := bleve.New("example.bleve", mapping)
	helper.PanicErr(err)

	pine.RegisterOnInterrupt(func() {
		_ = index.Close()
	})

	di.Set(ServiceSearchName, func(builder di.AbstractBuilder) (interface{}, error) {
		return index, err
	}, true)

}
