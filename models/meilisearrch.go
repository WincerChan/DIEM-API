package models

import (
	T "DIEM-API/tools"
	"log"

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
	log.Println(client)
	log.Println(q)
	r, err := client.Search(index).Search(q)
	T.CheckException(err, "Search failed")
	return r
}
