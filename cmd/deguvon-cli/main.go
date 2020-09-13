package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kuronosu/deguvon/pkg/db"
	"github.com/kuronosu/deguvon/pkg/scrape"
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
	states := make([]interface{}, len(container.States))
	for i, v := range container.States {
		states[i] = v
	}
	db.InsertStates(client, states)
	types := make([]interface{}, len(container.Types))
	for i, v := range container.Types {
		types[i] = v
	}
	db.InsertTypes(client, types)
	genres := make([]interface{}, len(container.Genres))
	for i, v := range container.Genres {
		genres[i] = v
	}
	db.InsertGenres(client, genres)
	animes := make([]interface{}, len(container.Animes))
	for i, v := range container.Animes {
		animes[i] = v
	}
	insertResult, err := db.InsertAnimes(client, animes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Base de datos llenada en %s con %d animes\n", time.Since(dillDbTime), len(insertResult.InsertedIDs))

	// client, _ := db.SetUp()
	// s, e := db.GetNextSequence(client, "users2")
	// fmt.Println(s, e)
}
