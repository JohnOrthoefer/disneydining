package main

import (
	"disneydining/internal/mail"
	"disneydining/internal/offers"
	"disneydining/internal/timeout"
	"gopkg.in/ini.v1"
	"log"
	"strconv"
	"time"
)

func main() {
	// Read the config file
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal("Failed to read config.ini file")
	}

	// Enable/Disable Timestamps
	if !cfg.Section("DEFAULT").Key("timestamps").MustBool(true) {
		clearTimestamps()
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
	browser := cfg.Section("browser")
	disney := cfg.Section("disney")

	for _, i := range cfg.Section("notify").Keys() {
		switch i.Name() {
		case "server":
			mail.SetSMTPHost(i.String())
		case "from":
			mail.SetFromAddr(i.String())
		default:
			mail.SetToAddr(i.Name(), i.String())
		}
	}

	InitContext(browser.Key("agent").String())

	defSize := dining.Section("DEFAULT").Key("size").String()
	defLocs := dining.Section("DEFAULT").Key("restaurants").Strings(",")
	defEnable := dining.Section("DEFAULT").Key("enabled").MustBool(true)
	defNotify := dining.Section("DEFAULT").Key("notify").String()
	//log.Printf("Notify = %s", defNotify)

	allOffers := offers.NewOffers()
	if cfg.Section("DEFAULT").HasKey("saveoffers") {
		offersName := cfg.Section("DEFAULT").Key("saveoffers").String()
		allOffers.LoadOffers(offersName)
		log.Printf("Loaded %d offers at %d locations from %s", allOffers.CountOffers(), len(allOffers), offersName)
	}

	for _, s := range dining.Sections() {
		searchName := s.Name()
		if searchName == "DEFAULT" {
			continue
		}
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

		searchLocs := append(defLocs, s.Key("restaurants").Strings(",")...)
		if len(searchLocs) < 1 {
			log.Printf("Locations Empty")
			continue
		}

      thisDate, _ := time.Parse("_2 Jan 2006 ", searchDate)
      this := FetchOffers(disney.Key("url").String(), searchDate, searchTime, searchSize)
		thisOffers := offers.GetOffersJSON(thisDate, this, searchTime, toInt(searchSize))
		log.Printf("Looking for %q, list of %d", searchLocs, len(thisOffers))
		for idx, offer := range thisOffers {
			if offers.StringIn(searchLocs, offer.RestaurantName()) {
				msg := mail.MakeMsg(offer)
				seatInt, _ := strconv.Atoi(searchSize)
				if offers.SameDate(offer.ByIndex(0), searchDate) &&
					seatInt == offer.Seats(0) {
					tellWho := s.Key("notify").MustString(defNotify)
					log.Printf("Found!!! (%s)  %s", tellWho, msg)
					if allOffers[idx].NewOffers(offer) {
						mail.Notify(tellWho, msg)
					} else {
						log.Printf("Squeleched")
					}
				} else {
					log.Printf("Mismatch search:%s for %d  found:%s for %d", searchDate, seatInt,
						offer.ByIndex(0).Format("01/02/2006"), offer.Seats(0))
				}
			}
		}
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
		log.Printf("Saving offers %d at %d locations to %s", allOffers.CountOffers(), len(allOffers), offersName)
		allOffers.SaveOffers(offersName)
	}
}

// vim: noai:ts=3:sw=3:set expandtab:
