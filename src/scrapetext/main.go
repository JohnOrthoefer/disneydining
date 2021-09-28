package main

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type DiningStruct struct {
	Name  string
	ID    int
	URL   string
	Avail []time.Time
}

func main() {
	scrape := "./out.html"
	if len(os.Args) > 1 {
		scrape = os.Args[1]
	}

	dining := make(map[int]*DiningStruct)

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
		var times []string

		single := sel.Eq(i)
		location := single.Find("h2.cardName").Contents().Text()

		u := single.Find("a")
		url, _ := u.Eq(-1).Attr("href")
		id, _ := u.Eq(-1).Attr("id")

		tId := strings.Split(id, ";")
		idNum, _ := strconv.Atoi(tId[0])
		if tId[1] != "entityType=restaurant" {
			continue
		}

		t := single.Find("span.buttonText")
		if t.Length() == 0 {
			continue
		}

		if _, ok := dining[idNum]; ok {
			log.Fatalf("dup ID %d", idNum)
		}
		dining[idNum] = &DiningStruct{
			Name: location,
			ID:   idNum,
			URL:  url,
		}
		t.Each(func(i int, s *goquery.Selection) {
			times = append(times, s.Text())
		})

		for _, v := range times {
			resTime, _ := time.Parse("01/02/2006 3:04 PM", searchDate+" "+v)
			dining[idNum].Avail = append(dining[idNum].Avail, resTime)
		}

		//		if url == "" {
		//			log.Printf("!!!! %s %q", location, times)
		//		} else {
		//			log.Printf("%q", dining[idNum])
		//		}

		//		v, _ := json.Marshal(dining)
		//		log.Print(v)
	}

	log.Printf("DB length %d", len(dining))
	jStr, _ := json.MarshalIndent(dining, "", "\t")
	log.Printf("%s", jStr)
	//	for _, j := range dining {
	//		jStr, _ := json.Marshal(j)
	//		log.Printf("%q", jStr)
	//	}
}
