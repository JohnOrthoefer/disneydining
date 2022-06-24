package offers

import (
   "encoding/json"
   "net/url"
   "strconv"
   "strings"
)

type LocStruct struct {
   ID string `json:"id"`
   Title string `json:"title"`
   URLFriendlyId string `json:"urlFriendlyId"`
   LocationType string `json:"locationType"`
}

type ResStruct struct {
   ID string `json:"id"`
   URLFriendlyId string `json:"urlFriendlyId"`
   LocationName string `json:"locationName"`
   ResURL string `json:"url"`
   Name string `json:"name"`
}

type CurlLocation struct {
   Error    string `json:"error"`
   Locations  []LocStruct `json:"locations"`
   Results []ResStruct `json:"results"`
}

func ParseEntities(results []byte) DiningMap {
   var inLoc CurlLocation
   json.Unmarshal(results, &inLoc)

   dining := NewOffers()

   for _, s := range inLoc.Results {
      sID := strings.Split(s.ID, ";")
      idNum, _ := strconv.Atoi(sID[0])
      v := dining[idNum]
      thisURL, _ := url.Parse("https://disneyworld.disney.go.com"+s.ResURL)
      t := &Restaurant{
         Name: s.Name,
         Loc:  []string{s.LocationName},
         ID:   idNum,
         URL:  thisURL,
      }
      v.Location = t
      v.Offers = nil 
      dining[idNum] = v
   }

   return dining
}

// vim: noai:ts=3:sw=3:set expandtab:
