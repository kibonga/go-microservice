package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const port = "4209"

type Config struct{}

func main() {
	app := &Config{}

	path, err := os.Getwd()
	if err != nil {
		return
	}
	fmt.Println("dir", path)

	log.Printf("starting broker service on port %s", port)

	// create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.Routes(),
	}

	// start the server
	fmt.Println("starting server")
	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("do i get here")
}
