package cmd

import (
	"context"
	"encoding/json"
	"lib-log/pkg/dto"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/uuid"
)

func Init(body *dto.LogDTO) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Create an index
	indexName := "test-index"
	body.Timestamp = time.Now().UTC()
	dados, _ := json.MarshalIndent(body, "", "   ")

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: uuid.NewString(),
		Body:       strings.NewReader(string(dados)),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error indexing document: %s", res.Status())
	}

	// Print the Elasticsearch response
	log.Printf("Elasticsearch response: %s", res.Status())
}
