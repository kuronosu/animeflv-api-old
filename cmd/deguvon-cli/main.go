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
var stopEllipsis bool
var c chan struct{} = make(chan struct{}) // event marker

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

func ellipsis(message string, delay time.Duration) {
	for !stopEllipsis {
		for _, r := range []string{".  ", ".. ", "..."} {
			fmt.Printf("\r%s %s", message, r)
			time.Sleep(delay)
		}
	}
	fmt.Fprint(os.Stdout, "\r \r")
	c <- struct{}{}
}

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
	containerI, _, _ := scrape.AllAnimesByPage()
	container := containerI.(scrape.AnimeSPContainer)

	dillDbTime := time.Now()
	db.InsertStates(client, container.States)
	db.InsertTypes(client, container.Types)
	db.InsertGenres(client, container.Genres)
	insertResult, err := db.InsertAnimes(client, container.Animes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Base de datos llenada en %s con %d animes\n", time.Since(dillDbTime), len(insertResult.InsertedIDs))
}

func intervalForLatestEpisodes() {
	fmt.Print("Connect to db")
	client, err := db.SetUp()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\rConnected    ")
	for {
		stopSpinner = false
		go spinner("Getting latest episodes", 100*time.Millisecond)
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
				states, _ := db.LoadStates(client)
				genres, _ := db.LoadGenres(client)
				types, _ := db.LoadTypes(client)
				container := scrape.AnimeSPContainer{
					States: states, Types: types,
					Genres: genres, Animes: []scrape.Anime{}}
				scrape.GetAnimes(relatedURLs, &container)
				db.UpdateOrInsertAnimes(client, container.Animes)
			}
		}
		stopSpinner = true
		<-c // wait spinner stop
		stopEllipsis = false
		fmt.Print("                         ")
		go ellipsis("Waiting", 333*time.Millisecond)
		time.Sleep(1 * time.Minute)
		stopEllipsis = true
		<-c // wait ellipsis stop
	}
}
