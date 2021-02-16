package models

import (
	T "DIEM-API/tools"

	meili "github.com/meilisearch/meilisearch-go"
)

var (
	client meili.ClientInterface
)

func InitMeiliClient(host, apiKey string) {
	client = meili.NewClient(meili.Config{
		Host:   host,
		APIKey: apiKey,
	})
}

func Execute(index string, q meili.SearchRequest) *meili.SearchResponse {
	r, err := client.Search(index).Search(q)
	T.CheckException(err, "Search failed")
	return r
}
