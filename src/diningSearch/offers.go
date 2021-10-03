package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
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

type DiningMap map[int]*DiningStruct

func CheckDate(when string) bool {
   w, err := time.Parse("01/_2/2006", when)
   if err != nil {
      log.Printf("Date Check error %s", err)
      return false
   }
   return w.After(time.Now())
}

func StringIn(set []string, t string) bool {
   for _, this := range set {
      this = strings.TrimSpace(this)
      this = strings.ToLower(this)
      t = strings.TrimSpace(t)
      t = strings.ToLower(t)
      if strings.Contains(t, this) {
         return true
      }
   }
   return false
}

func GetOffers(page string) DiningMap  {
	dining := make(DiningMap)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		log.Fatal(err)
	}

	searchTime := doc.Find("#searchTime-wrapper > div.select-toggle.hoverable > span > span").Contents().Eq(1).Text()

	log.Printf("Time: %s", searchTime)

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

      t := single.Find("div.groupedOffers.show > span > span > a")
		if t.Length() == 0 {
			continue
		}

		if _, ok := dining[idNum]; !ok {
		   dining[idNum] = &DiningStruct{
		   	Name: location,
		   	ID:   idNum,
		   	URL:  url,
		   }
      }
		t.Each(func(i int, s *goquery.Selection) {
         tempTime, _ := s.Attr("data-servicedatetime")
         times = append(times, tempTime)
		})

		for _, v := range times {
			resTime, _ := time.Parse("2006-01-02T15:04:05-07:00", v)
			dining[idNum].Avail = append(dining[idNum].Avail, resTime)
		}

	}

   return dining
}
