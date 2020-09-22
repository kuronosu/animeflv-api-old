package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// Api represents the api
type Api struct {
	router http.Handler
}

// Server represents the api
type Server interface {
	Router() http.Handler
}

func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wc := NewResponseWriterCounter(w)
		h.ServeHTTP(wc, r)
		fmt.Printf("%s %s %s %d %d\n",
			wc.Started().Format("2006-01-02 15:04:05"),
			r.URL,
			r.Method,
			wc.StatusCode(),
			wc.Count())
	})
}

// New create new server
func New(client *mongo.Client) Server {
	setDbClient(client)
	a := &Api{}
	r := mux.NewRouter()
	r.Use(LogMiddleware)
	r.HandleFunc(IndexPath, HandleIndex).Methods(http.MethodGet)
	r.HandleFunc(TypesPath, HandleTypes).Methods(http.MethodGet)
	r.HandleFunc(StatesPath, HandleStates).Methods(http.MethodGet)
	r.HandleFunc(GenresPath, HandleGenres).Methods(http.MethodGet)
	r.HandleFunc(AnimesPath, HandleAnimes).Methods(http.MethodGet)
	r.HandleFunc(DirectoryPath, HandleDirectory).Methods(http.MethodGet)
	r.HandleFunc(LatestEpisodesPath, HandleLatestEpisodes).Methods(http.MethodGet)
	r.HandleFunc(AnimeDetailsPath, HandleAnimeDetails).Methods(http.MethodGet)
	a.router = r
	return a
}

// Router return the api router
func (a *Api) Router() http.Handler {
	return CaselessMatcher(a.router)
}

func CaselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		if !strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path = r.URL.Path + "/"
		}
		next.ServeHTTP(w, r)
	})
}
