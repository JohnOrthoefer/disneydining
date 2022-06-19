package offers

import (
	"github.com/PuerkitoBio/goquery"
   "fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
   "sort"
	"time"
)

var disneyTZ *time.Location

const (
   adrDays = 60
)

// get the restaurant name
func (d DiningStruct) RestaurantName() string {
	return d.Location.Name
}

// get the restaurant url, the page about the place
func (d DiningStruct) RestaurantURL() string {
	return d.Location.URL.String()
}

// get the restaurant location, park or resort
func (d DiningStruct) RestaurantLocation(i int) string {
   if i > -1 && i < len(d.Location.Loc) {
      return d.Location.Loc[i]
   }
   return " "
}

// Get an offer time by index
func (d DiningStruct) ByIndex(i int) time.Time {
	return d.Offers[i].When
}

// tells you what dates are currently on file
func (d DiningStruct) GroupByDate() []time.Time {
   var rtn []time.Time
   tmp := make(map[time.Time]bool)

   for _, ent := range d.Offers {
      tmp[makeDate(ent.When)] = true
   }

   for i := range tmp {
      rtn = append(rtn, i)
   }
   return rtn
}

func (d DiningStruct) MealsByDate(t time.Time) []string {
   var rtn []string
   tmp := make(map[string]bool)

   for _, ent := range d.Offers {
      if makeDate(ent.When).Equal(makeDate(t)) {
         tmp[ent.Service] = true
      }
   }

   for i := range tmp {
      rtn = append(rtn, i)
   }
   return rtn
}

func (d DiningStruct) SeatsByMeal(t time.Time, meal string) []int {
   var rtn []int
   tmp := make(map[int]bool)
   for _, ent := range d.Offers {
      if makeDate(ent.When).Equal(makeDate(t)) && ent.Service == meal {
         tmp[ent.Seats] = true
      }
   }
   for i := range tmp {
      rtn = append(rtn, i)
   }
   return rtn
}

func (d DiningStruct) TimesByMealDate(t time.Time, meal string, seats int) []string {
   var rtn []string
   
   warningDuration, _ := time.ParseDuration("10m")
   dangerDuration, _  := time.ParseDuration("20m")

   sort.Slice(d.Offers, func(i, j int) bool {
      return d.Offers[i].When.Before(d.Offers[j].When)
   })

   for _, ent := range d.Offers {
      if makeDate(ent.When).Equal(makeDate(t)) && 
         ent.Service == meal && ent.Seats == seats {
         alert := "success"
         if time.Since(ent.Updated) > warningDuration {
            alert = "warning"
         }
         if time.Since(ent.Updated) > dangerDuration {
            alert = "danger"
         }
         rtn = append(rtn, 
            fmt.Sprintf("<span class=\"text-%s\">%s</span>", 
            alert, ent.When.Format("3:04 PM")))
      }
   }

   return rtn
}

func (d DiningStruct) FindOfferByTime(t time.Time, seats int) int {
   for i, ent := range d.Offers {
      if ent.When.Equal(t) && ent.Seats == seats { 
         return i
      }
   }
   return -1
}

func (d DiningStruct)NewOffers(src DiningStruct) bool {
   for _, ent := range src.Offers {
      if d.FindOfferByTime(ent.When, ent.Seats) < 0 {
         return true
      }
   }
   return false
}

func (d DiningMap)CountOffers()int {
   total := 0
   for _, ent := range d {
      total += len(ent.Offers)
   }
   return total
}

func (d DiningMap)PurgeOffers(pTime time.Duration) int {
   cnt := d.CountOffers()
   for i, ent := range d {
      var newAvail AvailMap
      for _, offer := range ent.Offers {
         if time.Since(offer.Updated) < pTime {
            newAvail = append(newAvail, offer)
         }
      }
      d[i] = DiningStruct{
         Location: ent.Location,
         Offers: newAvail,
      }
   }
   return (cnt - d.CountOffers())
}

// Get seats by index
func (d DiningStruct) Seats(i int) int {
	return d.Offers[i].Seats
}

func NewOffers() DiningMap {
	return make(DiningMap)
}

func splitLocation(s string) []string {
   k:=strings.Split(s, ",")
   for i := range k {
      k[i] = strings.TrimSpace(k[i])
   }
   return k
}

func GetOffers(page string) DiningMap {
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

func (d DiningMap) OfferExists(id int, t time.Time) bool {
	// does the location exist
	if _, ok := d[id]; !ok {
		return false
	}

	// does the time exist
	for _, otime := range d[id].Offers {
		if otime.When.Equal(t) {
			return true
		}
	}

	return false
}

func (d DiningMap) AddOffer(id int, avail Available) bool {
	// offer already exists
	if d.OfferExists(id, avail.When) {
		return false
	}
	foo := d[id]
	foo.Offers = append(foo.Offers, avail)
	return true
}

// vim: noai:ts=3:sw=3:set expandtab:
