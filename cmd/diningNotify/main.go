package main

import (
	"disneydining/internal/offers"
	"disneydining/internal/config"
	"log"
   "strings"
   "time"
)

func main() {
   config.ReadConfig("config.ini")

	// Enable/Disable Timestamps in log
	if !config.TimestampsEnabled() {
		clearTimestamps()
	}

	// info
	displayBuildInfo()

   // load the current offers
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

      // make sure the date parses
      thisDate := offers.NormalizeDate(searchDate)

      // if the date is in the past or too far in the future ignore
		if !offers.CheckDate(thisDate) {
			log.Printf("%s Skipping", searchName)
			continue
		}

      thisAfter := offers.DateAddTime(searchDate, s.SearchAfter()).Add(time.Minute * -1)
      thisBefore := offers.DateAddTime(searchDate, s.SearchBefore()).Add(time.Minute * 1)

      for _, r := range searchList {
         for _, sz := range searchSize {
            log.Printf("%s: Checking for %s@%s %s - %s", searchName, 
               fmtDate(thisDate), searchTime, sz, r)
            
            matches := allOffers.Match(
               offers.MatchQuery {
                  Date: thisDate, 
                  DateAfter: thisAfter,
                  DateBefore: thisBefore,
                  Meal: searchTime, 
                  Seats: offers.ToInt(sz), 
                  Name: r,
            })
            if (matches != nil) {
               doNotify(config.NotifyTransport(), s.UserToken(),
                  matches, config.SquelchFilename())
            }
         }
      }
   }
}

// vim: noai:ts=3:sw=3:set expandtab:
