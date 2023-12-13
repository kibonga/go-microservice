package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

// Define a port to listen to
const webPort = 420

// Config Define a default config type, which will store relevant data
// This type can be used as a default type for creating interface for methods called in main
type Config struct {
	DB sql.DB
}

func main() {
	//app := Config{}

	var server = &http.Server{
		Addr: fmt.Sprintf(":%d", webPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Hello from http server handler")
		}),
	}

	fmt.Println("Starting auth server...")
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
