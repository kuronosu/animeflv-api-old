package scrape

// HandleAnimeScrape handle the response by scraping the anime data
func HandleAnimeScrape(result RequestResult) interface{} {
	if result.ResponseErr != nil || result.DocumentErr != nil {
		return Anime{}
	}
	return GetAnime(result.Document)
}
