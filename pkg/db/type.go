package db

import (
	"github.com/kuronosu/animeflv-api/pkg/scrape"
	"go.mongodb.org/mongo-driver/mongo"
)

// Manager contains a db client and other utils
type Manager struct {
	Client *mongo.Client
}

func (manager *Manager) GetDB() *mongo.Database {
	return manager.Client.Database("deguvon")
}

func (manager *Manager) GetCollection(collectionName string) *mongo.Collection {
	return manager.GetDB().Collection(collectionName)
}

// Serial represents a sequence document
type Serial struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

// PaginatedAnimeResult result of pagination
type PaginatedAnimeResult struct {
	Page       int
	TotalPages int
	Count      int
	Animes     []scrape.Anime
}

// Options to make query
type Options struct {
	Page      int
	SortField string
	SortValue int
}

// FunctionDataHandler represents method to get data from db
type FunctionDataHandler = func(int) (interface{}, error)
