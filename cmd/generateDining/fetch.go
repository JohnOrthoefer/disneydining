package main

import (
   "fmt"
   "time"
   "bytes"
   "net/http"
   "io/ioutil"
)

var authToken *http.Cookie

func getAuth() *http.Cookie{
   if authToken != nil {
      fmt.Println("Token Found %s", authToken.String())
      return authToken
   }

   client := &http.Client{}

   url := "https://disneyworld.disney.go.com/finder/api/v1/authz/public"
   req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("{}")))
   checkErr(err)

   req.Header.Set("User-Agent", "Chrome/102.0.0.0")
   resp, err := client.Do(req)
   checkErr(err)

   defer resp.Body.Close()

   for _, c := range resp.Cookies() {
      if c.Name == "__d" {
         authToken = c
         fmt.Println("New Token %s", authToken.String())
         return authToken
      }
   }
   
   fmt.Println("No Token")
   return nil
}

func FetchEntities(d time.Time) []byte {
   client := &http.Client{}

   
   url := "https://disneyworld.disney.go.com/finder/api/v1/explorer-service/list-ancestor-entities/wdw/80007798;entityType=destination/"+d.Format("2006-01-02")+"/dining"
   req, err := http.NewRequest("GET", url, nil)
   checkErr(err)

   req.Header.Set("User-Agent", "Chrome/102.0.0.0")
   req.AddCookie(getAuth())
   resp, err := client.Do(req)
   checkErr(err)

   defer resp.Body.Close()

   body, err := ioutil.ReadAll(resp.Body)
   checkErr(err)

   return body
}

// vim: noai:ts=3:sw=3:set expandtab:
