package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type LogEntry struct {
	Id        string    `bson:"_id:omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Models struct {
	LogEntry LogEntry
}

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logger").Collection("logs")

	log.Println("inserting one", entry)
	log.Println("collection:", collection.Name())
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error inserting into logs", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logger").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("error finding all logs", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	var entry LogEntry
	for cursor.Next(ctx) {
		err := cursor.Decode(&entry)
		if err != nil {
			log.Println("error decoding log into slice", err)
			return nil, err
		}

		logs = append(logs, &entry)
	}

	return logs, nil
}

func (l *LogEntry) GetById(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logger").Collection("logs")

	logId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error parsing to mongo id", err)
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": logId}).Decode(&entry)
	if err != nil {
		log.Println("error finding log entry", err)
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logger").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logger").Collection("logs")

	logId, err := primitive.ObjectIDFromHex(l.Id)
	if err != nil {
		return nil, err
	}

	res, err := collection.UpdateOne(ctx, bson.M{"_id": logId}, bson.D{
		{"$set", bson.D{
			{"name", l.Name},
			{"data", l.Data},
			{"updated_at", time.Now()},
		},
		},
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
