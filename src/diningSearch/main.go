// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
   "os"
   "log"
   "gopkg.in/ini.v1"
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
   StartTimer(cfg.Section("DEFAULT").Key("timeout").MustString("10m"))

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

   InitContext(browser.Key("agent").String())

   defSize := dining.Section("DEFAULT").Key("size").String()
   defLocs := dining.Section("DEFAULT").Key("restaurants").Strings(",")
   defEnable := dining.Section("DEFAULT").Key("enabled").MustBool(true)

   var allOffers []DiningMap

   for _, s := range dining.Sections() {
      searchName := s.Name()
      if searchName == "DEFAULT" { continue }
      if !s.Key("enabled").MustBool(defEnable) { continue }

      searchDate := s.Key("date").String()
      searchTime := s.Key("time").String()
      searchSize := s.Key("size").MustString(defSize)
      log.Printf("%s: %s@%s size:%s", searchName, searchDate, searchTime, searchSize)

      if !CheckDate(searchDate) {
         log.Printf("Skipping")
         continue
      }

      searchLocs := append(defLocs, s.Key("restaurants").Strings(",")...)
      if len(searchLocs) < 1 {
         log.Printf("Locations Empty")
         continue
      }
      
      page := GetPage(disney.Key("url").String(), searchDate, searchTime, searchSize)
      offers := GetOffers(page)
      allOffers = append(allOffers, offers)
      log.Printf("Looking for %q, list of %d", searchLocs, len(offers))
      for _, offer := range offers {
         if StringIn(searchLocs, offer.Name) {
            s := MakeMsg(offer.Name, offer.URL, offer.Avail)
            if SameDate(offer.Avail[0], searchDate) {
               log.Printf("Found!!!  %s", s)
               Notify(s)
            } else {
               log.Printf("Date Mismatch %s  %s", searchDate, s)
            }
         }
      }
   }

   if cfg.Section("DEFAULT").HasKey("saveoffers") {
      offersName := cfg.Section("DEFAULT").Key("saveoffers").String()
      log.Printf("Saving offers to %s", offersName)
      SaveOffers(offersName, allOffers)
   }
   StopTimer()
}
