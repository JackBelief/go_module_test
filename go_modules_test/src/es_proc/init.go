package es_proc

import (
	//"context"
	"fmt"

	"github.com/olivere/elastic"
)

var GEsClient *elastic.Client
func Init()(err error) {
	GEsClient, err = NewEsClient()
	return err
}

var EsServerAddr = "http://127.0.0.1:9200/"
func NewEsClient()(*elastic.Client, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(EsServerAddr))
	if err != nil {
		fmt.Println("new es client fail", err.Error())
		return nil, err
	}

	return client, err
}

