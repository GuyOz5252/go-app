package main

import (
	"fmt"
	"log/slog"

	"github.com/GuyOz5252/go-app/internal/data"
	"github.com/GuyOz5252/go-app/internal/services"
	"github.com/GuyOz5252/go-app/pkg"
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
	db, err := data.NewPostgresSqlDb("postgres://go-app-server:Password1@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %s", err))
	}
	userRepository := data.NewSqlUserRepository(db)
	app := &application{
		logger:      pkg.NewLogger(),
		config: &config{
			address: ":8080",
		},
		UserService: services.NewUserService(userRepository),
	}

	server := app.newServer()

	// TODO: graceful shutdown
	app.logger.Info("Listening on port 8080")
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
