package main

// go run main.go > ../../internal/offers/locations.goX
// mv ../../internal/offers/locations.goX ../../internal/offers/locations.go

import (
   "encoding/json"
   "disneydining/internal/offers"
   "fmt"
   "io/ioutil"
   "net/url"
   "strconv"
   "strings"
)

type OffersAvail struct {
   Date     string `json:"date"`
   Time     string `json:"time"`
   Label    string `json:"label"`
   URL      string `json:"url"`
   ProductType string `json:"productType"`
}

type SingleLocation struct {
   UnavailableReason string   `json:"unavailableReason"`
   Title             string   `json:"title"`
   Offers            []OffersAvail `json:"offers"`
}

type CurlAvail struct {
   HasAvailability   bool     `json:"hasAvailability"`
   AvailabilitySearchDate string `json:"availabilitySearchDate"`
   Location          SingleLocation   `json:"singleLocation"`
}

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

type CurlReturn struct {
   Error    string `json:"error"`
   Availability map[string]CurlAvail `json:"availability"`
}

type CurlLocation struct {
   Error    string `json:"error"`
   Locations  []LocStruct `json:"locations"`
   Results []ResStruct `json:"results"`
}

func main() {
   testLocations := "data/locations.json"
   location, _ := ioutil.ReadFile(testLocations)

   var inLoc CurlLocation
   json.Unmarshal(location, &inLoc)

   dining := offers.NewOffers()

   for _, s := range inLoc.Results {
      sID := strings.Split(s.ID, ";")
      idNum, _ := strconv.Atoi(sID[0])
      v := dining[idNum]
      thisURL, _ := url.Parse("https://disneyworld.disney.go.com"+s.ResURL)
      t := &offers.Restaurant{
         Name: s.Name,
         Loc:  []string{s.LocationName},
         ID:   idNum,
         URL:  thisURL,
      }
      v.Location = t
      v.Offers = nil 
      dining[idNum] = v
   }
   dataLoc, _ := json.MarshalIndent(dining, "", "  ")
   fmt.Printf("package offers\n\nconst jsonLocations = `%s`\n", dataLoc)
}

// vim: noai:ts=3:sw=3:set expandtab:
