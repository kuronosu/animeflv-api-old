package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kuronosu/animeflv-api/pkg/db"
	"github.com/kuronosu/animeflv-api/pkg/scrape"
	"go.mongodb.org/mongo-driver/mongo"
)

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

type AnimesResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []scrape.Anime `json:"results"`
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

func assembleAnimesPageLink(result db.PaginatedAnimeResult, next bool) *string {
	newURI := AnimesPath + "?page=%d"
	if next && result.Page < result.TotalPages {
		newURI = fmt.Sprintf(newURI, result.Page+1)
	} else if !next && result.Page >= 2 {
		newURI = fmt.Sprintf(newURI, result.Page-1)
	} else {
		return nil
	}
	return &newURI
}

func HandleAnimes(w http.ResponseWriter, r *http.Request) {
	rawPage := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(rawPage)
	result, err := db.LoadAnimes(dbClient, page)
	if len(result.Animes) == 0 || err != nil {
		internalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, AnimesResponse{
		Count:    result.Count,
		Results:  result.Animes,
		Next:     assembleAnimesPageLink(result, true),
		Previous: assembleAnimesPageLink(result, false)}, http.StatusOK)
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
	animes, _ := db.LoadAllAnimes(dbClient)
	if len(animes) == 0 {
		internalError(w, "Error al cargar datos")
		return
	}
	animesMap := make(map[int]scrape.Anime)
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
	flvid, _ := strconv.Atoi(vars["flvid"])
	anime, err := db.LoadOneAnime(dbClient, flvid)
	if err != nil {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	JSONResponse(w, anime, http.StatusOK)
}
