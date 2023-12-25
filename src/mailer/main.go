package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	webPort = "80"
)

type Config struct {
}

func main() {
	//app := Config{}

	server := http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Hello from mailer handler fn!")
		}),
	}

	log.Println("starting mailer server...")
	err := server.ListenAndServe()
	if err != nil {
		log.Panic("failed to start mailer server...", err)
		return
	}
}
