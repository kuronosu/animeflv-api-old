package scrape

// HandleAnimeScrape handle the response by scraping the anime data
func HandleAnimeScrape(result RequestResult, container *AnimeSPContainer) interface{} {
	if !result.OK {
		return Anime{}
	}
	return GetAnime(result.Document, container)
}
