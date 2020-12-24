package db

import (
	"context"
	"fmt"
	"time"

	"github.com/kuronosu/animeflv-api/pkg/scrape"
	"github.com/kuronosu/animeflv-api/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

// SetUp Create mongo client
func SetUp(dbName string, connectionString string) (Manager, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return Manager{}, err
	}

	// Check the connections
	err = client.Ping(ctx, nil)
	if err != nil {
		return Manager{}, err
	}
	return Manager{Client: client, DBName: dbName}, nil
}

// CreateManager Create mongo client and launch fatal when error
func CreateManager(dbname string, connString string) Manager {
	utils.InfoLog("Connect to db")
	manager, err := SetUp(dbname, connString)
	if err != nil {
		utils.FatalLog(err)
	}
	utils.SuccessLog("Connected to db")
	return manager
}

// DropAll collection from db
func (manager *Manager) DropAll() {
	manager.GetCollection("states").Drop(ctx)
	manager.GetCollection("types").Drop(ctx)
	manager.GetCollection("genres").Drop(ctx)
	manager.GetCollection("animes").Drop(ctx)
	manager.GetCollection("latestEpisodes").Drop(ctx)
}

// InsertMany insert data in db
func (manager *Manager) InsertMany(coll string, data ...interface{}) (*mongo.InsertManyResult, error) {
	collection := manager.GetCollection(coll)
	insertManyResult, err := collection.InsertMany(ctx, data)
	return insertManyResult, err
}

// InsertStates insert list of states in states collection
func (manager *Manager) InsertStates(states []scrape.State) (*mongo.InsertManyResult, error) {
	statesI := make([]interface{}, len(states))
	for i, v := range states {
		statesI[i] = v
	}
	return manager.InsertMany("states", statesI...)
}

// InsertTypes insert list of types in types collection
func (manager *Manager) InsertTypes(types []scrape.Type) (*mongo.InsertManyResult, error) {
	typesI := make([]interface{}, len(types))
	for i, v := range types {
		typesI[i] = v
	}
	return manager.InsertMany("types", typesI...)
}

// InsertGenres insert list of genres in genres collection
func (manager *Manager) InsertGenres(genres []scrape.Genre) (*mongo.InsertManyResult, error) {
	genresI := make([]interface{}, len(genres))
	for i, v := range genres {
		genresI[i] = v
	}
	return manager.InsertMany("genres", genresI...)
}

// InsertAnimes insert list of animes in animes collection
func (manager *Manager) InsertAnimes(animes []scrape.Anime) (*mongo.InsertManyResult, error) {
	animesI := make([]interface{}, len(animes))
	for i, v := range animes {
		animesI[i] = v
	}
	return manager.InsertMany("animes", animesI...)
}

// SetLatestEpisodes drop the latestEpisodes after insert data in latestEpisodes collection
func (manager *Manager) SetLatestEpisodes(latestEpisodes []*scrape.LatestEpisode) (*mongo.InsertManyResult, error) {
	if len(latestEpisodes) != 20 {
		return nil, fmt.Errorf("Latest episodes length must be 20, it has %d", len(latestEpisodes))
	}
	collection := manager.GetCollection("latestEpisodes")
	e := collection.Drop(ctx)
	if e != nil {
		return nil, e
	}
	latestEpisodesI := make([]interface{}, len(latestEpisodes))
	for i, v := range latestEpisodes {
		latestEpisodesI[i] = v
	}
	return manager.InsertMany("latestEpisodes", latestEpisodesI...)
}

// UpdateOrInsertAnimes is very self-describing ... :)
func (manager *Manager) UpdateOrInsertAnimes(animes []scrape.Anime) ([]mongo.UpdateResult, []scrape.Anime, []error) {
	collection := manager.GetCollection("animes")
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
func (manager *Manager) LoadStates() ([]scrape.State, error) {
	coll := manager.GetCollection("states")
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

// LoadOneType from db
func (manager *Manager) LoadOneType(id int) (interface{}, error) {
	var result scrape.Type
	coll := manager.GetCollection("types")
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	return result, err
}

// LoadTypes from db
func (manager *Manager) LoadTypes() ([]scrape.Type, error) {
	coll := manager.GetCollection("types")
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
func (manager *Manager) LoadGenres() ([]scrape.Genre, error) {
	coll := manager.GetCollection("genres")
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

// LoadAnimes from db with pagination
func (manager *Manager) LoadAnimes(opts Options) (PaginatedAnimeResult, error) {
	coll := manager.GetCollection("animes")
	animeCount, err := coll.CountDocuments(context.TODO(), bson.D{{}})
	if err != nil {
		return PaginatedAnimeResult{}, err
	}
	const pageSize = 24
	totalPageCount := int(animeCount) / pageSize
	if int(animeCount)%pageSize > 0 {
		totalPageCount++
	}
	if opts.Page > totalPageCount {
		opts.Page = totalPageCount
	} else if opts.Page < 1 {
		opts.Page = 1
	}
	result := PaginatedAnimeResult{
		Page:       opts.Page,
		TotalPages: totalPageCount,
		Count:      int(animeCount),
		Animes:     []scrape.Anime{}}

	op := options.Find()
	op.SetLimit(pageSize)
	op.SetSkip(int64((opts.Page - 1) * pageSize))
	op.SetSort(bson.D{primitive.E{Key: opts.SortField, Value: opts.SortValue}})
	cur, _ := coll.Find(ctx, bson.D{{}}, op)
	for cur.Next(ctx) {
		var a scrape.Anime
		err := cur.Decode(&a)
		if err != nil {
			return result, err
		}
		result.Animes = append(result.Animes, a)
	}
	if err := cur.Err(); err != nil {
		return result, err
	}
	cur.Close(ctx)
	return result, nil
}

// LoadAllAnimes from db
func (manager *Manager) LoadAllAnimes() ([]scrape.Anime, error) {
	coll := manager.GetCollection("animes")
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

// LoadOneAnime from db
func (manager *Manager) LoadOneAnime(flvid int) (scrape.Anime, error) {
	var result scrape.Anime
	err := manager.GetCollection("animes").FindOne(ctx, bson.M{"_id": flvid}).Decode(&result)
	return result, err
}

// SearchAnimeByName from db
func (manager *Manager) SearchAnimeByName(name string) ([]scrape.Anime, error) {
	coll := manager.GetCollection("animes")
	patternName := `.*` + name + `.*`
	nameB := bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: patternName, Options: "i"}}}
	othernamesB := bson.M{"othernames": bson.M{"$regex": primitive.Regex{Pattern: patternName, Options: "i"}}}
	op := options.Find()
	op.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})
	cur, _ := coll.Find(ctx, bson.M{"$or": []interface{}{nameB, othernamesB}}, op)

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
func (manager *Manager) LoadLatestEpisodes() ([]scrape.LatestEpisode, error) {
	coll := manager.GetCollection("latestEpisodes")
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
