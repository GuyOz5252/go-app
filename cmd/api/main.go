package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/GuyOz5252/go-app/internal/data"
	"github.com/GuyOz5252/go-app/internal/services"
	"github.com/GuyOz5252/go-app/pkg"
	"github.com/go-chi/jwtauth/v5"
)

type Config struct {
	Address          string `yaml:"address"`
	ConnectionString string `yaml:"connection-string"`
	Auth             struct {
		JWTSecret       string `yaml:"jwt-secret"`
		TokenExpiration string `yaml:"token-expiration"`
	} `yaml:"auth"`
	Queries struct {
		User map[string]string `yaml:"user"`
	} `yaml:"queries"`
}

type application struct {
	logger          *slog.Logger
	config          *Config
	tokenAuth       *jwtauth.JWTAuth
	userService     *services.UserService
	tokenExpiration time.Duration
}

func newApplication() (*application, error) {
	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		configPath = "./config/config.yaml"
	}
	config, err := pkg.LoadConfig[Config](configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := data.NewPostgresSqlDb(config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	userRepository := data.NewSqlUserRepository(db, &config.Queries.User)
	tokenAuth := jwtauth.New("HS256", []byte(config.Auth.JWTSecret), nil)

	tokenExpiration, err := time.ParseDuration(config.Auth.TokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token expiration: %w", err)
	}

	app := &application{
		logger:          pkg.NewLogger(),
		config:          &config,
		tokenAuth:       tokenAuth,
		userService:     services.NewUserService(userRepository),
		tokenExpiration: tokenExpiration,
	}

	return app, nil
}

func main() {
	app, err := newApplication()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize application: %s", err))
	}

	server := app.newServer()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app.logger.Info("Listening on port 8080")
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	<-ctx.Done()
	app.logger.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		app.logger.Error("server forced to shutdown", "error", err)
	}

	app.logger.Info("server exited")
}
