package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/kuronosu/deguvon-server-go/pkg/scrape"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

// SetUp Create mongo client
func SetUp() (*mongo.Client, error) {
	host := "localhost"
	port := 27017

	clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// Check the connections
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// InsertStates insert list of states in states collection
func InsertStates(client *mongo.Client, states []interface{}) (*mongo.InsertManyResult, error) {
	collection := client.Database("deguvon").Collection("states")
	insertManyResult, err := collection.InsertMany(ctx, states)
	return insertManyResult, err
}

// InsertTypes insert list of types in types collection
func InsertTypes(client *mongo.Client, types []interface{}) (*mongo.InsertManyResult, error) {
	collection := client.Database("deguvon").Collection("types")
	insertManyResult, err := collection.InsertMany(ctx, types)
	return insertManyResult, err
}

// InsertGenres insert list of genres in genres collection
func InsertGenres(client *mongo.Client, genres []interface{}) (*mongo.InsertManyResult, error) {
	collection := client.Database("deguvon").Collection("genres")
	insertManyResult, err := collection.InsertMany(ctx, genres)
	return insertManyResult, err
}

// InsertAnimes insert list of animes in animes collection
func InsertAnimes(client *mongo.Client, animes []interface{}) (*mongo.InsertManyResult, error) {
	collection := client.Database("deguvon").Collection("animes")
	insertManyResult, err := collection.InsertMany(ctx, animes)
	return insertManyResult, err
}

// SetLatestEpisodes drop the latestEpisodes after insert data in latestEpisodes collection
func SetLatestEpisodes(client *mongo.Client, latestEpisodes []scrape.LatestEpisode) (*mongo.InsertManyResult, error) {
	collection := client.Database("deguvon").Collection("latestEpisodes")
	e := collection.Drop(ctx)
	if e != nil {
		return nil, e
	}
	collection = client.Database("deguvon").Collection("latestEpisodes")
	latestEpisodesInterface := make([]interface{}, len(latestEpisodes))
	for i, v := range latestEpisodes {
		latestEpisodesInterface[i] = v
	}
	insertManyResult, err := collection.InsertMany(ctx, latestEpisodesInterface)
	return insertManyResult, err
}

// UpdateOrInsertAnimes is very self-describing ... :)
func UpdateOrInsertAnimes(client *mongo.Client, animes []scrape.Anime) ([]*mongo.UpdateResult, []*mongo.InsertOneResult, []error) {
	collection := client.Database("deguvon").Collection("animes")
	animesInterface := make([]interface{}, len(animes))
	for i, v := range animes {
		animesInterface[i] = v
	}
	errors := []error{}
	updateResults := []*mongo.UpdateResult{}
	insertResults := []*mongo.InsertOneResult{}
	for _, anime := range animesInterface {
		r, err := collection.UpdateOne(ctx, bson.M{"_id": anime.(scrape.Anime).Flvid}, bson.D{{Key: "$set", Value: anime}})
		if err != nil {
			errors = append(errors, err)
			continue
		}
		updateResults = append(updateResults, r)
		if r.MatchedCount == 0 {
			r2, e := collection.InsertOne(ctx, anime)
			if e != nil {
				errors = append(errors, e)
				continue
			}
			insertResults = append(insertResults, r2)

		}
	}
	return updateResults, insertResults, errors
}
