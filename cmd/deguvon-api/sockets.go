package main

import (
	"bufio"
	"net"
	"strings"

	"github.com/kuronosu/animeflv-api/pkg/utils"
)

const (
	connHost = "localhost"
	connPort = "8081"
	connType = "tcp"
	less     = "LESS" // Latest Episodes Scraper Signal
)

func startSockets(evm *utils.EventContainer) {
	utils.InfoLog("Starting " + connType + " server on " + connHost + ":" + connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		utils.FatalLog("Error listening:", err.Error())
	}

	defer l.Close()

	for {
		// Listen for an incoming connection.
		c, err := l.Accept()
		if err != nil {
			return
		}
		// Handle connections concurrently in a new goroutine.
		go handleConnection(c, evm)
	}
}

func handleConnection(conn net.Conn, evm *utils.EventContainer) {
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')
	defer conn.Close()
	if err != nil {
		return
	}
	signals := strings.Split(string(buffer[:len(buffer)-1]), ":")
	if len(signals) == 2 {
		event := signals[0]
		payload := signals[1]
		switch event {
		case less:
			if payload == "true" {
				evm.Emit(latestEpisodesScraperEvent, true)
			} else if payload == "false" {
				evm.Emit(latestEpisodesScraperEvent, false)
			}
		}
	}
}
