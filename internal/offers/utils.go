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
			// move the whole thing
			dst[idx] = ent
			continue
		}
		// just move the times
      v := dst[idx]
		for _, tent := range ent.Offers {
         offset := dst[idx].FindOfferByTime(tent.When, tent.Seats)
         if offset == -1 {
            v.Offers = append(v.Offers, tent)
         } else {
            v.Offers[offset] = tent
         }
		}
      dst[idx] = v
	}
   return dst
}

// get what time it is at disney world NOW
func disneyToday() time.Time {
	n := time.Now().In(disneyTZ)
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, disneyTZ)
}

func NormalizeMeal(s string) string {
   lcs := strings.ToLower(strings.TrimSpace(s))
   switch lcs {
   case "breakfast":
      return "Breakfast"
   case "brunch":
      return "Brunch"
   case "lunch":
      return "Lunch"
   case "dinner":
      return "Dinner"
   }
   return lcs
}

func NormalizeDate(b string) time.Time {
   var rtn time.Time

   switch b := strings.TrimSpace(strings.ToLower(b)); b {
   case "today":
      rtn = disneyToday()
   case "tomorrow":
      rtn = disneyToday().Add(time.Hour * 24)
   default:
      w, err := time.ParseInLocation("_2 Jan 2006", b, disneyTZ)
      if err != nil {
         log.Printf("NormalizeDate error %s", err)
         rtn = disneyToday()
      }
      rtn = w
   }
   return rtn
}

func DateAddTime(d, h string) time.Time {
   rtn, err := time.ParseInLocation("_2 Jan 2006 15:04", d+" "+h, disneyTZ)
   if err != nil {
      log.Printf("DateAddTime error %s", err)
      return disneyToday()
   }
   return rtn
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
func CheckDate(w time.Time) bool {
	return w.After(disneyToday().AddDate(0, 0, -1)) && 
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
