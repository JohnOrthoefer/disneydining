// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"disneydining/internal/mail"
	"disneydining/internal/offers"
	"disneydining/internal/timeout"
	"gopkg.in/ini.v1"
	"log"
	"os"
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
	scheduleFile := "dining.ini"
	if len(os.Args) > 1 {
		scheduleFile = os.Args[1]
	}
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
				if offers.SameDate(offer.Avail[0].When, searchDate) {
					tellWho := s.Key("notify").MustString(defNotify)
					log.Printf("Found!!! (%s)  %s", tellWho, msg)
					mail.Notify(tellWho, msg)
				} else {
					log.Printf("Date Mismatch %s  %s", searchDate, msg)
				}
			}
		}
	}

	if cfg.Section("DEFAULT").HasKey("saveoffers") {
		offersName := cfg.Section("DEFAULT").Key("saveoffers").String()
		log.Printf("Saving offers to %s", offersName)
		offers.SaveOffers(offersName, allOffers)
	}
	timeout.StopTimer()
}

// vim: noai:ts=3:sw=3:set expandtab:
