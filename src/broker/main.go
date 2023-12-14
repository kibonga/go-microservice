package main

import (
	"fmt"
	"log"
	"modules/helpers"
	"net/http"
	"os"
)

const port = "420"

type Config struct{}

func main() {
	app := &Config{}
	app123 := &helpers.Config123{}
	app123.ReadJSON()

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
