package offers

import (
	"fmt"
	"sort"
	"strings"
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
	dangerDuration, _ := time.ParseDuration("20m")

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

func (d DiningStruct) NewOffers(src DiningStruct) bool {
	for _, ent := range src.Offers {
		if d.FindOfferByTime(ent.When, ent.Seats) < 0 {
			return true
		}
	}
	return false
}

func (d DiningMap) CountOffers() int {
	total := 0
	for _, ent := range d {
		total += len(ent.Offers)
	}
	return total
}

func (d DiningMap) PurgeOffers(pTime time.Duration) int {
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
			Offers:   newAvail,
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
	k := strings.Split(s, ",")
	for i := range k {
		k[i] = strings.TrimSpace(k[i])
	}
	return k
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
