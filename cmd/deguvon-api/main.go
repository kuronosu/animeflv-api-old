package main

import (
	"log"
	"os"
	"strconv"

	"github.com/kuronosu/animeflv-api/pkg/db"
	"github.com/kuronosu/animeflv-api/pkg/server"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
		log.Printf("Defaulting to port %d", port)
	}
	log.Println("Connect to db")
	manager, err := db.SetUp()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to db")
	log.Fatal(server.New(manager, port).Run())
}
