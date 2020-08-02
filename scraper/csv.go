package scraper

import (
	"github.com/gocarina/gocsv"
	"os"
	"regexp"
)

type Row struct {
	Title string `csv:"title"`
	Link  string `csv:"link"`
	Tags  string `csv:"tags"`
}

var database map[string][]Row

type csvDatabase struct{}

type Scraper interface {
	Database() map[string][]Row
}

func NewScraper() Scraper {
	return &csvDatabase{}
}

func (c *csvDatabase) Database() map[string][]Row {
	if database != nil {
		return database
	}
	clientsFile, err := os.OpenFile("./scraper/dataset/main.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()

	var rows []*Row

	if err := gocsv.UnmarshalFile(clientsFile, &rows); err != nil { // Load rows from file
		println(err.Error())
	}

	database = map[string][]Row{}

	for _, row := range rows {
		for _, tag := range row.tags() {
			entry := database[tag]
			if entry == nil {
				database[tag] = []Row{*row}
			} else {
				database[tag] = append(entry, *row)
			}
		}
	}
	return database
}

func (row Row) tags() []string {
	compile, _ := regexp.Compile(";")
	return compile.Split(row.Tags, -1)
}
