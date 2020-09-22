package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kuronosu/animeflv-api/pkg/db"
	"github.com/kuronosu/animeflv-api/pkg/server"
)

func main() {
	fmt.Print("Connect to db")
	client, err := db.SetUp()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\rConnected    ")
	s := server.New(client)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Fatal(http.ListenAndServe(":"+port, s.Router()))
}
