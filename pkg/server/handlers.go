package server

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/kuronosu/deguvon-server-go/pkg/db"
	"github.com/kuronosu/deguvon-server-go/pkg/scrape"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbClient *mongo.Client

func setDbClient(client *mongo.Client) {
	dbClient = client
}

func HandleTypes(w http.ResponseWriter, r *http.Request) {
	types, _ := db.LoadTypes(dbClient)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

func HandleStates(w http.ResponseWriter, r *http.Request) {
	states, _ := db.LoadStates(dbClient)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(states)
}

func HandleGenres(w http.ResponseWriter, r *http.Request) {
	genres, _ := db.LoadGenres(dbClient)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

func HandleAnimes(w http.ResponseWriter, r *http.Request) {
	animes, _ := db.LoadAnimes(dbClient)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animes)
}

func HandleLatestEpisodes(w http.ResponseWriter, r *http.Request) {
	latestEpisodes, _ := db.LoadLatestEpisodes(dbClient)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(latestEpisodes)
}

func HandleDirectory(w http.ResponseWriter, r *http.Request) {
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

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmplt, err := template.ParseFiles("tmpl/index.html")
	// t, err := template.ParseFiles("/pkg/server/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmplt.Execute(w, AllPathsWithoutIndex)
}
