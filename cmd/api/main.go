package main

import (
	"fmt"

	"github.com/GuyOz5252/go-app/pkg"
)

func main() {
	logger := pkg.NewLogger()
	server := newServer()

	// TODO: graceful shutdown
	logger.Info("Listening on port 8080")
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
