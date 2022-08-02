package offers

import (
   "io/ioutil"
   "log"
   "net/http"
   "net/url"
   "strings"
   "strconv"
   "time"
)

var offersURL *url.URL

func urlDate(t time.Time) string {
   return t.Format("2006-01-02")
}

func urlSize(t string) string {
   return strings.TrimSpace(t)
}

func ToInt(t string) int {
   r, _ := strconv.Atoi(urlSize(t))
   return r
}

func urlMeal(t string) string {
   t = strings.ToLower(t)
   if t == "breakfast" {
      return "mealPeriod=80000712"
   }
   if t == "bunch" {
      return "mealPeriod=80000713"
   }
   if t == "dinner" {
      return "mealPeriod=80000714"
   }
   if t == "lunch" {
      return "mealPeriod=80000717"
   }

   return "mealPeriod=80000717"
}

func SetOffersURL(u string) {
   t, err := url.Parse(u)
   checkErr(err)
   offersURL = t
}

func GetOffersURL() string {
   if offersURL == nil {
      log.Fatal("Offer URL not set")
   }
   return offersURL.String()
}


func FetchOffers(d time.Time, t, sz string) []byte {
   client := &http.Client{}

    url := GetOffersURL() +
      ";entityType=destination/" +
      urlDate(d) + "/" +
      urlSize(sz) + "/?" +
      urlMeal(t)

   req, err := http.NewRequest("GET", url, nil)

   if err != nil {
      log.Fatalln(err)
   }

   req.Header.Set("User-Agent", GetUserAgent())

   resp, err := client.Do(req)
   if err != nil {
      log.Fatalln(err)
   }

   defer resp.Body.Close()
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }

   return body
}

// vim: noai:ts=3:sw=3:set expandtab:

