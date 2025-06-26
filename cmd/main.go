package main

import (
	"log"

	"github.com/alejandro-cardenas-g/bullAndCowsApp/cmd/api"
	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/utils"
)

func main() {
	cfg := api.ApplicationConfig{
		Addr: utils.GetEnvironment().GetEnv("API_ADDR", ":3000"),
	}

	server := api.NewApplication(cfg)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
