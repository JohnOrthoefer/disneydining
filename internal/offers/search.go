package offers

import (
   "time"
//   "log"
   "strings"
)

func (ds DiningStruct) compareName(name string) bool {
   thisName := strings.ToLower(ds.Location.Name)
   queryName := strings.ToLower(name)
   return strings.Contains(thisName, queryName)
}

func (ds DiningStruct) compareDate(d time.Time) bool {
   queryDate := d.Format("20060102")
   for _, i := range ds.Offers {
      if queryDate == i.When.Format("20060102") {
         return true
      }
   }
   return false
}

func (ds DiningStruct) compareMeal(meal string) bool {
   m := strings.ToLower(meal) 
   for _, i := range ds.Offers {
      if strings.ToLower(i.Service) == m {
         return true
      }
   }
   return false
}

func (ds DiningStruct) compareSeats(s int) bool {
   for _, i := range ds.Offers {
      if s == i.Seats {
         return true
      }
   }
   return false
}


func (src DiningMap) Match(d time.Time, meal string, sz int, name string) DiningMap {

   var tmpAvail AvailMap
   rtn := NewOffers()
   queryDate := d.Format("20060102")
   queryMeal := strings.ToLower(meal)

   for _, i := range src {
      if i.compareName(name) {
         thisID := i.Location.ID
         for _, j := range i.Offers {
            if (queryDate == j.When.Format("20060102")) &&
               (queryMeal == strings.ToLower(j.Service)) &&
               (sz == j.Seats) {
               // we have a match
               tmpAvail = append(tmpAvail, j)
            }
         }
         rtn[thisID] = DiningStruct {
            Location: i.Location,
            Offers: tmpAvail,
         }
      }
   }

   return rtn
}
   

