package main

import (
	"fmt"
)

func main() {
	server := newServer()

	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
