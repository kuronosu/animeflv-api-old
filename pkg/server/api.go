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
	a := &Api{}
	r := mux.NewRouter()
	// r.HandleFunc("/", a.fetchGopher).Methods(http.MethodGet)
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("%s %s %s\n", r.URL, r.Method, time.Now().Format("2006-01-02 15:04:05"))
			h.ServeHTTP(w, r)
		})
	})
	r.HandleFunc("/types/", a.FetchTypes).Methods(http.MethodGet)
	r.HandleFunc("/states/", a.FetchStates).Methods(http.MethodGet)
	r.HandleFunc("/genres/", a.FetchGenres).Methods(http.MethodGet)
	r.HandleFunc("/animes/", a.FetchAnimes).Methods(http.MethodGet)
	r.HandleFunc("/latest/", a.FetchLatestEpisodes).Methods(http.MethodGet)
	r.HandleFunc("/directory/", a.FetchDirectory).Methods(http.MethodGet)
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
