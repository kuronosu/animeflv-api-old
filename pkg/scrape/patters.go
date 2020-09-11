package scrape

import (
	"regexp"
)

//AnimeScriptPattern regex pattern for anime_info in script
var AnimeScriptPattern = regexp.MustCompile(`var anime_info = \[.+\];`)

//EpisodeScriptPattern regex pattern for episodes in script
var EpisodeScriptPattern = regexp.MustCompile(`var episodes = \[.*\];`)

//ServersScriptPattern regex pattern for videos in script
var ServersScriptPattern = regexp.MustCompile(`var videos = {.+};`)
