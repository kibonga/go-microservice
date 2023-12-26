package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	webPort = "80"
)

type Config struct {
	Mailer Mail
}

func main() {
	app := Config{
		Mailer: createMail(),
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Println("starting mailer server...")
	err := server.ListenAndServe()
	if err != nil {
		log.Panic("failed to start mailer server...", err)
		return
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	return Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}
}
