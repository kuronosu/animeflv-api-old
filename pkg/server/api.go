package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// New create new server
func New(client *mongo.Client) Server {
	a := &API{DB: client}
	r := mux.NewRouter()
	r.Use(LogMiddleware)
	r.HandleFunc(IndexPath, a.HandleIndex).Methods(http.MethodGet)
	r.HandleFunc(TypesPath, a.HandleTypes).Methods(http.MethodGet)
	r.HandleFunc(StatesPath, a.HandleStates).Methods(http.MethodGet)
	r.HandleFunc(GenresPath, a.HandleGenres).Methods(http.MethodGet)
	r.HandleFunc(AnimesPath, a.HandleAnimes).Methods(http.MethodGet)
	r.HandleFunc(DirectoryPath, a.HandleDirectory).Methods(http.MethodGet)
	r.HandleFunc(LatestEpisodesPath, a.HandleLatestEpisodes).Methods(http.MethodGet)
	r.HandleFunc(AnimeDetailsPath, a.HandleAnimeDetails).Methods(http.MethodGet)
	a.router = r
	return a
}
