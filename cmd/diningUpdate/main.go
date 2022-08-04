package main

import (
	"disneydining/internal/offers"
	"disneydining/internal/timeout"
	"gopkg.in/ini.v1"
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

	// Read the config file
	cfg, err := ini.LoadSources(ini.LoadOptions{
	   IgnoreInlineComment:         true,
	   UnescapeValueCommentSymbols: true,
      }, "config.ini")
	if err != nil {
		log.Fatal("Failed to read config.ini file")
	}

	// Enable/Disable Timestamps in log
	if !cfg.Section("DEFAULT").Key("timestamps").MustBool(true) {
		clearTimestamps()
	}

   // Set the User Agent to your favorate browser
   if cfg.Section("browser").HasKey("agent") {
      ua := cfg.Section("browser").Key("agent").String()
      log.Printf("User Agent: %s", ua)
      offers.SetUserAgent(ua)
   }
      

	// info
	displayBuildInfo()

	// Start a Timer to make sure we get done
	timeout.StartTimer(cfg.Section("DEFAULT").Key("timeout").MustString("10m"))

	// Read the dining file This will get moved to the config file
	scheduleFile := cfg.Section("DEFAULT").Key("searchfile").MustString("./dining.ini")
	dining, err := ini.Load(scheduleFile)
	if err != nil {
		log.Fatal("Failed to read %s file", scheduleFile)
	}
	log.Printf("Schedule: %s", scheduleFile)

	// get params
	disney := cfg.Section("disney")
   if disney.HasKey("AuthURL") {
      offers.SetAuthURL(disney.Key("AuthURL").String())
   }
   if disney.HasKey("AuthCookie") {
      offers.SetAuthCookieName(disney.Key("AuthCookie").String())
   }
   offers.SetOffersURL(disney.Key("url").MustString(defaultQueryURL))


	defSize := dining.Section("DEFAULT").Key("size").String()
	defEnable := dining.Section("DEFAULT").Key("enabled").MustBool(true)

	allOffers := offers.NewOffers()
	if cfg.Section("DEFAULT").HasKey("saveoffers") {
		offersName := cfg.Section("DEFAULT").Key("saveoffers").String()
		allOffers.LoadOffers(offersName)
		log.Printf("Loaded %d offers at %d locations from %s", allOffers.CountOffers(), len(allOffers), offersName)
	}

   offersChan := make(chan diningInfo, 10)
   offersCnt  := 0

	for _, s := range dining.Sections() {
		searchName := s.Name()
      // default section is not a search section
		if searchName == "DEFAULT" {
			continue
		}
      // Section must be "enabled"
		if !s.Key("enabled").MustBool(defEnable) {
			continue
		}

      // get this sections information
		searchDate := s.Key("date").String()
		searchTime := offers.NormalizeMeal(s.Key("time").MustString("Lunch"))
		searchSize := s.Key("size").MustString(defSize)

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
	if cfg.Section("DEFAULT").HasKey("saveoffers") {
		retention, _ := time.ParseDuration("30m")
		log.Printf("Purged %d old entries",
			allOffers.PurgeOffers(cfg.Section("DEFAULT").Key("offerretention").MustDuration(retention)))
		offersName := cfg.Section("DEFAULT").Key("saveoffers").String()
		log.Printf("Saving %d offers at %d locations to %s", allOffers.CountOffers(), len(allOffers), offersName)
		allOffers.SaveOffers(offersName)
	}
}

// vim: noai:ts=3:sw=3:set expandtab:
