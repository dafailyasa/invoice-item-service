package elastic_repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dafailyasa/invoice-item-service/internal/entities"
	"github.com/dafailyasa/invoice-item-service/pkg/pagination"
	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"strconv"
)

type ElasticRepository struct {
	Client *es.Client
	index  *string
}

func NewElasticSearchRepository(client *es.Client, index *string) *ElasticRepository {
	return &ElasticRepository{
		Client: client,
		index:  index,
	}
}

func (es *ElasticRepository) Index(ctx context.Context, invoice entities.ElasticInvoice) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(invoice); err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      *es.index,
		Body:       &buf,
		DocumentID: strconv.FormatUint(invoice.ID, 10),
		Refresh:    "true",
	}

	resp, err := req.Do(ctx, es.Client)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return err
	}

	io.Copy(io.Discard, resp.Body)

	return nil
}
func (es *ElasticRepository) Search(ctx context.Context, pagination *pagination.PaginationRequest) ([]entities.ElasticInvoice, error) {
	should := es.buildShouldQuery(pagination)

	query := es.buildQuery(should, pagination)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{*es.index},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, es.Client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("search request error: %s", resp.String())
	}

	return es.parseResponse(resp, pagination)
}

func (es *ElasticRepository) buildShouldQuery(pagination *pagination.PaginationRequest) []interface{} {
	should := make([]interface{}, 0, 7)

	conditions := map[string]string{
		"invoiceId": pagination.InvoiceID,
		"itemCount": pagination.ItemCount,
		"customer":  pagination.Customer,
		"issueDate": pagination.IssueDate,
		"dueDate":   pagination.DueDate,
		"status":    pagination.Status,
		"subject":   pagination.Keyword,
	}

	for field, value := range conditions {
		if value != "" {
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					field: value,
				},
			})
		}
	}

	return should
}

func (es *ElasticRepository) buildQuery(should []interface{}, pagination *pagination.PaginationRequest) map[string]interface{} {
	query := map[string]interface{}{
		"query": map[string]interface{}{},
	}

	if len(should) == 0 {
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	} else {
		query["query"] = map[string]interface{}{
			"bool": map[string]interface{}{
				"should": should,
			},
		}
	}

	query["from"] = pagination.GetOffset()
	query["size"] = pagination.GetLimit()

	return query
}

func (es *ElasticRepository) parseResponse(resp *esapi.Response, pagination *pagination.PaginationRequest) ([]entities.ElasticInvoice, error) {
	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source entities.ElasticInvoice `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		return nil, err
	}

	res := make([]entities.ElasticInvoice, len(hits.Hits.Hits))
	for i, hit := range hits.Hits.Hits {
		res[i] = hit.Source
	}

	pagination.Total = hits.Hits.Total.Value

	return res, nil
}
