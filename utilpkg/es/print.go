package es

import (
	"encoding/json"
	"gopkg.in/olivere/elastic.v6"
	"log"
)

func PrintSource(query *elastic.BoolQuery) {
	src, _ := query.Source()
	data, err := json.MarshalIndent(src, "", "  ")
	log.Println("PrintSource data:", string(data), err)
}
