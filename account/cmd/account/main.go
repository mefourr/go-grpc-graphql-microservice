package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/mefourr/go-grpc-graphql-microservice/account"
	"github.com/tinrab/retry"
	"log"
	_ "net/http"
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

	var repository account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		repository, err = account.NewPostgresRepository(cfg.DataBaseURL)
		if err != nil {
			log.Println(err)
		}
		return nil
	})

	defer func(repository account.Repository) {
		err := repository.Close()
		if err != nil {
			log.Println(err)
		}
	}(repository)

	log.Println("It is listening on port: 8080...")
	service := account.NewService(repository)
	log.Fatal(account.ListenGRPC(service, 8080))
}
