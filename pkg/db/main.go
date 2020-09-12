package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetUp Create mongo client
func SetUp() (*mongo.Client, error) {
	host := "localhost"
	port := 27017

	clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port))
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, err
	}

	// Check the connections
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// InsertAnimes insert list of animes in animes collection
func InsertAnimes(client *mongo.Client, animes []interface{}) (*mongo.InsertManyResult, error) {
	collection := client.Database("deguvon").Collection("animes")
	insertManyResult, err := collection.InsertMany(context.TODO(), animes)
	return insertManyResult, err
}
