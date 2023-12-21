package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"logger/data"
	"net/http"
	"os"
	"time"
)

// Define constants
const (
	webPort  = "7070"
	rpcPort  = "42069"
	gRpcPort = "69420"
	mongoUrl = "mongodb://mongo:27018"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// Connect to MongoDb
	client, err := connectToMongoDb()
	if err != nil {
		log.Panic(err)
	}

	// Create Context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Create app config
	app := Config{
		Models: data.New(client),
	}

	// Start web server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Println("starting log server...")
	err = server.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

//func (app *Config) serve() {
//	server := &http.Server{
//		Addr:    fmt.Sprintf(":%s", webPort),
//		Handler: app.routes(),
//	}
//
//	log.Println("starting log server...")
//	server.ListenAndServe()
//}

func connectToMongoDb() (*mongo.Client, error) {
	// Create conn options
	clientOpts := options.Client().ApplyURI(mongoUrl)
	clientOpts.SetAuth(options.Credential{
		Username: os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		Password: os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
	})

	log.Printf("this is root username: %s", os.Getenv("LOG_USERNAME"))
	log.Printf("this is root password: %s", os.Getenv("LOG_PASSWORD"))

	// Connect
	conn, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Println("error connecting", err)
		return nil, err
	}
	log.Println("successfully connected to mongo db...")

	return conn, nil
}
