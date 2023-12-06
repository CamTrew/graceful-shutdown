package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	client  *mongo.Client
	errChan <-chan error
}

func newMongo() (*Mongo, error) {
	log.Println("Connecting to Mongo...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://mongo:%d", 27017)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to mongo")

	errChan := make(chan error)
	go monitor(client, errChan)

	return &Mongo{
		client:  client,
		errChan: errChan,
	}, nil
}

func monitor(client *mongo.Client, errChan chan error) {
	for {
		if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
			errChan <- errors.New("Lost connection to mongo")
			break
		}
		time.Sleep(5 * time.Second)
	}
}

func (m *Mongo) Err() <-chan error {
	return m.errChan
}

func (m *Mongo) Disconnect(ctx context.Context) error {
	if m.client == nil {
		return nil
	}
	return m.client.Disconnect(ctx)
}
