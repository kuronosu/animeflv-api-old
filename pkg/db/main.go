package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kuronosu/deguvon-server-go/pkg/scrape"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

// SetUp Create mongo client
func SetUp() (*mongo.Client, error) {
	connectionString := os.Getenv("MongoConnectionString")
	if connectionString == "" {
		connectionString = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
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
func InsertStates(client *mongo.Client, states []scrape.State) (*mongo.InsertManyResult, error) {
	statesI := make([]interface{}, len(states))
	for i, v := range states {
		statesI[i] = v
	}
	collection := client.Database("deguvon").Collection("states")
	insertManyResult, err := collection.InsertMany(ctx, statesI)
	return insertManyResult, err
}

// InsertTypes insert list of types in types collection
func InsertTypes(client *mongo.Client, types []scrape.Type) (*mongo.InsertManyResult, error) {
	typesI := make([]interface{}, len(types))
	for i, v := range types {
		typesI[i] = v
	}
	collection := client.Database("deguvon").Collection("types")
	insertManyResult, err := collection.InsertMany(ctx, typesI)
	return insertManyResult, err
}

// InsertGenres insert list of genres in genres collection
func InsertGenres(client *mongo.Client, genres []scrape.Genre) (*mongo.InsertManyResult, error) {
	genresI := make([]interface{}, len(genres))
	for i, v := range genres {
		genresI[i] = v
	}
	collection := client.Database("deguvon").Collection("genres")
	insertManyResult, err := collection.InsertMany(ctx, genresI)
	return insertManyResult, err
}

// InsertAnimes insert list of animes in animes collection
func InsertAnimes(client *mongo.Client, animes []scrape.Anime) (*mongo.InsertManyResult, error) {
	animesI := make([]interface{}, len(animes))
	for i, v := range animes {
		animesI[i] = v
	}
	collection := client.Database("deguvon").Collection("animes")
	insertManyResult, err := collection.InsertMany(ctx, animesI)
	return insertManyResult, err
}

// SetLatestEpisodes drop the latestEpisodes after insert data in latestEpisodes collection
func SetLatestEpisodes(client *mongo.Client, latestEpisodes []*scrape.LatestEpisode) (*mongo.InsertManyResult, error) {
	if len(latestEpisodes) != 20 {
		return nil, fmt.Errorf("Latest episodes length must be 20, it has %d", len(latestEpisodes))
	}
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
func UpdateOrInsertAnimes(client *mongo.Client, animes []scrape.Anime) ([]mongo.UpdateResult, []scrape.Anime, []error) {
	collection := client.Database("deguvon").Collection("animes")
	animesInterface := make([]interface{}, len(animes))
	for i, v := range animes {
		animesInterface[i] = v
	}
	errors := []error{}
	updateResults := []mongo.UpdateResult{}
	toInsert := []interface{}{}
	for _, anime := range animesInterface {
		r, err := collection.UpdateOne(ctx, bson.M{"_id": anime.(scrape.Anime).Flvid},
			bson.D{{Key: "$set", Value: anime}})
		if err != nil {
			errors = append(errors, err)
			continue
		}
		if r.MatchedCount == 0 {
			toInsert = append(toInsert, anime)
		} else {
			updateResults = append(updateResults, *r)
		}
	}
	var inResult *mongo.InsertManyResult
	var insertedAnimes []scrape.Anime
	if len(toInsert) > 0 {
		a, err := collection.InsertMany(ctx, toInsert)
		inResult = a
		if err != nil {
			errors = append(errors, err)
		}
		for _, id := range inResult.InsertedIDs {
			for _, anime := range animes {
				if id == anime.Flvid {
					insertedAnimes = append(insertedAnimes, anime)
				}
			}
		}
	}
	return updateResults, insertedAnimes, errors
}

// LoadStates from db
func LoadStates(client *mongo.Client) ([]scrape.State, error) {
	coll := client.Database("deguvon").Collection("states")
	cur, _ := coll.Find(ctx, bson.D{{}}, options.Find())
	var results []scrape.State
	for cur.Next(ctx) {
		var s scrape.State
		err := cur.Decode(&s)
		if err != nil {
			return results, err
		}
		results = append(results, s)
	}
	if err := cur.Err(); err != nil {
		return results, err
	}
	cur.Close(ctx)
	return results, nil
}

// LoadTypes from db
func LoadTypes(client *mongo.Client) ([]scrape.Type, error) {
	coll := client.Database("deguvon").Collection("types")
	var results []scrape.Type
	cur, err := coll.Find(ctx, bson.D{{}}, options.Find())
	if err != nil {
		return results, err
	}
	for cur.Next(ctx) {
		var s scrape.Type
		err := cur.Decode(&s)
		if err != nil {
			return results, err
		}
		results = append(results, s)
	}
	if err := cur.Err(); err != nil {
		return results, err
	}
	cur.Close(ctx)
	return results, nil
}

// LoadGenres from db
func LoadGenres(client *mongo.Client) ([]scrape.Genre, error) {
	coll := client.Database("deguvon").Collection("genres")
	cur, _ := coll.Find(ctx, bson.D{{}}, options.Find())
	var results []scrape.Genre
	for cur.Next(ctx) {
		var s scrape.Genre
		err := cur.Decode(&s)
		if err != nil {
			return results, err
		}
		results = append(results, s)
	}
	if err := cur.Err(); err != nil {
		return results, err
	}
	cur.Close(ctx)
	return results, nil
}

// LoadAnimes from db
func LoadAnimes(client *mongo.Client) ([]scrape.Anime, error) {
	coll := client.Database("deguvon").Collection("animes")
	cur, _ := coll.Find(ctx, bson.D{{}}, options.Find())
	var results []scrape.Anime
	for cur.Next(ctx) {
		var s scrape.Anime
		err := cur.Decode(&s)
		if err != nil {
			return results, err
		}
		results = append(results, s)
	}
	if err := cur.Err(); err != nil {
		return results, err
	}
	cur.Close(ctx)
	return results, nil
}

// LoadLatestEpisodes from db
func LoadLatestEpisodes(client *mongo.Client) ([]scrape.LatestEpisode, error) {
	coll := client.Database("deguvon").Collection("latestEpisodes")
	cur, _ := coll.Find(ctx, bson.D{{}}, options.Find())
	var results []scrape.LatestEpisode
	for cur.Next(ctx) {
		var s scrape.LatestEpisode
		err := cur.Decode(&s)
		if err != nil {
			return results, err
		}
		results = append(results, s)
	}
	if err := cur.Err(); err != nil {
		return results, err
	}
	cur.Close(ctx)
	return results, nil
}
