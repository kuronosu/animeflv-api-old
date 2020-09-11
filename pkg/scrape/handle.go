package scrape

import (
	"github.com/PuerkitoBio/goquery"
)

// HandleAnimeScrape handle the response by scraping the anime data
func HandleAnimeScrape(httpResp *HTTPResponse, doc *goquery.Document) interface{} {
	if httpResp.Err != nil {
		return Anime{}
	}
	return GetAnime(doc)
}
