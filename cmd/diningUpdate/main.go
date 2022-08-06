package main

import (
	"disneydining/internal/offers"
	"disneydining/internal/timeout"
	"disneydining/internal/config"
	"log"
	"time"
   "strings"
   "fmt"
)

type diningInfo struct {
   date       time.Time
   time       string
   size       int
   diningMap  offers.DiningMap
}

func inRequests(dl []diningInfo, r diningInfo) bool {
   for _, i := range dl {
      if (i.date == r.date) && (i.time == r.time) && (i.size == r.size) {
         return true
      }
   }
   return false
}

func fmtDiningInfo(di diningInfo) string {
   return fmt.Sprintf("%s@%s Size=%d",  di.date.Format("2 Jan 2006"), di.time, di.size)
}

func main() {
   var diningRequests []diningInfo

   config.ReadConfig("config.ini")

	// Enable/Disable Timestamps in log
	if !config.TimestampsEnabled() {
		clearTimestamps()
	}

   // Set the User Agent to your favorate browser
   offers.SetUserAgent(config.GetUserAgent())

	// info
	displayBuildInfo()

	// Start a Timer to make sure we get done
	timeout.StartTimer(config.GetRuntimeLimit())

	// get params
   offers.SetAuthURL(config.GetAuthURL())
   offers.SetAuthCookieName(config.GetAuthCookie())
   offers.SetOffersURL(config.GetQueryURL())

	allOffers := offers.NewOffers()
	if sf := config.OffersFilename(); sf != "" {
		allOffers.LoadOffers(sf)
		log.Printf("Loaded %d offers at %d locations from %s", allOffers.CountOffers(), len(allOffers), sf)
	}

   offersChan := make(chan diningInfo, 10)
   offersCnt  := 0

	for _, s := range config.DiningQueries() {
		searchName := s.SearchName()

		searchDate := s.SearchDate()
		searchTime := offers.NormalizeMeal(s.SearchTime())
		searchSize := s.SearchSize()

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

      for _, size := range strings.Fields(searchSize) {
         di := diningInfo{
            date: thisDate, 
            time: searchTime,
            size: offers.ToInt(size),
         }
         if inRequests(diningRequests, di) {
            log.Printf("%s: Skipping already requested: %s", searchName, fmtDiningInfo(di))
            continue
         }
         diningRequests = append(diningRequests, di)
         log.Printf("%s: %s starting", searchName, fmtDiningInfo(di))
         // shoot off a thread to fetch and parse the errors
         offersCnt += 1
         go func(di diningInfo) {
            this := offers.FetchOffers(di.date, di.time, di.size)
            di.diningMap = offers.GetOffersJSON(di.date, 
               this, di.time, di.size)
		      offersChan <- di
         }(di)
      }
   }

   log.Printf("threads running: %d", offersCnt)
   // wait for all the threads to checkin
   for offersCnt > 0 {
      thisOffers := <- offersChan
      offersCnt -= 1
      if len(thisOffers.diningMap) == 0 {
         log.Printf("%s - No Entries returned from thread", fmtDiningInfo(thisOffers))
         continue
      }
      log.Printf("%s- Entries Retrived: %d", fmtDiningInfo(thisOffers),
         len(thisOffers.diningMap))
		// once we've checked for new offers, add this search to the all
		allOffers = allOffers.Join(thisOffers.diningMap)
	}

	timeout.StopTimer()
	if offersName := config.OffersFilename(); offersName != "" {
		log.Printf("Purged %d old entries",
			allOffers.PurgeOffers(config.RetentionTime()))
		log.Printf("Saving %d offers at %d locations to %s", allOffers.CountOffers(), len(allOffers), offersName)
		allOffers.SaveOffers(offersName)
	}
}

// vim: noai:ts=3:sw=3:set expandtab:
