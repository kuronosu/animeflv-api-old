package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kuronosu/animeflv-api/pkg/db"
)

// InternalError make an json response with error message
func InternalError(w http.ResponseWriter, err string) {
	JSONResponse(w, ErrorResponse{err, http.StatusInternalServerError}, http.StatusInternalServerError)
}

// JSONResponse create an http response in json format
func JSONResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func assembleAnimesPageLink(result db.PaginatedAnimeResult, next bool) *string {
	newURI := AnimesPath + "?page=%d"
	if next && result.Page < result.TotalPages {
		newURI = fmt.Sprintf(newURI, result.Page+1)
	} else if !next && result.Page >= 2 {
		newURI = fmt.Sprintf(newURI, result.Page-1)
	} else {
		return nil
	}
	return &newURI
}
