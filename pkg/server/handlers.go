package server

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kuronosu/animeflv-api/pkg/db"
	"github.com/kuronosu/animeflv-api/pkg/scrape"
	"go.mongodb.org/mongo-driver/mongo"
)

func (api *API) genericDetails(w http.ResponseWriter, r *http.Request,
	dataHandler db.FunctionDataHandler, urlVarID string) {
	id, err := strconv.Atoi(mux.Vars(r)[urlVarID])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	_type, err := dataHandler(id)
	if err != nil {
		http.NotFound(w, r)
	}
	JSONResponse(w, _type, http.StatusOK)
}

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
	types, _ := db.LoadTypes(api.DBManager.Client)
	if len(types) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, types, http.StatusOK)
}

// HandleTypeDetails manage the type details endpoint
func (api *API) HandleTypeDetails(w http.ResponseWriter, r *http.Request) {
	api.genericDetails(w, r, api.DBManager.LoadOneType, "id")
}

// HandleStates manage the states endpoint
func (api *API) HandleStates(w http.ResponseWriter, r *http.Request) {
	states, _ := db.LoadStates(api.DBManager.Client)
	if len(states) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, states, http.StatusOK)
}

// HandleGenres manage the generes endpoint
func (api *API) HandleGenres(w http.ResponseWriter, r *http.Request) {
	genres, _ := db.LoadGenres(api.DBManager.Client)
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
	result, err := db.LoadAnimes(api.DBManager.Client, options)
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
	anime, err := getAnime(r, api.DBManager.Client)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	JSONResponse(w, anime, http.StatusOK)
}

// HandleLatestEpisodes manage the latest episodes endpoint
func (api *API) HandleLatestEpisodes(w http.ResponseWriter, r *http.Request) {
	latestEpisodes, _ := db.LoadLatestEpisodes(api.DBManager.Client)
	if len(latestEpisodes) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	JSONResponse(w, latestEpisodes, http.StatusOK)
}

// HandleDirectory manage the directory endpoint
func (api *API) HandleDirectory(w http.ResponseWriter, r *http.Request) {
	types, _ := db.LoadTypes(api.DBManager.Client)
	if len(types) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	states, _ := db.LoadStates(api.DBManager.Client)
	if len(states) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	genres, _ := db.LoadGenres(api.DBManager.Client)
	if len(genres) == 0 {
		InternalError(w, "Error al cargar datos")
		return
	}
	animes, _ := db.LoadAllAnimes(api.DBManager.Client)
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
	anime, err := getAnime(r, api.DBManager.Client)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	JSONResponse(w, EpisodesResponse{AnimeID: anime.Flvid, AnimeName: anime.Name,
		AnimeURL: anime.URL, Episodes: anime.Episodes}, http.StatusOK)
}

// HandleEpisodeDetails manage the episodes endpoint
func (api *API) HandleEpisodeDetails(w http.ResponseWriter, r *http.Request) {
	episodeRes, err := getEpisode(r, api.DBManager.Client)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	JSONResponse(w, episodeRes, http.StatusOK)
}

// HandleEpisodeVideo manage the videos endpoint
func (api *API) HandleEpisodeVideo(w http.ResponseWriter, r *http.Request) {
	server, found := mux.Vars(r)["server"]
	if !found {
		http.NotFound(w, r)
		return
	}
	server = strings.ToLower(server)
	if !scrape.ValidServer(server) {
		http.NotFound(w, r)
		return
	}
	lang, found := mux.Vars(r)["lang"]
	if !found {
		lang = "SUB"
	}
	lang = strings.ToUpper(lang)
	if !scrape.ValidLang(lang) {
		http.NotFound(w, r)
		return
	}

	episodeRes, err := getEpisode(r, api.DBManager.Client)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	v, err := scrape.GetVideoByURL(scrape.EpisodeURL(episodeRes.Episode.URL))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = v.Active(server, lang)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	JSONResponse(w, v, http.StatusOK)
}

// HandleAnimeSearch manage anime search endpoint
func (api *API) HandleAnimeSearch(w http.ResponseWriter, r *http.Request) {
	name, found := r.URL.Query()["name"]
	if found && len(name) > 0 && name[0] != "" {
		if animes, err := api.DBManager.SearchAnimeByName(name[0]); err == nil {
			JSONResponse(w, animes, http.StatusOK)
			return
		}
	}
	JSONResponse(w, []interface{}{}, http.StatusOK)
}

// Images

func handleImage(w http.ResponseWriter, r *http.Request, url string) {
	reqImg, err := scrape.Fetch(url)
	defer reqImg.Body.Close()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if reqImg.StatusCode != 200 {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Length", fmt.Sprint(reqImg.ContentLength))
	w.Header().Set("Content-Type", reqImg.Header.Get("Content-Type"))
	if _, err = io.Copy(w, reqImg.Body); err != nil {
		http.NotFound(w, r)
		return
	}
}

// HandleScreenshots manage screenshot request
func HandleScreenshots(w http.ResponseWriter, r *http.Request) {
	handleImage(w, r, "https://cdn.animeflv.net"+html.EscapeString(r.URL.Path))
}
