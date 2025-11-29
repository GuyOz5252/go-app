package main

import (
	"fmt"
	"log/slog"

	"github.com/GuyOz5252/go-app/internal/services"
)

type Config struct {
	Address          string `yaml:"address"`
	ConnectionString string `yaml:"connection-string"`
	Queries          struct {
		User map[string]string `yaml:"user"`
	} `yaml:"queries"`
}

type application struct {
	logger      *slog.Logger
	config      *Config
	userService *services.UserService
}

func main() {
	app := bootstrap()
	server := app.newServer()

	// TODO: graceful shutdown
	app.logger.Info("Listening on port 8080")
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
