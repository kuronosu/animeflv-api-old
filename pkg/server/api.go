package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kuronosu/animeflv-api/pkg/db"
)

const static = "/static/"

// New create new server
func New(manager db.Manager, port int) Server {
	a := &API{DBManager: manager, port: port}
	r := mux.NewRouter()
	r.PathPrefix(static).Handler(http.StripPrefix(static, http.FileServer(http.Dir("."+static))))
	r.HandleFunc("/", a.HandleIndex).Methods(http.MethodGet)
	r.HandleFunc(APIPath, a.HandleAPIIndex).Methods(http.MethodGet)
	r.HandleFunc(TypesPath, a.HandleTypes).Methods(http.MethodGet)
	r.HandleFunc(TypeDetailsPath, a.HandleTypeDetails).Methods(http.MethodGet)
	r.HandleFunc(StatesPath, a.HandleStates).Methods(http.MethodGet)
	r.HandleFunc(GenresPath, a.HandleGenres).Methods(http.MethodGet)
	r.HandleFunc(AnimesPath, a.HandleAnimes).Methods(http.MethodGet)
	r.HandleFunc(DirectoryPath, a.HandleDirectory).Methods(http.MethodGet)
	r.HandleFunc(LatestEpisodesPath, a.HandleLatestEpisodes).Methods(http.MethodGet)
	r.HandleFunc(AnimeDetailsPath, a.HandleAnimeDetails).Methods(http.MethodGet)
	r.HandleFunc(EpisodeListPath, a.HandleEpisodeList).Methods(http.MethodGet)
	r.HandleFunc(EpisodeDetailsPath, a.HandleEpisodeDetails).Methods(http.MethodGet)
	r.HandleFunc(VideoPath, a.HandleEpisodeVideo).Methods(http.MethodGet)
	r.HandleFunc(VideoLangPath, a.HandleEpisodeVideo).Methods(http.MethodGet)
	r.HandleFunc(SearchAnimePath, a.HandleAnimeSearch).Methods(http.MethodGet)
	a.router = r
	return a
}
