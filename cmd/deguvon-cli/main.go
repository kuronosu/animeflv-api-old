package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

// Application constants, defining host, port, and protocol.
const (
	connHost = "localhost"
	connPort = "8081"
	connType = "tcp"
	less     = "LESS" // Latest Episodes Scraper Signal
)

func main() {
	le := flag.Bool("le", true, "Api Rest server")
	required := []string{"le"}
	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			// or possibly use `log.Fatalf` instead of:
			fmt.Fprintf(os.Stderr, "Missing required -%s argument/flag\n", req)
			os.Exit(2) // the same exit code flag.Parse uses
		}
	}

	conn, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	flag.Parse()
	data := fmt.Sprintf("%s:%t\n", less, *le)
	fmt.Print(data)
	conn.Write([]byte(data))
	fmt.Println("Action sent to server")
}
