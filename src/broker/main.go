package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "420"

type Config struct{}

func main() {
	app := &Config{}

	// Create new instance of HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// Start the server
	log.Printf("starting broker service on port %s", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
