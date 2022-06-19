package offers

import (
   "testing"
)

func TestToday(t *testing.T) {
   t.Logf("Today- %s\n", disneyToday().Format("02 Jan 2006"))
   t.Logf("ADR (+%d)- %s\n", adrDays, 
      disneyToday().AddDate(0, 0, adrDays).Format("02 Jan 2006"))
}

func TestJSON(t *testing.T) {
   if 1 != 1 {
      t.Errorf("In Offer Main")
   }
}

// vim: noai:ts=3:sw=3:set expandtab:

