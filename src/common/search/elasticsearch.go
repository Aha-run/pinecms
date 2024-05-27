package search

import (
	"bytes"
	"encoding/json"

	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/xiusin/pinecms/src/common/helper"
)

type ElasticSearch struct {
	client *elasticsearch8.Client
}

func (e *ElasticSearch) Search(index string, query any) (any, error) {
	panic("not implemented") // TODO: Implement
}

func (e *ElasticSearch) Index(index string, doc map[string]any) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (e *ElasticSearch) Update(index string, id string, doc map[string]any) error {
	byts, err := json.Marshal(doc)
	if err != nil {
		return nil
	}
	_, err = e.client.Update(index, id, bytes.NewBuffer(byts))
	return err
}

func (e *ElasticSearch) Delete(index string, id string) error {
	_, err := e.client.Delete(index, id)
	return err
}

func NewElasticSearch() ISearch {
	es8, err := elasticsearch8.NewDefaultClient()
	helper.PanicErr(err)

	return &ElasticSearch{client: es8}
}
