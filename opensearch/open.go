package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"strings"

	"github.com/Millicom-MFS/kit-go/log"
	opensearch "github.com/opensearch-project/opensearch-go"
	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

const IndexName = "go-test-index1"

func main() {
	log.Default()
	// Initialize the client with SSL/TLS enabled.
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{"https://localhost:9200"},
		Username:  "master",
		Password:  "Z45!u7ft",
	})
	if err != nil {
		log.Error("cannot initialize", err)
		os.Exit(1)
	}

	// Print OpenSearch version information on console.
	log.Info(client.Info())

	// Define index mapping.
	mapping := strings.NewReader(`{
     'settings': {
       'index': {
            'number_of_shards': 4
            }
          }
     }`)

	// Create an index with non-default settings.
	res := opensearchapi.IndicesCreateRequest{
		Index: IndexName,
		Body:  mapping,
	}
	log.Info("creating index", res)

	// Add a document to the index.
	document := strings.NewReader(`{
        "title": "Moneyball",
        "director": "Bennett Miller",
        "year": "2011"
    }`)

	// docId := "3"
	req := opensearchapi.IndexRequest{
		Index:      IndexName,
		// DocumentID: docId,
		Body:       document,
	}
	insertResponse, err := req.Do(context.Background(), client)
	if err != nil {
		log.Error("failed to insert document ", err)
		os.Exit(1)
	}
	log.Info(insertResponse)

	// Search for the document.
	content := strings.NewReader(`{
       "size": 5,
       "query": {
           "multi_match": {
           "query": "miller",
           "fields": ["title^2", "director"]
           }
      }
    }`)

	search := opensearchapi.SearchRequest{
		Body: content,
	}

	searchResponse, err := search.Do(context.Background(), client)
	if err != nil {
		log.Error("failed to search document ", err)
		os.Exit(1)
	}
	log.Info(searchResponse)

	// // Delete the document.
	// delete := opensearchapi.DeleteRequest{
	// 	Index:      IndexName,
	// 	DocumentID: docId,
	// }

	// deleteResponse, err := delete.Do(context.Background(), client)
	// if err != nil {
	// 	log.Error("failed to delete document ", err)
	// 	os.Exit(1)
	// }
	// log.Info("deleting document: ", deleteResponse)


	// // Delete previously created index.
	// deleteIndex := opensearchapi.IndicesDeleteRequest{
	// 	Index: []string{IndexName},
	// }

	// deleteIndexResponse, err := deleteIndex.Do(context.Background(), client)
	// if err != nil {
	// 	log.Error("failed to delete index ", err)
	// 	os.Exit(1)
	// }
	// log.Info("deleting index", deleteIndexResponse)
}
