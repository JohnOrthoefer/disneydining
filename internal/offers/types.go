package offers

import (
   "net/url"
   "time"
)

type Available struct {
	When    time.Time
	Service string
	Seats   int
	URL     *url.URL
	Updated time.Time
}
type AvailMap []Available

type Restaurant struct {
	Name string
	Loc  []string
	ID   int
	URL  *url.URL
}

type DiningStruct struct {
	Location *Restaurant
	Offers   AvailMap
}

type DiningMap map[int]DiningStruct
