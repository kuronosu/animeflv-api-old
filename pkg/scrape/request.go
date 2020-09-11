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
func AsyncHTTPGets(urls []string, handler func(RequestResult) interface{}) []*RequestResult {
	ch := make(chan *RequestResult, len(urls)) // buffered
	for _, URL := range urls {
		go func(URL string) {
			resp, err := http.Get(AnimeURL(URL))
			ch <- &RequestResult{Response: resp, URL: URL, ResponseErr: err}
		}(URL)
	}

	results := []*RequestResult{}
	for {
		select {
		case result := <-ch:
			if result.ResponseErr == nil {
				result.Document, result.DocumentErr = GetDocument(result.Response, result.URL)
				result.ProcessedResponseData = handler(*result)
			}
			result.OK = result.ResponseErr != nil && result.DocumentErr != nil
			results = append(results, result)
			if len(results) == len(urls) {
				return results
			}
		}
	}
}
