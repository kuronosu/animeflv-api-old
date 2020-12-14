package main

import (
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
	manager := createDBManager()
	go startLatestEpisodesScraper(manager)
	utils.FatalLog(startWebServer(manager))
}

func createDBManager() db.Manager {
	utils.InfoLog("Connect to db")
	manager, err := db.SetUp()
	if err != nil {
		utils.FatalLog(err)
	}
	utils.SuccessLog("Connected to db")
	return manager
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
