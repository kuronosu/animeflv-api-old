package scrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// SUB Subtitled key
const SUB = "SUB"

// LAT Latino key
const LAT = "LAT"

// ESP Spanish key
const ESP = "ESP"

// Langs contains all the possible lang of a video
var Langs = map[string]string{SUB: "Subtitulado", LAT: "Latino", ESP: "Español"}

func getRawVideosData(doc *goquery.Document) (string, error) {
	script := doc.Find("script").FilterFunction(func(i int, sc *goquery.Selection) bool {
		text := sc.Text()
		return strings.Contains(text, "var videos") &&
			strings.Contains(text, "var anime_id") &&
			strings.Contains(text, "var episode_id") &&
			strings.Contains(text, "var episode_number")
	}).Text()
	script = strings.ReplaceAll(ServersScriptPattern.FindString(script), "var videos = ", "")
	if script == "" {
		return "nil", fmt.Errorf("Could not find the video data in the document")
	}
	return script[:len(script)-1], nil
}

// GetVideoByURL get the video data from episode url
func GetVideoByURL(url string) (*Video, error) {
	if !strings.Contains(url, AnimeFlvURL+`/ver/`) {
		return nil, fmt.Errorf("The url '%s' is not valid episode url", url)
	}
	doc, err := FetchDocument(url)
	if err != nil {
		return nil, err
	}
	rawVideos, err := getRawVideosData(doc)
	if err != nil {
		return nil, err
	}
	var servers map[string][]VideoServerData
	err = json.Unmarshal([]byte(rawVideos), &servers)
	if err != nil {
		return nil, err
	}
	return &Video{Servers: servers}, nil
}

// Active get the active url from a server
func (v *Video) Active(server, lang string) error {
	switch strings.ToLower(server) {
	case "gocdn":
		return v.Gocdn(lang)
	}
	return fmt.Errorf("%s is not a valid server", server)
}

// Gocdn activate the gocdn video server
func (v *Video) Gocdn(lang string) error {
	if !ValidLang(lang) {
		return fmt.Errorf("%s is not a valid lang", lang)
	}
	// if v.Servers.
	if servers, found := v.Servers[lang]; found {
		for _, video := range servers {
			if video.Server == "gocdn" {
				subs := strings.Split(video.Code, "#")
				if len(subs) > 0 {
					resp, err := Fetch("https://streamium.xyz/gocdn.php?v=" + subs[len(subs)-1])
					if err != nil {
						return err
					}
					bodyBytes, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return err
					}
					var data map[string]string
					json.Unmarshal(bodyBytes, &data)
					url, found := data["file"]
					if !found {
						return fmt.Errorf("could not get the url of the video")
					}
					v.ActiveURL = url
					return nil
				}
			}
		}
	}
	return fmt.Errorf("gocdn is not available in the language %s", lang)
}
