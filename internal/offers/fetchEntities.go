package offers

import (
   "log"
   "time"
   "bytes"
   "net/http"
   "io"
   "os"
   "net/url"
   "strings"
)

var authCookieName string
var authToken *http.Cookie
var authURL *url.URL

func SetAuthURL(u string) error {
   authURL, err := url.Parse(u)
   if err != nil {
      log.Printf("Failed to Parse %s", u)
      return err
   }
   log.Printf("Auth URL Set: %s", authURL.String())
   return nil
}

func GetAuthURL() string {
   if authURL == nil {
      // this is correct for July 2022
      return "https://disneyworld.disney.go.com/finder/api/v1/authz/public"
   }
   return authURL.String()
}


// Right now only saving need "Auth cookie"
// in the future it should save all needed cookies as defined int he configfile
func SetAuthCookieName(n string) {
   authCookieName = strings.TrimSpace(n)
}

func GetAuthCookieName() string {
   if authCookieName == "" {
      // corrext for July 2022
      return "__d"
   }
   return authCookieName
}

func getAuth() *http.Cookie{
   if authToken != nil {
      return authToken
   }

   client := &http.Client{}

   req, err := http.NewRequest("POST", GetAuthURL(), bytes.NewBuffer([]byte("{}")))
   checkErr(err)

   req.Header.Set("User-Agent", GetUserAgent())
   resp, err := client.Do(req)
   checkErr(err)

   defer resp.Body.Close()

   for _, c := range resp.Cookies() {
      if c.Name == GetAuthCookieName() {
         authToken = c
         return authToken
      }
   }
   
   log.Printf("No Token")
   return nil
}

func FetchEntities(d time.Time) []byte {
   client := &http.Client{}

   
   url := "https://disneyworld.disney.go.com/finder/api/v1/explorer-service/list-ancestor-entities/wdw/80007798;entityType=destination/"+d.Format("2006-01-02")+"/dining"
   req, err := http.NewRequest("GET", url, nil)
   checkErr(err)

   req.Header.Set("User-Agent", GetUserAgent())
   req.AddCookie(getAuth())
   resp, err := client.Do(req)
   checkErr(err)

   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   checkErr(err)

   if false {
      // Saves Entities File as retrived
      os.WriteFile("saveEnt.json", body, 0644)
   }

   return body
}

// vim: noai:ts=3:sw=3:set expandtab:
