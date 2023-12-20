package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// Define constants
const (
	webPort  = "80"
	rpcPort  = "42069"
	gRpcPort = "69420"
	mongoUrl = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct{}

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
}

func connectToMongoDb() (*mongo.Client, error) {
	// Create conn options
	clientOpts := options.Client().ApplyURI(mongoUrl)
	clientOpts.SetAuth(options.Credential{
		Username: "admin",
		Password: "assword",
	})

	// Connect
	conn, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Println("error connecting", err)
		return nil, err
	}

	return conn, nil
}
