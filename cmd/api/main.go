package main

import (
	"fmt"
)

func main() {
	server := newServer()

	// TODO: log server start
	// TODO: graceful shutdown
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
