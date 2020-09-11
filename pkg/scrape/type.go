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
	Number float32
	Eid    int
	URL    string
	Img    string
}

// Anime represents raw data of anime from animeflv
type Anime struct {
	// anime_info var of script
	Flvid           string //OK
	Name            string //OK
	Slug            string //OK
	NextEpisodeDate string //OK
	// Other anime info
	URL        string   //OK
	State      string   //OK
	Typea      string   //OK
	Genres     []string //OK
	OtherNames []string //OK
	Synopsis   string   //OK
	Score      float32  //OK
	Votes      int      //OK
	// Images
	Cover  string //OK
	Banner string //OK
	// Relations
	Relations []Relation //OK
	Episodes  []Episode
}

// HTTPResponse represents http Response
type HTTPResponse struct {
	URL      string
	Response *http.Response
	Err      error
}

// Document content
type Document struct {
	Document *goquery.Document
	Err      error
}

// Result contains the response and the handled response
type Result struct {
	HTTPResponse    HTTPResponse
	HandledResponse interface{}
	Document        Document
}

// NewEpisode create a episode instance
func NewEpisode(Number float32, Eid int, Flvid string, animeSlug string) Episode {
	return Episode{
		Number: Number,
		Eid:    Eid,
		URL:    fmt.Sprintf("/ver/%s-%s", animeSlug, FloatToString(Number)),
		Img:    fmt.Sprintf("https://cdn.animeflv.net/screenshots/%s/%s/th_3.jpg", Flvid, FloatToString(Number)),
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
func NewAnime(Flvid string, Name string, Slug string, NextEpisodeDate string, URL string, State string, typea string, Genres []string, Synopsis string, Score float32, Votes int, Cover string, Banner string, Relations []Relation, Episodes []Episode, OtherNames []string) Anime {
	return Anime{
		Flvid:           Flvid,
		Name:            Name,
		Slug:            Slug,
		NextEpisodeDate: NextEpisodeDate,
		URL:             URL,
		State:           State,
		Typea:           typea,
		Genres:          Genres,
		Synopsis:        Synopsis,
		Score:           Score,
		Votes:           Votes,
		Cover:           Cover,
		Banner:          Banner,
		Relations:       Relations,
		Episodes:        Episodes,
		OtherNames:      OtherNames,
	}
}
