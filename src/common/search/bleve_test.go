package search

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/blevesearch/bleve/v2"
	gse "github.com/vcaesar/gse-bleve"
)

type TestItem struct {
	Id      int64
	Title   string
	Keyword string
	Pubtime time.Time
	Status  bool
	Extra   any
}

func TestBleve(t *testing.T) {
	opt := gse.Option{
		Index: "test.blv",
		Dicts: "embed, zh",
		Stop:  "",
		Opt:   "search-hmm",
		Trim:  "trim",
	}

	index, err := gse.New(opt)
	if err != nil {
		fmt.Println("new mapping error is: ", err)
		return
	}

	err = index.Index("1", TestItem{
		Id:      1,
		Title:   "我是标题",
		Keyword: "我,标题,关键字",
		Pubtime: time.Now(),
		Status:  true,
		Extra: map[string]any{
			"cash": 1,
			"time": time.Now(),
		},
	})

	if err != nil {
		fmt.Println("index error: ", err)
	}

	query := "标题"
	req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	req.Highlight = bleve.NewHighlight()
	res, err := index.Search(req)

	// 处理查询结果
	for _, hit := range res.Hits {
		var p TestItem
		// 从索引中获取匹配的文档
		fmt.Println("hit.ID", hit.ID)
		document, err := index.GetInternal([]byte(hit.ID))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(document))
		fmt.Println(json.Unmarshal(document, &p))
		fmt.Printf("Found: %+v\n", p)
	}
	os.RemoveAll("test.blv")
}
