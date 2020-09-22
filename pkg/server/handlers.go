package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/kuronosu/animeflv-api/pkg/db"
	"github.com/kuronosu/animeflv-api/pkg/scrape"
	"go.mongodb.org/mongo-driver/mongo"
)

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

func internalError(w http.ResponseWriter, err string) {
	JSONResponse(w, ErrorResponse{err, http.StatusInternalServerError}, http.StatusInternalServerError)
}

func JSONResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

var dbClient *mongo.Client

func setDbClient(client *mongo.Client) {
	dbClient = client
}

func HandleTypes(w http.ResponseWriter, r *http.Request) {
	types, _ := db.LoadTypes(dbClient)
	if len(types) == 0 {
		internalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, types, http.StatusOK)
}

func HandleStates(w http.ResponseWriter, r *http.Request) {
	states, _ := db.LoadStates(dbClient)
	if len(states) == 0 {
		internalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, states, http.StatusOK)
}

func HandleGenres(w http.ResponseWriter, r *http.Request) {
	genres, _ := db.LoadGenres(dbClient)
	if len(genres) == 0 {
		internalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, genres, http.StatusOK)
}

func HandleAnimes(w http.ResponseWriter, r *http.Request) {
	animes, _ := db.LoadAnimes(dbClient)
	if len(animes) == 0 {
		internalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, animes, http.StatusOK)
}

func HandleLatestEpisodes(w http.ResponseWriter, r *http.Request) {
	latestEpisodes, _ := db.LoadLatestEpisodes(dbClient)
	if len(latestEpisodes) == 0 {
		internalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, latestEpisodes, http.StatusOK)
}

func HandleDirectory(w http.ResponseWriter, r *http.Request) {
	types, _ := db.LoadTypes(dbClient)
	if len(types) == 0 {
		internalError(w, "Error al cargar datos")
		return
	}
	states, _ := db.LoadStates(dbClient)
	if len(states) == 0 {
		internalError(w, "Error al cargar datos")
		return
	}
	genres, _ := db.LoadGenres(dbClient)
	if len(genres) == 0 {
		internalError(w, "Error al cargar datos")
		return
	}
	animes, _ := db.LoadAnimes(dbClient)
	if len(animes) == 0 {
		internalError(w, "Error al cargar datos")
		return
	}
	animesMap := make(map[string]scrape.Anime)
	for _, anime := range animes {
		animesMap[anime.Flvid] = anime
	}
	JSONResponse(w, scrape.Directory{States: states, Types: types, Genres: genres, Animes: animesMap}, http.StatusOK)
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmplt, err := template.ParseFiles(filepath.Join("tmpl", "index.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	tmplt.Execute(w, AllPathsWithoutIndex)
}

func HandleAnimeDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	anime, err := db.LoadOneAnime(dbClient, vars["flvid"])
	if err != nil {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	JSONResponse(w, anime, http.StatusOK)
}
