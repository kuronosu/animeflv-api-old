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

// InsertStates insert list of states in states collection
func InsertStates(client *mongo.Client, states []interface{}) (*mongo.InsertManyResult, error) {
	collection := client.Database("deguvon").Collection("states")
	insertManyResult, err := collection.InsertMany(context.TODO(), states)
	return insertManyResult, err
}

// InsertTypes insert list of types in types collection
func InsertTypes(client *mongo.Client, types []interface{}) (*mongo.InsertManyResult, error) {
	collection := client.Database("deguvon").Collection("types")
	insertManyResult, err := collection.InsertMany(context.TODO(), types)
	return insertManyResult, err
}

// InsertGenres insert list of genres in genres collection
func InsertGenres(client *mongo.Client, genres []interface{}) (*mongo.InsertManyResult, error) {
	collection := client.Database("deguvon").Collection("genres")
	insertManyResult, err := collection.InsertMany(context.TODO(), genres)
	return insertManyResult, err
}

// InsertAnimes insert list of animes in animes collection
func InsertAnimes(client *mongo.Client, animes []interface{}) (*mongo.InsertManyResult, error) {
	collection := client.Database("deguvon").Collection("animes")
	insertManyResult, err := collection.InsertMany(context.TODO(), animes)
	return insertManyResult, err
}
