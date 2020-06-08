package es

import (
	"gopkg.in/olivere/elastic.v6"
)

var client *elastic.Client

func Init() (err error) {
	client, err = elastic.NewClient(
		elastic.SetURL([]string{"localhost:9200"}...),
		elastic.SetSniff(false),
	)
	return
}

func ES() *elastic.Client {
	return client
}

func Bulk() *elastic.BulkService {
	return client.Bulk()
}
