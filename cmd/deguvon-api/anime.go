package main

import (
	"fmt"
	"time"

	"github.com/kuronosu/animeflv-api/pkg/db"
	"github.com/kuronosu/animeflv-api/pkg/scrape"
	"github.com/kuronosu/animeflv-api/pkg/utils"
)

type animeExecutionsState struct {
	latestEpisodesScraperLoopFlag      bool
	isLatestEpisodesScraperLoopRunning bool
}

func (aexs *animeExecutionsState) startLES(manager db.Manager) {
	aexs.latestEpisodesScraperLoopFlag = true
	if !aexs.isLatestEpisodesScraperLoopRunning {
		go aexs.latestEpisodesScraperLoop(manager)
	}
}

func (aexs *animeExecutionsState) stopLES() {
	aexs.latestEpisodesScraperLoopFlag = false
}

func (aexs *animeExecutionsState) latestEpisodesScraperLoop(manager db.Manager) {
	aexs.isLatestEpisodesScraperLoopRunning = true
	for aexs.latestEpisodesScraperLoopFlag {
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
	aexs.isLatestEpisodesScraperLoopRunning = false
}

func createDirectory(manager db.Manager) {
	containerI, _, pgErr, err := scrape.AllAnimesByPage()
	if err != nil {
		utils.ErrorLog(fmt.Sprintf("Error al scrapear las paginas %s, %s", fmt.Sprint(pgErr), err))
		return
	}
	container := containerI.(scrape.AnimeSPContainer)

	manager.DropAll()

	dillDbTime := time.Now()
	// manager.InsertMany("states", container.States...)
	_, err = manager.InsertStates(container.States)
	if err != nil {
		utils.FatalLog("InsertStates ", err)
	}
	_, err = manager.InsertTypes(container.Types)
	if err != nil {
		utils.FatalLog("InsertTypes ", err)
	}
	_, err = manager.InsertGenres(container.Genres)
	if err != nil {
		utils.FatalLog("InsertGenres ", err)
	}
	insertResult, err := manager.InsertAnimes(container.Animes)
	if err != nil {
		utils.FatalLog("InsertAnimes ", err)
	}
	le, _, e := scrape.FetchLatestEpisodes()
	if e == nil {
		manager.SetLatestEpisodes(le)
	}
	utils.SuccessLog(fmt.Sprintf("Base de datos llenada en %s con %d animes", time.Since(dillDbTime), len(insertResult.InsertedIDs)))
}
