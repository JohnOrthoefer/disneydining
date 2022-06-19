package offers

import (
	"log"
	"strings"
	"time"
)

func makeDate(a time.Time)time.Time {
   return time.Date(a.Year(), a.Month(), a.Day(), 0, 0, 0, 0, disneyTZ)
}

// Join Dining Map, src with Dining Map dst
func (dst DiningMap) Join(src DiningMap) DiningMap {
	for idx, ent := range src {
		if _, ok := dst[idx]; !ok {
         log.Printf("Join- %s (id:%d) does not exist in dst", ent.RestaurantName(), idx)
			// move the whole thing
			dst[idx] = ent
			continue
		}
		// just move the times
      v := dst[idx]
      start := len(v.Offers)
		for _, tent := range ent.Offers {
         offset := dst[idx].FindOfferByTime(tent.When, tent.Seats)
         if offset == -1 {
            v.Offers = append(v.Offers, tent)
         } else {
            v.Offers[offset] = tent
         }
		}
      dst[idx] = v
      log.Printf("Join-  %s (%d) added %d entries", ent.RestaurantName(), idx, (len(v.Offers)-start))
	}
   return dst
}

// get what time it is at disney world NOW
func disneyToday() time.Time {
	n := time.Now().In(disneyTZ)
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, disneyTZ)
}

func NormalizeDate(b string) time.Time {
   w, err := time.ParseInLocation("_2 Jan 2006 ", b, disneyTZ)
   if err != nil {
      log.Printf("Date Check error %s", err)
      return disneyToday()
   }
   return w
}

// checks that a and b are the same date
func SameDate(a time.Time, b string) bool {
	w, err := time.ParseInLocation("_2 Jan 2006", b, disneyTZ)
	if err != nil {
		log.Printf("Date Check error %s", err)
		return false
	}
	n := time.Date(a.Year(), a.Month(), a.Day(), 0, 0, 0, 0, disneyTZ)
	return w.Equal(n)
}

// Check that string is after today
func CheckDate(when string) bool {
	w, err := time.ParseInLocation("_2 Jan 2006", when, disneyTZ)
	if err != nil {
		log.Printf("Date Check error %s", err)
		return false
	}

	return w.After(disneyToday()) && 
      w.Before(disneyToday().AddDate(0,0,adrDays))
}

// Match if a any of the substrings, set, appear in string, t.
func StringIn(set []string, t string) bool {
	for _, this := range set {
		this = strings.ToLower(strings.TrimSpace(this))
		t = strings.TrimSpace(t)
		t = strings.ToLower(t)
		if strings.Contains(t, this) {
			return true
		}
	}
	return false
}
