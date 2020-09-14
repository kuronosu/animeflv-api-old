package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kuronosu/deguvon-server-go/pkg/db"
	"github.com/kuronosu/deguvon-server-go/pkg/scrape"
)

var stopSpinner bool

func spinner(message string, delay time.Duration) {
	for !stopSpinner {
		for _, r := range `-\|/` {
			fmt.Printf("\r%s %c", message, r)
			time.Sleep(delay)
		}
	}
	fmt.Fprint(os.Stdout, "\r \r")
	c <- struct{}{}
}

var c chan struct{} = make(chan struct{}) // event marker

func main() {
	createFlag := flag.Bool("c", false, "Create directory")
	latestEpisodesFlag := flag.Bool("le", false, "Latest episodes")
	helpFlag := flag.Bool("h", true, "Help")
	flag.Parse()

	if *createFlag {
		createDirectory()
	} else if *latestEpisodesFlag {
		intervalForLatestEpisodes()
	} else if *helpFlag {
		fmt.Println("-h: Help\n-c: Create directory\n-le: Latest episodes")
	}
}

func createDirectory() {
	client, err := db.SetUp()
	if err != nil {
		log.Fatal(err)
	}
	// stopSpinner = false
	// go spinner("Obteniendo animes", 100*time.Millisecond)
	containerI, _, _ := scrape.AllAnimesByPage()
	container := containerI.(scrape.AnimeSPContainer)
	// stopSpinner = true
	// <-c // wait spinner stop

	dillDbTime := time.Now()
	db.InsertStates(client, container.States)
	db.InsertTypes(client, container.Types)
	db.InsertGenres(client, container.Genres)
	insertResult, err := db.InsertAnimes(client, container.Animes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Base de datos llenada en %s con %d animes\n", time.Since(dillDbTime), len(insertResult.InsertedIDs))

	// client, _ := db.SetUp()
	// s, e := db.GetNextSequence(client, "users2")
	// fmt.Println(s, e)
}

func intervalForLatestEpisodes() {
	client, err := db.SetUp()
	if err != nil {
		log.Fatal(err)
	}
	for {
		le, a, e := scrape.FetchLatestEpisodes()
		if e == nil {
			db.SetLatestEpisodes(client, le)
			_, in, _ := db.UpdateOrInsertAnimes(client, a.Animes)
			if len(in) > 0 {
				relatedURLs := []string{}
				for _, anime := range in {
					for _, rel := range anime.Relations {
						relatedURLs = append(relatedURLs, rel.URL)
					}
				}
				container := scrape.AnimeSPContainer{
					States: []scrape.State{}, Types: []scrape.Type{},
					Genres: []scrape.Genre{}, Animes: []scrape.Anime{}}
				scrape.GetAnimes(relatedURLs, &container)
				db.UpdateOrInsertAnimes(client, container.Animes)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
