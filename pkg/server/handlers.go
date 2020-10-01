package server

import (
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kuronosu/animeflv-api/pkg/db"
	"github.com/kuronosu/animeflv-api/pkg/scrape"
	"go.mongodb.org/mongo-driver/mongo"
)

var baseTemplate = filepath.Join("tmpl", "base.html")

func getAnime(r *http.Request, client *mongo.Client) (scrape.Anime, error) {
	vars := mux.Vars(r)
	flvid, err := strconv.Atoi(vars["flvid"])
	if err != nil {
		return scrape.Anime{}, err
	}
	return db.LoadOneAnime(client, flvid)
}

func getEpisode(r *http.Request, client *mongo.Client) (*EpisodeResponse, error) {
	eNumber, err := strconv.ParseFloat(mux.Vars(r)["eNumber"], 64)
	if err != nil {
		return nil, err
	}
	if eNumber < 0 {
		return nil, errors.New("The episode number must be greater than or equal to zero")
	}
	anime, err := getAnime(r, client)
	if err != nil {
		return nil, err
	}
	for _, episode := range anime.Episodes {
		if episode.Number == eNumber {
			return &EpisodeResponse{AnimeID: anime.Flvid, AnimeName: anime.Name,
				AnimeURL: anime.URL, Episode: episode}, nil
		}
	}
	return nil, errors.New("Episode not found")
}

// HandleIndex manage the index route
func (api *API) HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmplt, err := template.ParseFiles(filepath.Join("tmpl", "index.html"), baseTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	tmplt.ExecuteTemplate(w, "base", nil)
	// tmplt.Execute(w, AllPathsWithoutIndex)
}

// HandleAPIIndex manage the api base route
func (api *API) HandleAPIIndex(w http.ResponseWriter, r *http.Request) {
	tmplt, err := template.ParseFiles(filepath.Join("tmpl", "api_index.html"), baseTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	tmplt.ExecuteTemplate(w, "base", AllPathsWithoutIndex)
	// tmplt.Execute(w, AllPathsWithoutIndex)
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
	sortField, sortValue := validSortField(r.URL.Query().Get("order"))
	options := db.Options{Page: page, SortField: sortField, SortValue: sortValue}
	result, err := db.LoadAnimes(api.DB, options)
	if len(result.Animes) == 0 || err != nil {
		InternalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, AnimesResponse{
		Count:    result.Count,
		Results:  result.Animes,
		Next:     assembleAnimesPageLink(result, true, options),
		Previous: assembleAnimesPageLink(result, false, options)}, http.StatusOK)
}

// HandleAnimeDetails manage the anime details endpoint
func (api *API) HandleAnimeDetails(w http.ResponseWriter, r *http.Request) {
	anime, err := getAnime(r, api.DB)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	JSONResponse(w, anime, http.StatusOK)
}

// HandleLatestEpisodes manage the latest episodes endpoint
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

// HandleEpisodeList manage the episodes endpoint
func (api *API) HandleEpisodeList(w http.ResponseWriter, r *http.Request) {
	anime, err := getAnime(r, api.DB)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	JSONResponse(w, EpisodesResponse{AnimeID: anime.Flvid, AnimeName: anime.Name,
		AnimeURL: anime.URL, Episodes: anime.Episodes}, http.StatusOK)
}

// HandleEpisodeDetails manage the episodes endpoint
func (api *API) HandleEpisodeDetails(w http.ResponseWriter, r *http.Request) {
	episodeRes, err := getEpisode(r, api.DB)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	JSONResponse(w, episodeRes, http.StatusOK)
}
