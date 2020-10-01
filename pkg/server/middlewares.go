package server

import (
	"log"
	"net/http"
	"strings"
)

// LogMiddleware prints the data of request and response
func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wc := NewResponseWriterCounter(w)
		h.ServeHTTP(wc, r)
		log.Printf("%s %s %d %d\n",
			r.URL,
			r.Method,
			wc.StatusCode(),
			wc.Count())
	})
}

// CaselessMatcher modify the request url path to be case insensitive and add / at the end if it is not already
func CaselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// TrailingSlashes redirect when has end slash
func TrailingSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, trimSuffix(r.URL.Path, "/"), http.StatusSeeOther)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
