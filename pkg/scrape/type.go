package scrape

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Relation represents raw relation between 2 animes
type Relation struct {
	Name string //OK
	URL  string //OK
	Rel  string //OK
}

// Episode represents raw data of episode from animeflv
type Episode struct {
	Number float64
	Eid    int
	URL    string
	Img    string
}

// Genre represents series genre
type Genre struct {
	ID   int `bson:"_id"`
	Name string
}

// Type represents series type
type Type struct {
	ID   int `bson:"_id"`
	Name string
}

// State represents anime state
type State struct {
	ID   int `bson:"_id"`
	Name string
}

// Anime represents raw data of anime from animeflv
type Anime struct {
	// anime_info var of script
	Flvid           int    `bson:"_id"` //OK
	Name            string //OK
	Slug            string //OK
	NextEpisodeDate string //OK
	// Other anime info
	URL        string   //OK
	State      int      //OK
	Type       int      //OK
	Genres     []int    //OK
	OtherNames []string //OK
	Synopsis   string   //OK
	Score      float64  //OK
	Votes      int      //OK
	// Images
	Cover  string //OK
	Banner string //OK
	// Relations
	Relations []Relation //OK
	Episodes  []Episode
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
	States []State
	Types  []Type
	Genres []Genre
	Animes map[int]Anime
}

// LatestEpisode represent the info of the latest episode
type LatestEpisode struct {
	URL   string
	Image string
	Capi  string
	Anime int
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
