package main

import (
	"disneydining/internal/mail"
	"disneydining/internal/offers"
	"disneydining/internal/timeout"
	"disneydining/internal/squelch"
	"gopkg.in/ini.v1"
	"log"
   "strconv"
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

   squelchOffers := squelch.NewSquelch()
   squelchOffers.Load("./squlech.json")

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

	var allOffers []offers.DiningMap

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

		page := GetPage(disney.Key("url").String(), searchDate, searchTime, searchSize)
		thisOffers := offers.GetOffers(page)
		allOffers = append(allOffers, thisOffers)
		log.Printf("Looking for %q, list of %d", searchLocs, len(thisOffers))
		for _, offer := range thisOffers {
			if offers.StringIn(searchLocs, offer.Name) {
				msg := mail.MakeMsg(offer.Name, offer.URL, offer.Avail)
            seatInt, _ := strconv.Atoi(searchSize)
				if offers.SameDate(offer.Avail[0].When, searchDate) &&
               seatInt == offer.Avail[0].Seats {
					tellWho := s.Key("notify").MustString(defNotify)
					log.Printf("Found!!! (%s)  %s", tellWho, msg)
               if !squelchOffers.Mute(msg) {
                  squelchOffers.Add(msg)
					   mail.Notify(tellWho, msg)
               } else {
                  log.Printf("Squeleched")
               }
				} else {
					log.Printf("Mismatch search:%s for %d  found:%s for %d", searchDate, seatInt, 
                  offer.Avail[0].When.Format("01/02/2006"), offer.Avail[0].Seats)
				}
			}
		}
	}

	if cfg.Section("DEFAULT").HasKey("saveoffers") {
		offersName := cfg.Section("DEFAULT").Key("saveoffers").String()
		log.Printf("Saving offers to %s", offersName)
		offers.SaveOffers(offersName, allOffers)
	}
   if cfg.Section("DEFAULT").HasKey("restaurantlist") {
      listName := cfg.Section("DEFAULT").Key("restaurantlist").String()
		log.Printf("Saving Restaurants to %s", listName)
      offers.SaveRestaurants(listName)
   }
   squelchOffers.Save("./squlech.json")
	timeout.StopTimer()
}

// vim: noai:ts=3:sw=3:set expandtab:
