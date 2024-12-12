package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

type AppConfig struct {
	AccountURL string `config:"ACCOUNT_SERVICE_URL"`
	CatalogRRL string `config:"CATALOG_SERVICE_URL"`
	OrderURL   string `config:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQlServer(cfg.AccountURL, cfg.CatalogRRL, cfg.OrderURL)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/graphql", handler.New(s.ToExecutableSchema()))
	http.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
