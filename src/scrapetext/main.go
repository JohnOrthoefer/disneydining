package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape() {
	// Request the HTML page.
	res, err := http.Get("http://metalsucks.net")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".left-content article .post-title").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("a").Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}

func main() {
	scrape := "./out.html"
	f, _ := os.Open(scrape)
	defer f.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	// find <li class="card dining show">
	//  or maybe <div class="cardLink finderCard hasLink" role="link>
	//  or maybe <div class="cardLinkContainer">
	// #noAvailability-disneyPicks-default > li:nth-child(1) > div.cardLink.finderCard.hasLink > div.cardLinkContainer
	// Name  <h2 class="cardName"> NAME </h2>
	// locate <div class="groupedOffers show">
	// Time <span class="buttonText"> TIME </span>
	sel := doc.Find("div.cardLinkContainer")

	for i := range sel.Nodes {
		single := sel.Eq(i)
		location := single.Find("h2.cardName").Contents().Text()
		times := single.Find("span.buttonText")
		for j := range times.Nodes {
			log.Printf("%s %s", location, times.Eq(j).Contents().Text())
		}
	}
}
