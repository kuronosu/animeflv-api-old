package scrape

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetDirectoryPageCount return the count of page in the animeflv directory
func GetDirectoryPageCount() (int, error) {
	doc, err := FetchDocument(DirectoryURL)
	if err != nil {
		return 0, err
	}
	count := doc.Find("ul.pagination").First().Find("a").Eq(-2).Text()
	return strconv.Atoi(count)
}

// GetAnimeURLSFromDirectory as
func GetAnimeURLSFromDirectory() ([]string, error) {
	pages, err := GetDirectoryPageCount()
	if err != nil {
		return nil, err
	}
	var urls []string
	for _, page := range MakeRange(1, pages) {
		tmpUrls, err := GetAnimeURLSFromDirectoryPage(page)
		if err == nil {
			urls = append(urls, tmpUrls...)
		}
	}
	return urls, nil
}

// GetAnimeURLSFromDirectoryPage return all anime urls on a page
func GetAnimeURLSFromDirectoryPage(page int) ([]string, error) {
	// fmt.Println("getting urls from page ", page)
	doc, err := FetchDocument(DirectoryURLPage(page))
	if err != nil {
		return nil, err
	}
	var urls []string
	doc.Find("article.Anime").Each(func(index int, element *goquery.Selection) {
		val, exists := element.Find("a").First().Attr("href")
		if exists {
			urls = append(urls, val)
		}
	})
	return urls, nil
}

//GetAnime scrape anime data from document
func GetAnime(doc *goquery.Document, states *[]State, types *[]Type, genres *[]Genre) Anime {
	anime := Anime{
		Type:       getType(doc, types),
		State:      getState(doc, states),
		Genres:     getGenres(doc, genres),
		OtherNames: getOtherNames(doc),
		Synopsis:   getSynopsis(doc),
		Score:      getScore(doc),
		Votes:      getVotes(doc),
		Cover:      getCover(doc),
		Banner:     getBanner(doc),
		Relations:  getRelations(doc),
	}
	jsContent := getScript(doc)
	setAnimeDataFromScript(&anime, jsContent)
	setEpisodesFromScript(&anime, jsContent)
	return anime
}

func getScript(document *goquery.Document) string {
	script := document.Find("script").FilterFunction(func(i int, sc *goquery.Selection) bool {
		return strings.Contains(sc.Text(), "var anime_info = ")
	})

	return script.Text()
}

func getType(document *goquery.Document, types *[]Type) int {
	typeString := strings.Trim(document.Find("span.Type").Text(), " ")
	for _, t := range *types {
		if t.Name == typeString {
			return t.ID
		}
	}
	_type := Type{ID: len(*types), Name: typeString}
	*types = append(*types, _type)
	return _type.ID
}

func getState(document *goquery.Document, states *[]State) int {
	stateString := strings.Trim(document.Find("span.fa-tv").Text(), " ")
	for _, s := range *states {
		if s.Name == stateString {
			return s.ID
		}
	}
	state := State{ID: len(*states), Name: stateString}
	*states = append(*states, state)
	return state.ID
}

func getGenres(document *goquery.Document, genres *[]Genre) []int {
	var genresIDs []int
	document.Find("nav.Nvgnrs").Find("a").Each(func(index int, element *goquery.Selection) {
		genreString := strings.Trim(element.Text(), " ")
		for _, g := range *genres {
			if g.Name == genreString {
				genresIDs = append(genresIDs, g.ID)
				return
			}
		}
		genre := Genre{ID: len(*genres), Name: genreString}
		*genres = append(*genres, genre)
		genresIDs = append(genresIDs, genre.ID)
	})
	return genresIDs
}

func getOtherNames(document *goquery.Document) []string {
	var names []string
	document.Find("span.TxtAlt").Each(func(index int, element *goquery.Selection) {
		name := strings.Trim(element.Text(), " ")
		if cfEmail := element.Find(".__cf_email__"); cfEmail.Length() > 0 {
			decoded := DecodeEmail(cfEmail.AttrOr("data-cfemail", ""))
			name = strings.ReplaceAll(name, EmailProtected, decoded)
		}
		names = append(names, name)
	})
	return names
}

func getScore(document *goquery.Document) float32 {
	str := strings.Trim(document.Find("span#votes_prmd").Text(), " ")
	score, _ := strconv.ParseFloat(str, 32)
	return float32(score)
}

func getVotes(document *goquery.Document) int {
	str := strings.Trim(document.Find("span#votes_nmbr").Text(), " ")
	score, _ := strconv.Atoi(str)
	return score
}

func getCover(document *goquery.Document) string {
	return strings.Trim(document.Find("div.AnimeCover").Find("div.Image").Find("img").AttrOr("style", ""), " ")
}

func getBanner(document *goquery.Document) string {
	s := strings.Trim(strings.ReplaceAll(document.Find("div.Bg").AttrOr("style", ""), "background-image:url(", ""), " ")
	if len(s) > 0 {
		return s[0 : len(s)-1]
	}
	return ""
}

func getSynopsis(document *goquery.Document) string {
	var description string
	document.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); name == "description" {
			description = s.AttrOr("content", "")
		}
	})
	return description
}

func getRelations(document *goquery.Document) []Relation {
	relations := []Relation{}
	document.Find("ul.ListAnmRel").Find("li.fa-chevron-circle-right").Each(func(i int, el *goquery.Selection) {
		a := el.Find("a")
		name := a.Text()
		url := a.AttrOr("href", "")
		if cfEmail := a.Find(".__cf_email__"); cfEmail.Length() > 0 {
			decoded := DecodeEmail(cfEmail.AttrOr("data-cfemail", ""))
			name = strings.ReplaceAll(name, EmailProtected, decoded)
		}
		rel := strings.Trim(el.Contents().Not("a").Text(), " ")
		if len(rel) >= 2 {
			rel = rel[1 : len(rel)-1]
		}
		relations = append(relations, NewRelation(name, url, rel))
	})
	return relations
}

// setAnimeDataFromScript set Flvid Name Slug NextEpisodeDate from script
func setAnimeDataFromScript(a *Anime, script string) {
	rawAnimeInfo := AnimeScriptPattern.FindString(script)
	rawAnimeInfo = strings.ReplaceAll(rawAnimeInfo, "var anime_info = ", "")
	rawAnimeInfo = rawAnimeInfo[0 : len(rawAnimeInfo)-1]
	var animeInfo []string
	_ = json.Unmarshal([]byte(rawAnimeInfo), &animeInfo)
	if len(animeInfo) > 0 && len(animeInfo) <= 4 {
		a.Flvid = animeInfo[0]
		a.Name = animeInfo[1]
		a.Slug = animeInfo[2]
		a.URL = "/anime/" + a.Slug
	}
	if len(animeInfo) == 4 {
		a.NextEpisodeDate = animeInfo[3]
	}
}

// setEpisodesFromScript set Flvid Name Slug NextEpisodeDate from script
func setEpisodesFromScript(a *Anime, script string) {
	rawEpisodes := EpisodeScriptPattern.FindString(script)
	rawEpisodes = strings.ReplaceAll(rawEpisodes, "var episodes = ", "")
	rawEpisodes = rawEpisodes[0 : len(rawEpisodes)-1]
	var episodesData [][]float32
	_ = json.Unmarshal([]byte(rawEpisodes), &episodesData)
	var episodes []Episode
	for _, ep := range episodesData {
		episodes = append(episodes, NewEpisode(ep[0], int(ep[1]), a.Flvid, a.Slug))
	}
	sort.SliceStable(episodes, func(i, j int) bool {
		return episodes[i].Number < episodes[j].Number
	})
	a.Episodes = episodes
}
