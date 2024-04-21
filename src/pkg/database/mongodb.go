package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoConnection(uri string) *mongo.Client {
	counts := 0
	for {
		conn, err := connect(uri)
		if err == nil {
			log.Println("Connected to Mongo!")
			return conn
		}
		counts++
		if counts == 3 {
			log.Println("Failed to connect to Mongo")
			return nil
		}
		log.Println(err)
		time.Sleep(5 * time.Second)
		log.Println("Retrying to connect to Mongo...")
	}
}

func connect(uri string) (*mongo.Client, error) {
	log.Println("Connecting to Mongo...", uri)
	option := options.Client().ApplyURI(uri)
	conn, err := mongo.Connect(context.Background(), option)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
