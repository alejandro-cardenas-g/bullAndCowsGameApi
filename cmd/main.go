package main

import (
	"time"

	"github.com/alejandro-cardenas-g/bullAndCowsApp/cmd/api"
	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/utils"
)

func main() {
	cfg := api.ApplicationConfig{
		Addr:            utils.GetEnvironment().GetEnv("API_ADDR", ":3000"),
		GracefulTimeout: time.Second * 15,
	}

	server := api.NewApplication(cfg)

	server.Run()
}
