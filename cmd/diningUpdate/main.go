package main

import (
	"disneydining/internal/offers"
	"disneydining/internal/timeout"
	"gopkg.in/ini.v1"
	"log"
	"time"
   "strings"
)

type diningChannel struct {
   date       time.Time
   time       string
   size       int
   diningMap  offers.DiningMap
}

func main() {
	// Read the config file
	//cfg, err := ini.Load("config.ini")
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

   offersChan := make(chan diningChannel)
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
		log.Printf("%s: %s@%s size:%s", searchName, searchDate, searchTime, searchSize)

      // if the date is in the past or too far in the future ignore
		if !offers.CheckDate(searchDate) {
			log.Printf("Skipping")
			continue
		}

      // make sure the date parses
      thisDate, err := time.Parse("_2 Jan 2006 ", searchDate)
      if err != nil {
         log.Printf("Could not parse %s.. Skipping\n", searchDate)
         continue
      }

      for _, size := range strings.Fields(searchSize) {
         // shoot off a thread to fetch and parse the errors
         offersCnt += 1
         go func(dt time.Time, tm string, s int) {
            log.Printf("%s %s Size = %d", dt.Format("2 Jan 2006"), tm, s)
            this := offers.FetchOffers(dt, tm, s)
		      offersChan <- diningChannel{
               date: dt,
               time: tm,
               size: s,
               diningMap: offers.GetOffersJSON(dt, this, tm, s),
            }
         }(thisDate, searchTime, offers.ToInt(size))
      }
   }

   log.Printf("threads running: %d", offersCnt)
   // wait for all the threads to checkin
   for offersCnt > 0 {
      thisOffers := <- offersChan
      offersCnt -= 1
      if len(thisOffers.diningMap) == 0 {
         log.Printf("%s@%s- No Entries returned from thread", thisOffers.date.Format("2 Jan 2006"), thisOffers.time)
         continue
      }
      log.Printf("%s@%s for %d- Entries Retrived: %d", thisOffers.date.Format("2 Jan 2006"), thisOffers.time, thisOffers.size, len(thisOffers.diningMap))
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
