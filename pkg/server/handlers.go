package server

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kuronosu/animeflv-api/pkg/db"
	"github.com/kuronosu/animeflv-api/pkg/scrape"
)

// HandleIndex manage the index route
func (api *API) HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmplt, err := template.ParseFiles(filepath.Join("tmpl", "index.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	tmplt.Execute(w, AllPathsWithoutIndex)
}

// HandleTypes manage the types endpoint
func (api *API) HandleTypes(w http.ResponseWriter, r *http.Request) {
	types, _ := db.LoadTypes(api.DB)
	if len(types) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, types, http.StatusOK)
}

// HandleStates manage the states endpoint
func (api *API) HandleStates(w http.ResponseWriter, r *http.Request) {
	states, _ := db.LoadStates(api.DB)
	if len(states) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, states, http.StatusOK)
}

// HandleGenres manage the generes endpoint
func (api *API) HandleGenres(w http.ResponseWriter, r *http.Request) {
	genres, _ := db.LoadGenres(api.DB)
	if len(genres) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, genres, http.StatusOK)
}

// HandleAnimes manage the animelist endpoint
func (api *API) HandleAnimes(w http.ResponseWriter, r *http.Request) {
	rawPage := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(rawPage)
	result, err := db.LoadAnimes(api.DB, page)
	if len(result.Animes) == 0 || err != nil {
		InternalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, AnimesResponse{
		Count:    result.Count,
		Results:  result.Animes,
		Next:     assembleAnimesPageLink(result, true),
		Previous: assembleAnimesPageLink(result, false)}, http.StatusOK)
}

// HandleAnimeDetails manage the anime details endpoint
func (api *API) HandleAnimeDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flvid, _ := strconv.Atoi(vars["flvid"])
	anime, err := db.LoadOneAnime(api.DB, flvid)
	if err != nil {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	JSONResponse(w, anime, http.StatusOK)
}

// HandleLatestEpisodes manage the  latest episodes endpoint
func (api *API) HandleLatestEpisodes(w http.ResponseWriter, r *http.Request) {
	latestEpisodes, _ := db.LoadLatestEpisodes(api.DB)
	if len(latestEpisodes) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, latestEpisodes, http.StatusOK)
}

// HandleDirectory manage the directory endpoint
func (api *API) HandleDirectory(w http.ResponseWriter, r *http.Request) {
	types, _ := db.LoadTypes(api.DB)
	if len(types) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	states, _ := db.LoadStates(api.DB)
	if len(states) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	genres, _ := db.LoadGenres(api.DB)
	if len(genres) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	animes, _ := db.LoadAllAnimes(api.DB)
	if len(animes) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	animesMap := make(map[int]scrape.Anime)
	for _, anime := range animes {
		animesMap[anime.Flvid] = anime
	}
	JSONResponse(w, scrape.Directory{States: states, Types: types, Genres: genres, Animes: animesMap}, http.StatusOK)
}
