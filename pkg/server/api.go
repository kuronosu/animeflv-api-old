package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

// New create new server
func New(client *mongo.Client) Server {
	setDbClient(client)
	setUpTemplatePath()
	a := &Api{}
	r := mux.NewRouter()
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("%s %s %s\n", r.URL, r.Method, time.Now().Format("2006-01-02 15:04:05"))
			h.ServeHTTP(w, r)
		})
	})
	r.HandleFunc(IndexPath, HandleIndex).Methods(http.MethodGet)
	r.HandleFunc(TypesPath, HandleTypes).Methods(http.MethodGet)
	r.HandleFunc(StatesPath, HandleStates).Methods(http.MethodGet)
	r.HandleFunc(GenresPath, HandleGenres).Methods(http.MethodGet)
	r.HandleFunc(AnimesPath, HandleAnimes).Methods(http.MethodGet)
	r.HandleFunc(DirectoryPath, HandleDirectory).Methods(http.MethodGet)
	r.HandleFunc(LatestEpisodesPath, HandleLatestEpisodes).Methods(http.MethodGet)
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
