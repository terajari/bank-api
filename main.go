package main

import (
	"log"

	"github.com/terajari/bank-api/delivery"
	"github.com/terajari/bank-api/manager"
	"github.com/terajari/bank-api/utils"
)

func main() {
	cfg, err := utils.LoadConfig("./.env")
	if err != nil {
		log.Fatal("cannot load config")
	}

	infra, err := manager.NewInfraManager(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	repoManager, err := manager.NewRepositoryManager(infra)
	if err != nil {
		log.Fatal(err)
	}
	usecaseManager, err := manager.NewUsecaseManager(repoManager)
	if err != nil {
		log.Fatal(err)
	}
	server, err := delivery.NewServer(cfg, usecaseManager)
	if err != nil {
		log.Fatal(err)
	}
	server.SetupRouter()
	server.Start(cfg.HTTPServer)

}
