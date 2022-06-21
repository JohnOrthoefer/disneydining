package offers

import (
	"encoding/json"
   "fmt"
	"log"
	"net/url"
//	"regexp"
	"strconv"
	"strings"
	"time"
)

type OffersAvail struct {
   Date     string `json:"date"`
   Time     string `json:"time"`
   Label    string `json:"label"`
   URL      string `json:"url"`
   ProductType string `json:"productType"`
}

type SingleLocation struct {
   UnavailableReason string   `json:"unavailableReason"`
   Title             string   `json:"title"`
   Offers            []OffersAvail `json:"offers"`
}

type CurlAvail struct {
   HasAvailability   bool     `json:"hasAvailability"`
   AvailabilitySearchDate string `json:"availabilitySearchDate"`
   Location          SingleLocation   `json:"singleLocation"`
}

type CurlReturn struct {
   Error    string `json:"error"`
   Availability map[string]CurlAvail `json:"availability"`
}

func shortID(r string) int {
   id := strings.Split(r, ";")
   idNum, _ := strconv.Atoi(id[0])
   return idNum
}

// JSON List of offers return a DiningMap
func GetOffersJSON(availIn []byte, meal string, seats int) DiningMap {
   var inMap CurlReturn
   dining := NewOffers()

	// When parced should be good enough
   timeNow := time.Now()

   // Parse the JSON
   json.Unmarshal(availIn, &inMap)
   json.Unmarshal([]byte(jsonLocations), &dining)

//  	data, _ := json.MarshalIndent(inMap, "", "  ")
//	fmt.Printf("%s\n", data)
   
   for l, s:= range inMap.Availability {
      if s.HasAvailability {
         idNum := shortID(l)
         v, ok := dining[idNum]
         if !ok {
            log.Printf("No Match for ID: %d, Skipping", idNum)
            continue
         }
         fmt.Printf("%d %q\n", idNum, v.Location.Name)
         for _, t := range s.Location.Offers {
            w, _ := time.ParseInLocation("2006-01-02T15:04:05", t.Date+"T"+t.Time, disneyTZ)
         	avail := Available{
		   		When:    w,
			   	Service: meal,
			   	Seats:   seats,
			   	Updated: timeNow,
			   }
            avail.URL, _ = url.Parse(v.Location.URL.Scheme+"://"+v.Location.URL.Host+t.URL)
            v.Offers = append(v.Offers, avail)
            fmt.Printf("%q\n", v.Offers)
         }
         dining[idNum] = v
      }
   }

   for i := range dining {
      if dining[i].Offers == nil {
         delete(dining, i)
      }
   }
  	data, _ := json.MarshalIndent(dining, "", "  ")
	fmt.Printf("%s\n", data)

   return dining

/*
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
*/
}

// vim: noai:ts=3:sw=3:set expandtab:
