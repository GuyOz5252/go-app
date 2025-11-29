package main

import (
	"fmt"

	"github.com/GuyOz5252/go-app/internal/data"
	"github.com/GuyOz5252/go-app/internal/services"
	"github.com/GuyOz5252/go-app/pkg"
)

func bootstrap() *application {
	db, err := data.NewPostgresSqlDb("postgres://go-app-server:Password1@localhost:5432/go-app-db?sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %s", err))
	}
	userRepository := data.NewSqlUserRepository(db)
	app := &application{
		logger: pkg.NewLogger(),
		config: &config{
			address: ":8080",
		},
		UserService: services.NewUserService(userRepository),
	}

	return app
}
