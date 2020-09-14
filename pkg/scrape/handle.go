package scrape

// HandleAnimeScrape handle the response by scraping the anime data
func HandleAnimeScrape(result RequestResult, container *AnimeSPContainer) interface{} {
	if !result.OK {
		return Anime{}
	}
	return GetAnime(result.Document, &container.States, &container.Types, &container.Genres)
}

// HandleEpisodeScrape handle the response by scraping the episode page
func HandleEpisodeScrape(result RequestResult) interface{} {
	if !result.OK {
		return ""
	}
	return GetAnimeURLByEpisodeURL(result.Document)
}
