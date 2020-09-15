package server

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kuronosu/deguvon-server-go/pkg/db"
)

var dbClient *mongo.Client

func setDbClient(client *mongo.Client) {
	dbClient = client
}

func (a *Api) FetchTypes(w http.ResponseWriter, r *http.Request) {
	types, _ := db.LoadTypes(dbClient)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

func (a *Api) FetchStates(w http.ResponseWriter, r *http.Request) {
	states, _ := db.LoadStates(dbClient)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(states)
}

func (a *Api) FetchGenres(w http.ResponseWriter, r *http.Request) {
	genres, _ := db.LoadGenres(dbClient)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

func (a *Api) FetchAnimes(w http.ResponseWriter, r *http.Request) {
	animes, _ := db.LoadAnimes(dbClient)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animes)
}

func (a *Api) FetchLatestEpisodes(w http.ResponseWriter, r *http.Request) {
	animes, _ := db.LoadLatestEpisodes(dbClient)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animes)
}
