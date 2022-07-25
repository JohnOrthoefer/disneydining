package offers

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)


// HTML Page return a DiningMap
func GetOffersHTML(page string) DiningMap {
	dining := NewOffers()
	// When parced should be good enough
	timeNow := time.Now()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		log.Fatal(err)
	}

	// gets the list of "cardLinks"
	sel := doc.Find("div.cardLink.finderCard.hasLink")

	// find the Meal that the search matched
	meal := doc.Find("#searchTime-wrapper > div.select-toggle.hoverable > span > span").Eq(0).Contents().Text()
	meal = meal[30:]
	meal = meal[:len(meal)-20]

	// find how many "seats" the search was for
	seatsRE := regexp.MustCompile(`[0-9]+`)
	seatsTxt := doc.Find("#partySize-wrapper > div.select-toggle.hoverable > span > span").Eq(0).Contents().Text()
	seats, err := strconv.Atoi(seatsRE.FindString(seatsTxt))
	log.Printf("seatsTxt: %s; Match: %d", seatsTxt, seats)

	// loop though each "cardLink"
	for i := range sel.Nodes {
		// get the cursor on the current card
		single := sel.Eq(i)

		// location: name of the restaurant
		location := single.Find("h2.cardName").Contents().Text()

		u := single.Find("a").Eq(-1)
		// url: link to the restaurant page or the default dining link
		url, err := url.Parse(u.AttrOr("href", "https://disneyworld.disney.go.com/dining/"))
		if err != nil {
			log.Printf("URL for %s did not parse", location)
		}

		// idNum: is unique ID for each location/event
		id, _ := u.Attr("id")
		tId := strings.Split(id, ";")
		idNum, _ := strconv.Atoi(tId[0])

		if tId[1] != "entityType=restaurant" {
			log.Printf("Skipping %s", tId[1])
			continue
		}

		locName := splitLocation(single.Find("div.descriptionLines > span:nth-child(3)").Contents().Text())
      
		// idNum should only show up once.
		v, ok := dining[idNum]
		if !ok {
			//log.Printf("%s: href=%s", location, url.Path)
			t := &Restaurant{
				Name: location,
				Loc:  locName,
				ID:   idNum,
				URL:  url,
			}
			v.Location = t
			v.Offers = nil
		} else {
			log.Printf("Duplicate entries ID: %d, Skipping", idNum)
			continue
		}

		// see if there are any offers to poulate
		t := single.Find("div.groupedOffers.show > span > span > a")
		if t.Length() == 0 {
			log.Printf("No offers for %s", location)
			continue
		}

		t.Each(func(i int, s *goquery.Selection) {
			tempTime, _ := s.Attr("data-servicedatetime")
         log.Printf("ServciceDateTime = %s\n", tempTime)
			w, _ := time.Parse("2006-01-02T15:04:05-07:00", tempTime)
			tempLink, _ := s.Attr("data-bookinglink")
			tempURL, _ := url.Parse(tempLink)
			avail := Available{
				When:    w,
				Service: meal,
				Seats:   seats,
				URL:     tempURL,
				Updated: timeNow,
			}
			v.Offers = append(v.Offers, avail)
		})
		dining[idNum] = v
	}

	// log.Printf("getOffers: %q", dining)
	return dining
}

// vim: noai:ts=3:sw=3:set expandtab:
