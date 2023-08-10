package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	//prepare a file to write data
	file, err := os.Create("export.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//create a csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	//header in csv file
	headers := []string{"Title", "Price"}
	writer.Write(headers)

	//create a colly collector
	c := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
	)

	type Book struct {
		Title string
		Price string
	}

	// deal with multi page
	c.OnHTML(".next > a", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	// find the book title and price information
	c.OnHTML(".product_pod", func(e *colly.HTMLElement) {
		book := Book{}
		book.Title = e.ChildAttr(".image_container img", "alt")
		book.Price = e.ChildText(".price_color")
		row := []string{book.Title, book.Price}
		writer.Write(row)
	})

	// is current connection ok?
	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	// current url info
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	startUrl := fmt.Sprintf("https://books.toscrape.com/")
	c.Visit(startUrl)
}
