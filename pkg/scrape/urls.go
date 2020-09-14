package scrape

import "fmt"

// AnimeFlvURL url for animeflv
const AnimeFlvURL = "https://animeflv.net"

// DirectoryURL url for animeflv directory
const DirectoryURL = AnimeFlvURL + "/browse"

// DirectoryURLPage url for animeflv directory with page parameter
func DirectoryURLPage(page int) string {
	return fmt.Sprintf(DirectoryURL+"?page=%d", page)
}

// AnimeURL url for specific anime
func AnimeURL(path string) string {
	return fmt.Sprintf(AnimeFlvURL+"%s", path)
}

// EpisodeURL url for specific episode
func EpisodeURL(path string) string {
	return fmt.Sprintf(AnimeFlvURL+"%s", path)
}

// UserAgent to make requests
const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36 Edg/85.0.564.44"
