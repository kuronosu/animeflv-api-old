package server

import (
	"encoding/json"
	"net/http"

	"github.com/kuronosu/deguvon-server-go/pkg/db"
	"github.com/kuronosu/deguvon-server-go/pkg/scrape"
	"go.mongodb.org/mongo-driver/mongo"
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
	latestEpisodes, _ := db.LoadLatestEpisodes(dbClient)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(latestEpisodes)
}

func (a *Api) FetchDirectory(w http.ResponseWriter, r *http.Request) {
	types, _ := db.LoadTypes(dbClient)
	states, _ := db.LoadStates(dbClient)
	genres, _ := db.LoadGenres(dbClient)
	animes, _ := db.LoadAnimes(dbClient)
	animesMap := make(map[string]scrape.Anime)
	for _, anime := range animes {
		animesMap[anime.Flvid] = anime
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scrape.Directory{
		States: states, Types: types,
		Genres: genres, Animes: animesMap})
}
