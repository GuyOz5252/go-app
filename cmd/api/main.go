package main

import "log"

func main() {
	app := &application{
		config: config{
			adress: ":8080",
		},
	}

	if err := run(app); err != nil {
		log.Fatal(err)
	}
}
