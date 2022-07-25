package offers

import (
	"log"
	"time"
)

func init() {
	tz, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal("Can not load US/Eastern Time Zone")
	}
	disneyTZ = tz
   log.Printf("Dining Window after %s, and before %s\n", 
      disneyToday().Format("02 Jan 06"), 
      disneyToday().AddDate(0, 0, adrDays).Format("02 Jan 06"))
}

// vim: noai:ts=3:sw=3:set expandtab:
