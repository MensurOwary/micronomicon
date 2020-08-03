package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"regexp"
	"strings"
)

type extractLink func(baseUrl string, e *colly.HTMLElement) string
type extractTitle func(e *colly.HTMLElement) string
type extractFields func(title string) string

// Scraped represents the scraped endpoint
type Scraped struct {
	BaseURL       string
	VisitURL      string
	FileName      string
	Selector      string
	ExtractLink   extractLink
	ExtractTitle  extractTitle
	ExtractFields extractFields
}

// Scrape method scrapes the given resource (Scraped)
func (scraped Scraped) Scrape() {
	handle, err := os.Create(scraped.FileName)

	if err != nil {
		panic("File creation failed")
	}
	defer handle.Close()

	collector := colly.NewCollector()

	collector.OnHTML(scraped.Selector, func(e *colly.HTMLElement) {
		title, link, joinedFields := doScrape(e, scraped)

		_, _ = fmt.Fprintf(handle, "%s,%s,%s\n", title, link, joinedFields)
	})
	_ = collector.Visit(scraped.VisitURL)
}

func doScrape(e *colly.HTMLElement, scraped Scraped) (string, string, string) {
	title := scraped.ExtractTitle(e)
	link := scraped.ExtractLink(scraped.BaseURL, e)
	joinedFields := scraped.ExtractFields(title)
	return title, link, joinedFields
}

func handleStopWords(fields []string) []string {
	var tags []string
	for _, field := range fields {
		if !isStopWord(field) {
			tags = append(tags, strings.ToLower(field))
		}
	}
	return tags
}

var nonText = regexp.MustCompile(`^[0-9-.,\s]+`)

func handleTitle(text string) string {
	text = nonText.ReplaceAllString(text, "")
	return strings.TrimSpace(strings.ReplaceAll(text, ",", ""))
}
