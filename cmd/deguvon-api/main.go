package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kuronosu/deguvon-server-go/pkg/db"
	"github.com/kuronosu/deguvon-server-go/pkg/server"
)

func main() {
	fmt.Print("Connect to db")
	client, err := db.SetUp()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\rConnected    ")
	s := server.New(client)
	log.Fatal(http.ListenAndServe(":8080", s.Router()))
}
