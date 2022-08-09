package main

import (
   "encoding/json"
   "os"
	"disneydining/internal/offers"
   "log"
   "fmt"
   "time"
)

type MatchItem struct {
   ID int
   Seats int
   When time.Time
   Updated time.Time
}

type Squelch struct {
   Items []MatchItem
   Updated time.Time
}

type apikey struct {
   Url string `json:"Url"`
   Token string `json:"Token"`
}

func New() *Squelch {
   return new(Squelch)
}

func (s *Squelch)Add(item MatchItem) {
   s.Items = append(s.Items, item)
   s.Updated = time.Now()
}

func (s *Squelch)Mute(item MatchItem)bool {
   for _, i := range s.Items {
      if (i.ID == item.ID) && (i.Seats == item.Seats) && (i.When == item.When) {
         i.Updated = item.Updated
         return true
      }
   }
   return false
}

func loadApiKeys(name string) map[string]apikey {
   j, err := os.ReadFile(name)
   if err != nil {
      log.Printf("loadApiKeys error: %s", err)
      return nil
   }

   rtn := make(map[string]apikey)

   json.Unmarshal(j, &rtn)
   return rtn
}

func (s *Squelch)Load(name string) error{
   j, err := os.ReadFile(name)
   if err != nil {
      return err
   }
   json.Unmarshal(j, s)
   return nil
}

func (s *Squelch)Save(name string) {
   // clear out the expired data
   newS := New()
   timeout, _ := time.ParseDuration("1h")
   for _, i := range s.Items {
      if time.Since(i.Updated) < timeout {
         newS.Add(i)
      }
   }
   newS.Updated = time.Now()

   data, _ := json.Marshal(newS)
   os.WriteFile(name, data, 0644)
}

func doNotify(user string, ds offers.DiningMap, sqlFilename string) {
   sqLst := New()
   sqLst.Load(sqlFilename)

   // load the API keys/token
   apiStorage := loadApiKeys("keys.json")

   for _, m := range ds {
      log.Printf("Searching %s", m.Location.Name)
      for _, o := range m.Offers {
         thisMatch := MatchItem {
            ID: m.Location.ID,
            Seats: o.Seats,
            When: o.When,
            Updated: o.Updated,
         }
         if !sqLst.Mute(thisMatch) {
            sqLst.Add(thisMatch)
            str := fmt.Sprintf("!!Match!! %s (seats:%d@%s) - %s",
               m.Location.Name, o.Seats, o.Service,
               o.When.Format(time.RFC1123))
            log.Printf("%s", str)
            pushover(user, str, apiStorage["pushover"])
         }
      }
   }
   sqLst.Save(sqlFilename)
}

// vim: noai:ts=3:sw=3:set expandtab:
