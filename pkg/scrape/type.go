package scrape

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Relation represents raw relation between 2 animes
type Relation struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Rel  string `json:"rel"`
}

// Episode represents raw data of episode from animeflv
type Episode struct {
	Number float64 `json:"number"`
	Eid    int     `json:"eid"`
	URL    string  `json:"url"`
	Img    string  `json:"img"`
}

// Genre represents series genre
type Genre struct {
	ID   int    `bson:"_id" json:"id"`
	Name string `json:"name"`
}

// Type represents series type
type Type struct {
	ID   int    `bson:"_id" json:"id"`
	Name string `json:"name"`
}

// State represents anime state
type State struct {
	ID   int    `bson:"_id" json:"id"`
	Name string `json:"name"`
}

// Anime represents raw data of anime from animeflv
type Anime struct {
	// anime_info var of script
	Flvid           int    `bson:"_id" json:"flvid"` //OK
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	NextEpisodeDate string `json:"nextEpisodeDate"`
	// Other anime info
	URL        string   `json:"url"`
	State      int      `json:"state"`
	Type       int      `json:"type"`
	Genres     []int    `json:"genres"`
	OtherNames []string `json:"otherNames"`
	Synopsis   string   `json:"synopsis"`
	Score      float64  `json:"score"`
	Votes      int      `json:"votes"`
	// Images
	Cover  string `json:"cover"`
	Banner string `json:"banner"`
	// Relations
	Relations []Relation `json:"relations"`
	Episodes  []Episode  `json:"episodes"`
}

// RequestResult contains the response, document and the processed data
type RequestResult struct {
	URL                   string
	Response              *http.Response
	ResponseErr           error
	Document              *goquery.Document
	DocumentErr           error
	ProcessedResponseData interface{}
	OK                    bool
}

// AnimeSPContainer contains states types genres array pointers
type AnimeSPContainer struct {
	States []State
	Types  []Type
	Genres []Genre
	Animes []Anime
}

// Directory contains all directory data
type Directory struct {
	States []State       `json:"states"`
	Types  []Type        `json:"types"`
	Genres []Genre       `json:"genres"`
	Animes map[int]Anime `json:"animes"`
}

// LatestEpisode represent the info of the latest episode
type LatestEpisode struct {
	URL   string `json:"url"`
	Image string `json:"image"`
	Capi  string `json:"capi"`
	Anime int    `json:"anime"`
}

// NewEpisode create a episode instance
func NewEpisode(Number float64, Eid int, Flvid int, animeSlug string) Episode {
	return Episode{
		Number: Number,
		Eid:    Eid,
		URL:    fmt.Sprintf("/ver/%s-%s", animeSlug, FloatToString(Number)),
		Img:    fmt.Sprintf("https://cdn.animeflv.net/screenshots/%d/%s/th_3.jpg", Flvid, FloatToString(Number)),
	}
}

// NewRelation create a relation instance
func NewRelation(Name string, URL string, Rel string) Relation {
	return Relation{
		Name: Name,
		URL:  URL,
		Rel:  Rel,
	}
}

// NewAnime create a anime instance
// func NewAnime(Flvid string, Name string, Slug string, NextEpisodeDate string, URL string, State string, typea string, Genres []string, Synopsis string, Score float32, Votes int, Cover string, Banner string, Relations []Relation, Episodes []Episode, OtherNames []string) Anime {
// 	return Anime{
// 		Flvid:           Flvid,
// 		Name:            Name,
// 		Slug:            Slug,
// 		NextEpisodeDate: NextEpisodeDate,
// 		URL:             URL,
// 		State:           State,
// 		Typea:           typea,
// 		Genres:          Genres,
// 		Synopsis:        Synopsis,
// 		Score:           Score,
// 		Votes:           Votes,
// 		Cover:           Cover,
// 		Banner:          Banner,
// 		Relations:       Relations,
// 		Episodes:        Episodes,
// 		OtherNames:      OtherNames,
// 	}
// }

// VideoServerData contains the data of a video server
type VideoServerData struct {
	Server      string `json:"server"`
	Title       string `json:"title"`
	AllowMobile bool   `json:"allow_mobile"`
	Code        string `json:"code"`
}

type Video struct {
	ActiveURL string
	Servers   map[string][]VideoServerData
}
