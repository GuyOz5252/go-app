package main

import (
	"fmt"
	"log/slog"

	"github.com/GuyOz5252/go-app/internal/services"
)

type config struct {
	address string
}

type application struct {
	logger      *slog.Logger
	config      *config
	UserService *services.UserService
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
