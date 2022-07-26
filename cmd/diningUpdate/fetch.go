package main

import (
   "io/ioutil"
   "log"
   "net/http"
   "strings"
   "strconv"
   "time"
)

func urlDate(t string) string {
   w, err := time.Parse("_2 Jan 2006 ", t)
   if err != nil {
      log.Printf("Date Check error %s", err)
      return time.Now().Format("2006-01-02")
   }
   return w.Format("2006-01-02")
}

func urlSize(t string) string {
   return strings.TrimSpace(t)
}

func toInt(t string) int {
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

func FetchOffers( u, d, t, sz string) []byte {
   client := &http.Client{}

    url := u+
      ";entityType=destination/" +
      urlDate(d) + "/" +
      urlSize(sz) + "/?" +
      urlMeal(t)

   req, err := http.NewRequest("GET", url, nil)

   if err != nil {
      log.Fatalln(err)
   }

   req.Header.Set("User-Agent", "Chrome/102.0.0.0")

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

