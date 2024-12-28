package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/mefourr/go-grpc-graphql-microservice/order"
	"github.com/tinrab/retry"
	"log"
	_ "net/http"
	"time"
)

type Config struct {
	DataBaseURL string `config:"DATABASE_URL"`
	AccountURL  string `config:"ACCOUNT_SERVICE_URL"`
	CatalogURL  string `config:"CATALOG_SERVICE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var repository order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		repository, err = order.NewPostgresRepository(cfg.DataBaseURL)
		if err != nil {
			log.Println(err)
		}
		return nil
	})

	defer func(repository order.Repository) {
		err := repository.Close()
		if err != nil {
			log.Println(err)
		}
	}(repository)

	log.Println("It is listening on port: 8080...")
	service := order.NewService(repository)
	log.Fatal(order.ListenGRPC(service, cfg.AccountURL, cfg.CatalogURL, 8080))
}
