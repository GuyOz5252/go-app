package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/GuyOz5252/go-app/internal/data"
	"github.com/GuyOz5252/go-app/internal/handlers"
	"github.com/GuyOz5252/go-app/internal/services"
	"github.com/GuyOz5252/go-app/internal/services/websocket"
	"github.com/GuyOz5252/go-app/pkg"
	"github.com/go-chi/jwtauth/v5"
)

type Config struct {
	Address          string `mapstructure:"address"`
	ConnectionString string `mapstructure:"connection-string"`
	Auth             struct {
		JWTSecret       string        `mapstructure:"jwt-secret"`
		TokenExpiration time.Duration `mapstructure:"token-expiration"`
	} `mapstructure:"auth"`
	Queries struct {
		User map[string]string `mapstructure:"user"`
		Chat map[string]string `mapstructure:"chat"`
	} `mapstructure:"queries"`
}

func newApplication() (*application, error) {
	config, err := pkg.LoadConfig[Config]("../../config")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := data.NewPostgresSqlDb(config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	tokenAuth := jwtauth.New("HS256", []byte(config.Auth.JWTSecret), nil)

	healthHandler := handlers.NewHealthHandler()

	userRepository := data.NewSqlUserRepository(db, &config.Queries.User)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService, tokenAuth, config.Auth.TokenExpiration)

	hub := websocket.NewHub()
	go hub.Run()

	chatService := services.NewChatService(nil, hub)
	chatHandler := handlers.NewChatHandler(chatService)
	websocketHandler := handlers.NewWebSocketHandler(hub)

	app := &application{
		config:           config,
		logger:           pkg.NewLogger(),
		db:               db,
		tokenAuth:        tokenAuth,
		healthHandler:    healthHandler,
		userHandler:      userHandler,
		chatHandler:      chatHandler,
		websocketHandler: websocketHandler,
	}

	return app, nil
}

func main() {
	app, err := newApplication()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize application: %s", err))
	}
	defer app.db.Close()

	server := app.newHttpServer(app.mount())

	app.logger.Info("Listening on port 8080", slog.String("address", app.config.Address))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
