package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GuyOz5252/go-app/pkg"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	ConnectionString string `yaml:"connection-string"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up|down>")
	}
	cmd := os.Args[1]

	configPath := "./config/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := pkg.LoadConfig[Config](configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	m, err := migrate.New("file://db/migrations", cfg.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run up migrations: %v", err)
		}
		fmt.Println("Migrations up ran successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run down migrations: %v", err)
		}
		fmt.Println("Migrations down ran successfully")
	default:
		log.Fatalf("unknown command: %s", cmd)
	}
}
