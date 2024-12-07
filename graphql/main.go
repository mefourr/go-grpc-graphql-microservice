package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/kelseyhightower/envconfig"
)

type AppCnfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogRRL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppCnfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQlServer(cfg.AccountURL, cfg.CatalogRRL, cfg.OrderURL)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/graphql", handler.GraphQL(s.ToExecutableSchema()))
	http.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}