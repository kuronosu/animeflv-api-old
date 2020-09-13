package scrape

import (
	"fmt"
	"net/http"
	"os"
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

// AsyncHTTPGetsAnimes make request to multiples Urls asynchronously
func AsyncHTTPGetsAnimes(urls []string, container *AnimeSPContainer) []*RequestResult {
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
			if result.OK = result.ResponseErr == nil; result.OK {
				result.Document, result.DocumentErr = GetDocument(result.Response, result.URL)
				result.OK = result.OK && result.DocumentErr == nil
				result.ProcessedResponseData = HandleAnimeScrape(*result, container)
			}
			results = append(results, result)
			if len(results) == len(urls) {
				return results
			}
		}
	}
}

// AllAnimesByPage get all the animes by making asynchronous requests one page at a time
func AllAnimesByPage() (interface{}, []RequestResult, []int) {
	start := time.Now()
	errs := []RequestResult{}
	pagesErr := []int{}
	pages, err := GetDirectoryPageCount()
	container := AnimeSPContainer{
		States: []State{},
		Types:  []Type{},
		Genres: []Genre{},
		Animes: []Anime{},
	}
	if err != nil {
		return container, errs, pagesErr
	}
	for _, page := range MakeRange(1, pages) {
		start2 := time.Now()
		urls, err := GetAnimeURLSFromDirectoryPage(page)
		errcp := 0
		animeCount := 0
		if err == nil {
			results := AsyncHTTPGetsAnimes(urls, &container)
			for _, result := range results {
				if !result.OK {
					errcp++
					errs = append(errs, *result)
					continue
				}
				animeCount++
				container.Animes = append(container.Animes, result.ProcessedResponseData.(Anime))
			}
			// time.Sleep(500 * time.Millisecond)
			fmt.Fprint(os.Stdout, fmt.Sprintf("\r \rScraped page %d of %d with %d animes and %d errors in %s. Total animes %d in %s",
				page, pages, animeCount, errcp, time.Since(start2), len(container.Animes), time.Since(start)))
		} else {
			pagesErr = append(pagesErr, page)
		}
	}
	// return errc
	fmt.Fprint(os.Stdout, fmt.Sprintf("\r \rCompleted... Pages: %d Animes: %d, Erros: %d in Time %s                                            \n",
		pages, len(container.Animes), len(errs), time.Since(start)))
	return container, errs, pagesErr
}
