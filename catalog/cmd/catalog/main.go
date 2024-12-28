package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/mefourr/go-grpc-graphql-microservice/catalog"
	"github.com/tinrab/retry"
	"log"
	"time"
)

type Config struct {
	DataBaseURL string `config:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var repository catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		repository, err = catalog.NewElasticRepository(cfg.DataBaseURL)
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	defer repository.Close()

	log.Println("Catalog is listening on port: 8080...")
	service := catalog.NewService(repository)
	log.Fatal(catalog.ListenGRPC(service, 8080))
}
