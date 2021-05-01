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

// LAT Latin spanish key
const LAT = "LAT"

// ESP Spain spanish key
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
	if !ValidLang(lang) {
		return fmt.Errorf("%s is not a valid lang", lang)
	}
	switch strings.ToLower(server) {
	case "gocdn":
		return v.Gocdn(lang)
	case "fembed":
		return v.Fembed(lang)
	}
	return fmt.Errorf("%s is not a valid server", server)
}

// NotAvailableLanguage error
func NotAvailableLanguage(server string, lang string) error {
	return fmt.Errorf("%s is not available in the language %s", server, lang)
}

// CheckLanguageByServer check if a video server supports a language
func (v *Video) CheckLanguageByServer(server string, lang string) (VideoServerData, error) {
	if servers, found := v.Servers[lang]; found {
		for _, video := range servers {
			if video.Server == server {
				return video, nil
			}
		}
	}
	return VideoServerData{}, NotAvailableLanguage("Gocdn", lang)
}

// Gocdn activate the gocdn video server
func (v *Video) Gocdn(lang string) error {
	video, err := v.CheckLanguageByServer("gocdn", lang)
	if err != nil {
		return err
	}
	subs := strings.Split(video.Code, "#")
	if len(subs) == 0 {
		return fmt.Errorf("Cant extract the code of video")
	}
	resp, err := Fetch("https://streamium.xyz/gocdn.php?v=" + subs[len(subs)-1])
	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var data map[string]string
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return err
	}
	url, found := data["file"]
	if !found {
		return fmt.Errorf("could not get the url of the video")
	}
	v.ActiveURL = url
	return nil
}

// Fembed activate the fembed video server
func (v *Video) Fembed(lang string) error {
	video, err := v.CheckLanguageByServer("fembed", lang)
	if err != nil {
		return err
	}
	resp, err := FetchPost(strings.Replace(video.Code, "/v/", "/api/source/", 1))
	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var fdata FembedResponse
	err = json.Unmarshal(bodyBytes, &fdata)
	if !fdata.Success || len(fdata.Data) == 0 || err != nil {
		return fmt.Errorf("Request was not succeeded")
	}
	v.ActiveURL = fdata.Data[0].File
	return nil
}
