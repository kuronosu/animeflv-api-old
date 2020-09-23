package server

import (
	"net/http"

	"github.com/kuronosu/animeflv-api/pkg/scrape"
	"go.mongodb.org/mongo-driver/mongo"
)

// API represents the api
type API struct {
	router http.Handler
	DB     *mongo.Client
}

// Server represents the api
type Server interface {
	Router() http.Handler
}

// Router return the api router
func (a *API) Router() http.Handler {
	return CaselessMatcher(a.router)
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
