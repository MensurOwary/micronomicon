package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"regexp"
	"strings"
)

type ExtractLink func(baseUrl string, e *colly.HTMLElement) string
type ExtractTitle func(e *colly.HTMLElement) string
type ExtractFields func(title string) string

type Scraped struct {
	BaseUrl       string
	VisitUrl      string
	FileName      string
	Selector      string
	ExtractLink   ExtractLink
	ExtractTitle  ExtractTitle
	ExtractFields ExtractFields
}

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
	_ = collector.Visit(scraped.VisitUrl)
}

func doScrape(e *colly.HTMLElement, scraped Scraped) (string, string, string) {
	title := scraped.ExtractTitle(e)
	link := scraped.ExtractLink(scraped.BaseUrl, e)
	joinedFields := scraped.ExtractFields(title)
	return title, link, joinedFields
}

func handleStopWords(fields []string) []string {
	var tags []string
	for _, field := range fields {
		if !IsStopWord(field) {
			tags = append(tags, strings.ToLower(field))
		}
	}
	return tags
}

var compile = regexp.MustCompile(`^[0-9-.,\s]+`)

func handleTitle(text string) string{
	text = compile.ReplaceAllString(text, "")
	return strings.TrimSpace(strings.ReplaceAll(text, ",", ""))
}
