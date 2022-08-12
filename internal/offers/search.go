package offers

import (
   "time"
//   "log"
   "strings"
)

type MatchQuery struct {
   Date time.Time
   DateAfter time.Time
   DateBefore time.Time
   Meal string
   Seats int
   Name string
}

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

func (src DiningMap) Match(q MatchQuery) DiningMap {

   rtn := NewOffers()
   queryDate := q.Date.Format("20060102")
   queryMeal := strings.ToLower(q.Meal)

   for _, i := range src {
      if i.compareName(q.Name) {
         var tmpAvail AvailMap
         thisID := i.Location.ID
         for _, j := range i.Offers {
/*
            log.Printf("When: %s, Before(%s)=%t, After(%s)=%t",
               j.When.String(), 
               q.DateBefore.String(), j.When.Before(q.DateBefore),
               q.DateAfter.String(), j.When.After(q.DateAfter))
*/
            if (queryDate == j.When.Format("20060102")) &&
               (queryMeal == strings.ToLower(j.Service)) &&
               (j.When.Before(q.DateBefore)) &&
               (j.When.After(q.DateAfter)) &&
               (q.Seats == j.Seats) {
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
   

