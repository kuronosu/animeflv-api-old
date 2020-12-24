package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kuronosu/animeflv-api/pkg/db"
	"github.com/kuronosu/animeflv-api/pkg/scrape"
	"github.com/kuronosu/animeflv-api/pkg/server"
	"github.com/kuronosu/animeflv-api/pkg/utils"
)

func main() {
	manager := db.CreateManager("animeflv", getAnimeFLVConnectionString())
	createFlag := flag.Bool("c", false, "Create directory")
	latestEpisodesFlag := flag.Bool("le", false, "Latest episodes")
	apiServerFlag := flag.Bool("ars", true, "Api Rest server")
	flag.Parse()

	if *createFlag {
		createDirectory(manager)
	}
	if *latestEpisodesFlag {
		utils.InfoLog("Last episode scraper ON")
		go startLatestEpisodesScraper(manager)
	}
	if *apiServerFlag {
		utils.FatalLog(startWebServer(manager))
	}
}

func getAnimeFLVConnectionString() string {
	connectionString := os.Getenv("AnimeFLVConnectionString")
	if connectionString == "" {
		connectionString = "mongodb://localhost:27017"
	}
	return connectionString
}

func startWebServer(manager db.Manager) error {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
		utils.WarningLog(fmt.Sprintf("Defaulting to port %d", port))
	}
	return server.New(manager, port).Run()
}

func startLatestEpisodesScraper(manager db.Manager) {
	for {
		utils.ColoredLog(utils.LightPurple, "Getting latest episodes")
		le, a, e := scrape.FetchLatestEpisodes()
		if e == nil {
			manager.SetLatestEpisodes(le)
			_, in, _ := manager.UpdateOrInsertAnimes(a.Animes)
			if len(in) > 0 {
				relatedURLs := []string{}
				for _, anime := range in {
					for _, rel := range anime.Relations {
						relatedURLs = append(relatedURLs, rel.URL)
					}
				}
				states, _ := manager.LoadStates()
				genres, _ := manager.LoadGenres()
				types, _ := manager.LoadTypes()
				container := scrape.AnimeSPContainer{
					States: states, Types: types,
					Genres: genres, Animes: []scrape.Anime{}}
				scrape.GetAnimes(relatedURLs, &container)
				manager.UpdateOrInsertAnimes(container.Animes)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func createDirectory(manager db.Manager) {
	containerI, _, _ := scrape.AllAnimesByPage()
	container := containerI.(scrape.AnimeSPContainer)

	manager.DropAll()

	dillDbTime := time.Now()
	// manager.InsertMany("states", container.States...)
	manager.InsertStates(container.States)
	manager.InsertTypes(container.Types)
	manager.InsertGenres(container.Genres)
	insertResult, err := manager.InsertAnimes(container.Animes)
	if err != nil {
		utils.FatalLog(err)
	}
	le, _, e := scrape.FetchLatestEpisodes()
	if e == nil {
		manager.SetLatestEpisodes(le)
	}
	utils.SuccessLog(fmt.Sprintf("Base de datos llenada en %s con %d animes", time.Since(dillDbTime), len(insertResult.InsertedIDs)))
}
