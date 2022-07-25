package main

import (
   "fmt"
//   "io/ioutil"
   "time"
   "encoding/json"
)

func main() {
   d, err := time.Parse("2006-01-02", "2022-06-23")
   checkErr(err)

   entities := FetchEntities(d)

/*
   testLocations := "data/locations.json"
   location, _ := ioutil.ReadFile(testLocations)
*/

   diningMap := ParseEntities(entities)

   dataLoc, _ := json.MarshalIndent(diningMap, "", "  ")
   fmt.Printf("%s\n", dataLoc)
}

// vim: noai:ts=3:sw=3:set expandtab:

