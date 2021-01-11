package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/kuronosu/animeflv-api/pkg/db"
	"github.com/kuronosu/animeflv-api/pkg/server"
	"github.com/kuronosu/animeflv-api/pkg/utils"
)

const (
	latestEpisodesScraperEvent = "StopLatestEpisodes"
)

func main() {
	createFlag := flag.Bool("c", false, "Create directory")
	latestEpisodesFlag := flag.Bool("le", false, "Latest episodes")
	apiServerFlag := flag.Bool("ars", true, "Api Rest server")
	flag.Parse()

	manager := db.CreateManager("animeflv", getAnimeFLVConnectionString()) // DB manager
	evm := utils.EventContainer{}                                          // Events manager
	initEvents(&evm, manager)

	go startSockets(&evm) // Start socket server

	if *createFlag {
		createDirectory(manager)
	}
	if *latestEpisodesFlag {
		evm.Emit(latestEpisodesScraperEvent, true)
	}
	if *apiServerFlag {
		utils.FatalLog(startWebServer(manager))
	}
}

func initEvents(evm *utils.EventContainer, manager db.Manager) {
	aexs := animeExecutionsState{false, false}
	evm.AddEventListener(latestEpisodesScraperEvent, func(state interface{}) {
		switch state.(type) {
		case bool:
			if state.(bool) { // want start and is stopped
				if !aexs.latestEpisodesScraperLoopFlag {
					utils.WarningLog("Last episode scraper ON")
				}
				aexs.startLES(manager)
			} else if !state.(bool) { // want stop and is started
				if aexs.isLatestEpisodesScraperLoopRunning {
					utils.WarningLog("Scraper will be stop after the next iteration")
				}
				aexs.stopLES()
			}
		}
	})
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
