package offers

import (
   "os"
   "testing"
   "time"
)

func TestToday(t *testing.T) {
   t.Logf("Today- %s\n", disneyToday().Format("02 Jan 2006"))
   t.Logf("ADR (+%d)- %s\n", adrDays, 
      disneyToday().AddDate(0, 0, adrDays).Format("02 Jan 2006"))
}

func TestJSON(t *testing.T) {
   testOffers := "data/15Aug2022.json"
   content, err := os.ReadFile(testOffers)
   if err != nil {
      t.Errorf("%s", err)
   }
   GetOffersJSON(time.Now(), content, "Dinner", 7)
}

// vim: noai:ts=3:sw=3:set expandtab:

