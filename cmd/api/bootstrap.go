package main

import (
	"fmt"

	"github.com/GuyOz5252/go-app/internal/data"
	"github.com/GuyOz5252/go-app/internal/services"
	"github.com/GuyOz5252/go-app/pkg"
)

func bootstrap() *application {
	config, err := pkg.LoadConfig[Config]("./../../config/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %s", err))
	}
	db, err := data.NewPostgresSqlDb(config.ConnectionString)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %s", err))
	}

	userRepository := data.NewSqlUserRepository(db, &config.Queries.User)
	app := &application{
		logger: pkg.NewLogger(),
		config: &config,
		userService: services.NewUserService(userRepository),
	}

	return app
}
