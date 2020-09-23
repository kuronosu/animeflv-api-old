package server

import (
	"fmt"
	"net/http"
	"strings"
)

// LogMiddleware prints the data of request and response
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

func CaselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		if !strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path = r.URL.Path + "/"
		}
		next.ServeHTTP(w, r)
	})
}
