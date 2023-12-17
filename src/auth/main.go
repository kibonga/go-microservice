package main

import (
	"auth/data"
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
	DB     *sql.DB
	Models data.Models
}

func main() {
	// Connect to database
	conn := connectToDb()
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	var server = &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.routes(),
	}

	fmt.Println("Starting auth server...")
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
