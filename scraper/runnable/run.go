package scraper

import (
	"github.com/gocolly/colly"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

// Start method starts the scraping process
func Start(enabled bool) {
	var wg sync.WaitGroup
	wg.Add(4)
	if enabled {
		go kotlin(&wg)
		go react(&wg)
		go springBoot(&wg)
		go learnXY(&wg)
		wg.Wait()
		merge()
	}
}

func kotlin(wg *sync.WaitGroup) {
	Scraped{
		BaseURL:  "https://kotlinlang.org",
		VisitURL: "https://kotlinlang.org/docs/reference/",
		FileName: "./scraper/dataset/kotlin-docs.csv",
		Selector: "a.tree-item-title",
		ExtractLink: func(baseUrl string, e *colly.HTMLElement) string {
			return baseUrl + e.Attr("href")
		},
		ExtractTitle: func(e *colly.HTMLElement) string {
			return handleTitle(e.Text)
		},
		ExtractFields: func(title string) string {
			fields := handleStopWords(strings.Fields(title))
			fields = append(fields, "kotlin")
			return strings.Join(fields, ";")
		},
	}.Scrape()
	log.Println("Finished scraping Kotlin documentation")
	wg.Add(-1)
}

func react(wg *sync.WaitGroup) {
	Scraped{
		BaseURL:  "https://reactjs.org",
		VisitURL: "https://reactjs.org/docs/getting-started.html",
		FileName: "./scraper/dataset/react-docs.csv",
		Selector: "nav div ul li a[href]",
		ExtractLink: func(baseUrl string, e *colly.HTMLElement) string {
			return baseUrl + e.Attr("href")
		},
		ExtractTitle: func(e *colly.HTMLElement) string {
			return handleTitle(e.Text)
		},
		ExtractFields: func(title string) string {
			fields := handleStopWords(strings.Fields(title))
			fields = append(fields, "react")
			return strings.Join(fields, ";")
		},
	}.Scrape()
	log.Println("Finished scraping React documentation")
	wg.Add(-1)
}

func springBoot(wg *sync.WaitGroup) {
	Scraped{
		BaseURL:  "https://docs.spring.io/spring-boot/docs/current/reference/html/getting-started.html",
		VisitURL: "https://docs.spring.io/spring-boot/docs/current/reference/html/getting-started.html",
		FileName: "./scraper/dataset/spring-boot-docs.csv",
		Selector: "#toc ul li a[href]",
		ExtractLink: func(baseUrl string, e *colly.HTMLElement) string {
			return baseUrl + e.Attr("href")
		},
		ExtractTitle: func(e *colly.HTMLElement) string {
			return handleTitle(e.Text)
		},
		ExtractFields: func(title string) string {
			fields := handleStopWords(strings.Fields(title))
			fields = append(fields, "spring-boot", "spring")
			return strings.Join(fields, ";")
		},
	}.Scrape()
	log.Println("Finished scraping Spring Boot documentation")
	wg.Add(-1)
}

func learnXY(wg *sync.WaitGroup) {
	Scraped{
		BaseURL:  "https://learnxinyminutes.com",
		VisitURL: "https://learnxinyminutes.com",
		FileName: "./scraper/dataset/learnxinyminutes.csv",
		Selector: "td.name a[href]",
		ExtractLink: func(baseUrl string, e *colly.HTMLElement) string {
			return baseUrl + e.Attr("href")
		},
		ExtractTitle: func(e *colly.HTMLElement) string {
			return handleTitle(e.Text)
		},
		ExtractFields: func(title string) string {
			fields := handleStopWords(strings.Fields(title))
			return strings.Join(fields, ";")
		},
	}.Scrape()
	log.Println("Finished scraping Learn X in Y Minutes")
	wg.Add(-1)
}

func merge() {
	if os.Remove("./scraper/dataset/main.csv") != nil {
		log.Fatal("Could not delete the main.csv file")
	}

	open, err := os.Open("./scraper/dataset")
	if err != nil {
		log.Fatal("Could not open the folder")
	}

	info, err := open.Readdir(0)
	if err != nil {
		log.Fatal("Could not read the directory")
	}
	create, _ := os.Create("./scraper/dataset/main.csv")
	_, _ = create.Write([]byte("title,link,tags\n")) // add headers
	for _, fi := range info {
		filename := fi.Name()
		filePath := "./scraper/dataset/" + filename

		log.Printf("Handling %s started", filename)

		file, _ := os.Open(filePath)
		_, _ = io.Copy(create, file)
		log.Printf("Handling %s finished", filename)
	}
}
