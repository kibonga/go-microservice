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
	mongoUrl = "mongodb://mongo:27017"
)

type Config struct {
	Models data.Models
}

var client *mongo.Client

func main() {
	// Connect to MongoDb
	mongoClient, err := connectToMongoDb()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

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
		Username: "admin",
		Password: "assword",
	})

	log.Printf("this is root username: %s", os.Getenv("LOG_USERNAME"))
	log.Printf("this is root password: %s", os.Getenv("LOG_PASSWORD"))
	log.Println("client options", clientOpts)

	// Create new Mongo client
	mongoClient, err := mongo.Connect(context.TODO(), clientOpts)
	log.Println("this is mongo conn:", mongoClient)
	log.Println("testing(ping)", mongoClient.Ping(context.TODO(), nil))
	if err != nil {
		log.Println("error connecting", err)
		return nil, err
	}
	log.Println("successfully connected to mongo db...")

	return mongoClient, nil
}
