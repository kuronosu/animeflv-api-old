package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kuronosu/animeflv-api/pkg/scrape"
	"go.mongodb.org/mongo-driver/mongo"
)

// API represents the api
type API struct {
	router http.Handler
	port   int
	DB     *mongo.Client
}

// Server represents the api
type Server interface {
	Router() http.Handler
	Run() error
}

// Router return the api router
func (a *API) Router() http.Handler {
	return LogMiddleware(TrailingSlashes(CaselessMatcher(a.router)))
}

// Run start the server
func (a *API) Run() error {
	log.Printf("Listen to :%d", a.port)
	return http.ListenAndServe(fmt.Sprint(":", a.port), a.Router())
}

// ErrorResponse rendered in json
type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

// AnimesResponse rendered anime data in json
type AnimesResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []scrape.Anime `json:"results"`
}

// EpisodesResponse rendered list of episodes data with extra anime data in json
type EpisodesResponse struct {
	AnimeID   int              `json:"animeID"`
	AnimeName string           `json:"animeName"`
	AnimeURL  string           `json:"animeURL"`
	Episodes  []scrape.Episode `json:"episodes"`
}

// EpisodeResponse rendered episode data with extra anime data in json
type EpisodeResponse struct {
	AnimeID   int            `json:"animeID"`
	AnimeName string         `json:"animeName"`
	AnimeURL  string         `json:"animeURL"`
	Episode   scrape.Episode `json:"episode"`
}
