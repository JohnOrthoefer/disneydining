package main

import (
	"disneydining/internal/offers"
	"disneydining/internal/timeout"
	"gopkg.in/ini.v1"
	"log"
	"time"
)

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

		searchDate := s.Key("date").String()
		searchTime := s.Key("time").String()
		searchSize := s.Key("size").MustString(defSize)
		log.Printf("%s: %s@%s size:%s", searchName, searchDate, searchTime, searchSize)

		if !offers.CheckDate(searchDate) {
			log.Printf("Skipping")
			continue
		}

      thisDate, _ := time.Parse("_2 Jan 2006 ", searchDate)
      this := offers.FetchOffers(searchDate, searchTime, searchSize)
		thisOffers := offers.GetOffersJSON(thisDate, this, searchTime, offers.ToInt(searchSize))
      log.Printf("Entries Retrived: %d", len(thisOffers))
		// once we've checked for new offers, add this search to the all
		allOffers = allOffers.Join(thisOffers)

		if cfg.Section("DEFAULT").HasKey("saveoffers") {
			offersName := cfg.Section("DEFAULT").Key("saveoffers").String()
			log.Printf("Saving offers to %s", offersName)
			allOffers.SaveOffers(offersName)
		}
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
