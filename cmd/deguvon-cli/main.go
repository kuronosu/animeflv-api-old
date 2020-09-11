package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kuronosu/deguvon/pkg/scrape"
)

func setUpLog() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)
	log.Print("Logging to a file in Go!")
}

func scrapeUrls() {
	urls, err := scrape.GetAnimeURLSFromDirectory()
	if err == nil {
		f, err := os.Create("urls.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		b := 0
		for _, url := range urls {
			l, err := f.WriteString(url + "\n")
			if err != nil {
				fmt.Println(err)
				return
			}
			b += l
		}
		fmt.Println(b, "bytes written successfully")
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func loadURLs() []string {
	file, err := os.Open("urls.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	urls := []string{}
	for scanner.Scan() {
		urls = append(urls, scrape.AnimeFlvURL+scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return urls
}

func scrapeByPage() {
	start := time.Now()
	pages, err := scrape.GetDirectoryPageCount()
	if err != nil {
		return
	}
	errc := 0
	errs := []scrape.Result{}
	for _, page := range scrape.MakeRange(1, pages) {
		start2 := time.Now()
		urls, err := scrape.GetAnimeURLSFromDirectoryPage(page)
		errcp := 0
		animes := []scrape.Anime{}
		if err == nil {
			results := scrape.AsyncHTTPGets(urls, scrape.HandleAnimeScrape)
			for _, result := range results {
				if result.HTTPResponse.Err != nil || result.Document.Err != nil {
					errc++
					errcp++
					errs = append(errs, *result)
					continue
				}
				animes = append(animes, result.HandledResponse.(scrape.Anime))
			}
			// time.Sleep(500 * time.Millisecond)
			fmt.Printf("Scraped page: #%d\twith %d animes and %d errors\tin %s\n", page, len(animes), errcp, time.Since(start2))
		} else {
			fmt.Println("Error obteniendo las urls de los animes en la pagina: ", page)
		}
	}
	fmt.Println(errs)
	// return errc
	fmt.Printf("Completado en %s con %d errores", time.Since(start), errc)
}

func main() {

	// setUpLog()
	scrapeByPage()
	// anime, err := scrape.GetAnime("/anime/koi-to-producer-evollove")
	// anime, err := scrape.GetAnime("/anime/the-idolmster")
	// anime, err := scrape.GetAnime("/anime/over-drie")
	// anime, err := scrape.GetAnime("/anime/chihayafuru-3")
	// anime, err := scrape.GetAnime("/anime/kami-no-tou")
	// if err != nil {
	// 	fmt.Println("Error ", err)
	// }
	// fmt.Println("Anime ", anime)

	// errCount := 0
	// errs := []error{}
	// urls, _ := scrape.GetAnimeURLSFromDirectory()
	// count := 1
	// for _, url := range urls {
	// 	fmt.Print(fmt.Sprintf("Anime #%d ", count))
	// 	count++
	// 	_, err := scrape.GetAnime(url)
	// 	if err != nil {
	// 		errCount++
	// 		errs = append(errs, err)
	// 	}
	// }
	// fmt.Println(errs)
}

func startScrape() {
	val, err := scrape.GetDirectoryPageCount()
	if err != nil {
		fmt.Println("Error ", err)
	}
	fmt.Println(val)
}

func saturateAnimeFLV() {
	start := time.Now()
	for i := 0; i < 1000; i++ {
		reqFlv()
	}
	fmt.Println("1000 request in ", time.Since(start))
}

func reqFlv() {
	start := time.Now()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	response, err := client.Get(scrape.AnimeFlvURL)
	if err != nil {
		return
	}
	defer response.Body.Close()
	fmt.Println("Number of bytes copied to STDOUT: ", response.StatusCode, " In ", time.Since(start))
}
