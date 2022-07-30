package offers

import (
	"encoding/json"
	"log"
	"net/url"
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
func GetOffersJSON(dt time.Time, availIn []byte, meal string, seats int) DiningMap {
   var inMap CurlReturn
   dining := ParseEntities(FetchEntities(dt))

	// When parced should be good enough
   timeNow := time.Now()

   // Parse the JSON
   json.Unmarshal(availIn, &inMap)

   for l, s:= range inMap.Availability {
      if s.HasAvailability {
         idNum := shortID(l)
         v, ok := dining[idNum]
         if !ok {
            log.Printf("No Match for ID: %d, Skipping", idNum)
            continue
         }
         for _, t := range s.Location.Offers {
            w, _ := time.ParseInLocation("2006-01-02T15:04:05", t.Date+"T"+t.Time, disneyTZ)
         	avail := Available{
		   		When:    w,
			   	Service: NormalizeMeal(meal),
			   	Seats:   seats,
			   	Updated: timeNow,
			   }
            avail.URL, _ = url.Parse(v.Location.URL.Scheme+"://"+v.Location.URL.Host+t.URL)
            v.Offers = append(v.Offers, avail)
         }
         dining[idNum] = v
      }
   }

   for i := range dining {
      if dining[i].Offers == nil {
         delete(dining, i)
      }
   }

   return dining
}

// vim: noai:ts=3:sw=3:set expandtab:
