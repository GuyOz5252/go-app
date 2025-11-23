package main

import "net/http"

type application struct {
	config config
}

type config struct {
	adress string
}

func run(app *application) error {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr: app.config.adress,
		Handler: mux,
	}

	return server.ListenAndServe()
}
