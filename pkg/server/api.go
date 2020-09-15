package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// Api represents the api
type Api struct {
	router http.Handler
}

// Server represents the api
type Server interface {
	Router() http.Handler
}

// New create new server
func New(client *mongo.Client) Server {
	setDbClient(client)
	a := &Api{}
	r := mux.NewRouter()
	// r.HandleFunc("/", a.fetchGopher).Methods(http.MethodGet)
	r.HandleFunc("/types", a.FetchTypes).Methods(http.MethodGet)
	r.HandleFunc("/states", a.FetchStates).Methods(http.MethodGet)
	r.HandleFunc("/genres", a.FetchGenres).Methods(http.MethodGet)
	r.HandleFunc("/animes", a.FetchAnimes).Methods(http.MethodGet)
	r.HandleFunc("/latest", a.FetchLatestEpisodes).Methods(http.MethodGet)
	a.router = r
	return a
}

// Router reutnr the api router
func (a *Api) Router() http.Handler {
	return a.router
}
