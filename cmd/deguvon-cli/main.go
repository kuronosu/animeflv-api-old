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
	animes, _, _ := scrape.AllAnimesByPage()
	// stopSpinner = true
	// <-c // wait spinner stop
	dillDbTime := time.Now()
	insertResult, err := db.InsertAnimes(client, animes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Base de datos llenada en %s con %d animes\n", time.Since(dillDbTime), len(insertResult.InsertedIDs))
}
