package search

import (
	"github.com/blevesearch/bleve/v2/document"
	gse "github.com/vcaesar/gse-bleve"
	"github.com/xiusin/pine"
	"github.com/xiusin/pine/di"
	"github.com/xiusin/pinecms/src/application/controllers"
	"github.com/xiusin/pinecms/src/common/helper"
)

func NewBleve() {
	index, err := gse.New(gse.Option{
		Index: "example.bleve",
		Dicts: "embed, zh",
		Stop:  "",
		Opt:   "search-hmm",
		Trim:  "trim",
	})
	helper.PanicErr(err)
	pine.RegisterOnInterrupt(func() { _ = index.Close() })

	doc := document.NewDocument("1")
	// doc.AddField(document.NewTextField("title", nil, []byte("hello world!")))
	// if err := index.Index(doc.ID(), doc); err != nil {
	// 	helper.PanicErr(err)
	// }

	di.Set(controllers.ServiceSearchName, func(builder di.AbstractBuilder) (any, error) {
		return index, nil
	}, true)

}

// func searchArticles(index bleve.Index, query string) ([]Article, error) {
//     searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
//     searchResult, err := index.Search(searchRequest)
//     if err != nil {
//         return nil, err
//     }
//     var articles []Article
//     for _, hit := range searchResult.Hits {
//         var article Article
//         article.ID, _ = strconv.Atoi(hit.ID)
//         article.Title = hit.Fields["title"].(string)
//         article.Content = hit.Fields["content"].(string)
//         articles = append(articles, article)
//     }
//     return articles, nil
// }
