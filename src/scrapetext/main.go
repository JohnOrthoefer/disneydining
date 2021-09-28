package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strings"
)

func main() {
	scrape := "./out.html"
	if len(os.Args) > 1 {
		scrape = os.Args[1]
	}

	f, _ := os.Open(scrape)
	defer f.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	searchDate, _ := doc.Find("#diningAvailabilityForm-searchDate").Attr("value")
	searchTime := doc.Find("#searchTime-wrapper > div.select-toggle.hoverable > span > span").Contents().Eq(1).Text()
	searchSize := doc.Find("#partySize-wrapper > div.select-toggle.hoverable > span > span").Contents().Eq(1).Text()

	log.Printf("Date: %q", searchDate)
	log.Printf("Time: %s", searchTime)
	log.Printf("Size: %s", searchSize)

	// find <li class="card dining show">
	//  or maybe <div class="cardLink finderCard hasLink" role="link>
	//  or maybe <div class="cardLinkContainer">
	// #noAvailability-disneyPicks-default > li:nth-child(1) > div.cardLink.finderCard.hasLink > div.cardLinkContainer
	// Name  <h2 class="cardName"> NAME </h2>
	// locate <div class="groupedOffers show">
	// Time <span class="buttonText"> TIME </span>
	//	sel := doc.Find("div.cardLinkContainer")
	sel := doc.Find("div.cardLink.finderCard.hasLink")

	for i := range sel.Nodes {
		single := sel.Eq(i)
		location := single.Find("h2.cardName").Contents().Text()
		u := single.Find("a")
		url, _ := u.Eq(-1).Attr("href")
		id, _ := u.Eq(-1).Attr("id")
		tId := strings.Split(id, ";")
		id = tId[0]
		if tId[1] != "entityType=restaurant" {
			continue
		}
		var times []string
		t := single.Find("span.buttonText")
		if t.Length() == 0 {
			continue
		}
		t.Each(func(i int, s *goquery.Selection) {
			times = append(times, s.Text())
		})
		if url == "" {
			log.Printf("!!!! %s %q", location, times)
		} else {
			log.Printf("%s - %s %s %q", id, location, url, times)
		}
		//		for j := range times.Nodes {
		//			log.Printf("%s %s", location, times.Eq(j).Contents().Text())
		//		}
	}
}
