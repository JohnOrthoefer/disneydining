package main

import (
	"disneydining/internal/offers"
	"disneydining/internal/config"
	"log"
	"time"
   "strings"
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
      thisDate, err := time.Parse("_2 Jan 2006 ", searchDate)
      if err != nil {
         log.Printf("%s: Could not parse %s.. Skipping\n", searchName, searchDate)
         continue
      }
      for _, r := range searchList {
         for _, sz := range searchSize {
            log.Printf("%s: %s@%s %s - %s", searchName, 
               fmtDate(thisDate), searchTime, sz, r)
         }
      }
   }
}

// vim: noai:ts=3:sw=3:set expandtab:
