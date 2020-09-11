package scrape

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Fetch make request and return response
func Fetch(URL string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(URL)
	if err != nil {
		resp.Body.Close()
	}
	return resp, err
}

//GetDocument extract the goquery.Document from the response
func GetDocument(resp *http.Response, url string) (*goquery.Document, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("The %s page was not found (404 error)", url)
	}
	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return document, err
}

// FetchDocument return the goquery document from url
func FetchDocument(URL string) (*goquery.Document, error) {
	resp, err := Fetch(URL)
	if err != nil {
		return nil, err
	}
	return GetDocument(resp, URL)
}

// AsyncHTTPGets make request to multiples Urls asynchronously
func AsyncHTTPGets(Urls []string, handler func(*HTTPResponse, *goquery.Document) interface{}) []*Result {
	ch := make(chan *HTTPResponse, len(Urls)) // buffered
	for _, URL := range Urls {
		go func(URL string) {
			resp, err := http.Get(AnimeURL(URL))
			ch <- &HTTPResponse{Response: resp, URL: URL, Err: err}
		}(URL)
	}

	results := []*Result{}
	for {
		select {
		case response := <-ch:
			var hResponse interface{}
			var doc *goquery.Document
			var err error = nil
			// var err error = nil
			if response.Err == nil {
				doc, err = GetDocument(response.Response, response.URL)
				if err == nil {
					hResponse = handler(response, doc)
				}
			}
			res := Result{
				HTTPResponse:    *response,
				HandledResponse: hResponse,
				Document:        Document{Document: doc, Err: err},
			}
			results = append(results, &res)
			if len(results) == len(Urls) {
				return results
			}
		}
	}
}
