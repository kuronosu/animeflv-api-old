package db

import "github.com/kuronosu/animeflv-api/pkg/scrape"

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
