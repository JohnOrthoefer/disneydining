package offers

import (
	"encoding/json"
   "os"
)

func (d DiningMap) SaveOffers(n string) {
	//    log.Printf("Saving ... %q", d)
	data, _ := json.MarshalIndent(d, "", " ")
	os.WriteFile(n, data, 0644)
}

func (d DiningMap) LoadOffers(n string) {
	j, _ := os.ReadFile(n)
	json.Unmarshal(j, &d)
}

// vim: noai:ts=3:sw=3:set expandtab:

