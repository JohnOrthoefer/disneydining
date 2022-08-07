package main

import (
	"disneydining/internal/offers"
	"disneydining/internal/config"
	"log"
   "strings"
   "time"
)

func main() {
   //var diningRequests []diningInfo

   config.ReadConfig("config.ini")

	// Enable/Disable Timestamps in log
	if !config.TimestampsEnabled() {
		clearTimestamps()
	}

	// info
	displayBuildInfo()

	allOffers := offers.NewOffers()
	if sf := config.OffersFilename(); sf != "" {
		allOffers.LoadOffers(sf)
		log.Printf("Loaded %d offers at %d locations from %s", allOffers.CountOffers(), len(allOffers), sf)
	}

	for _, s := range config.DiningQueries() {
		searchName := s.SearchName()
		searchDate := s.SearchDate()
		searchTime := offers.NormalizeMeal(s.SearchTime())
		searchSize := strings.Fields(s.SearchSize())
      searchList := s.RestaurantList()

      // if the date is in the past or too far in the future ignore
		if !offers.CheckDate(searchDate) {
			log.Printf("%s Skipping", searchName)
			continue
		}

      // make sure the date parses
      thisDate := offers.NormalizeDate(searchDate)

      for _, r := range searchList {
         for _, sz := range searchSize {
            log.Printf("%s: %s@%s %s - %s", searchName, 
               fmtDate(thisDate), searchTime, sz, r)
            matches := allOffers.Match(thisDate, searchTime, offers.ToInt(sz), r)
            if (matches != nil) {
               for _, m := range matches {
                  for _, o := range m.Offers {
                     log.Printf("Match %s (seats:%d@%s) - %s", 
                        m.Location.Name, o.Seats, o.Service, 
                        o.When.Format(time.RFC1123))
                  }
               }
            }
         }
      }
   }
}

// vim: noai:ts=3:sw=3:set expandtab:
