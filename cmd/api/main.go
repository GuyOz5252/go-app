package main

import (
	"fmt"
	"log/slog"

	"github.com/GuyOz5252/go-app/internal/data"
	"github.com/GuyOz5252/go-app/internal/services"
	"github.com/GuyOz5252/go-app/pkg"
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

func newApplication() (*application, error) {
	config, err := pkg.LoadConfig[Config]("./../../config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	
	db, err := data.NewPostgresSqlDb(config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	userRepository := data.NewSqlUserRepository(db, &config.Queries.User)
	app := &application{
		logger:      pkg.NewLogger(),
		config:      &config,
		userService: services.NewUserService(userRepository),
	}

	return app, nil
}

func main() {
	app, err := newApplication()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize application: %s", err))
	}

	server := app.newServer()

	// TODO: graceful shutdown
	app.logger.Info("Listening on port 8080")
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
