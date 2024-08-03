package config

import (
	"fmt"
	es "github.com/elastic/go-elasticsearch/v8"
)

func NewESClient(cfg ElasticSearchConfig) (*es.Client, error) {
	ESCfg := es.Config{
		Addresses:         []string{cfg.Host},
		Username:          cfg.UserName,
		Password:          cfg.Password,
		EnableDebugLogger: cfg.Debug,
		MaxRetries:        3,
	}

	client, err := es.NewClient(ESCfg)
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		res, err := client.Indices.Create(cfg.Index)
		if err != nil {
			return nil, err
		}
		if res.StatusCode == 400 {
			fmt.Println("ignore this messages index already created")
		}
	}

	return client, nil
}
