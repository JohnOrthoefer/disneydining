package offers

import (
    "github.com/PuerkitoBio/goquery"
    "log"
    "strconv"
    "strings"
    "time"
    "encoding/json"
    "io/ioutil"
)

type AvailStruct struct {
    When time.Time
    Link string
}

type DiningStruct struct {
    Name  string
    Loc   string
    Meal  string
    ID    int
    URL   string
    Avail []AvailStruct
}

type DiningMap map[int]*DiningStruct
var disneyTZ *time.Location
var disneyToday time.Time

// checks that a and b are the same date
func SameDate(a time.Time, b string) bool {
    w, err := time.ParseInLocation("01/_2/2006", b, disneyTZ)
    if err != nil {
        log.Printf("Date Check error %s", err)
        return false
    }
    n := time.Date(a.Year(), a.Month(), a.Day(), 0, 0, 0, 0, disneyTZ)
    return w.Equal(n)
}
   

// Check that string is after today
func CheckDate(when string) bool {
    w, err := time.ParseInLocation("01/_2/2006", when, disneyTZ)
    if err != nil {
        log.Printf("Date Check error %s", err)
        return false
    }
    n := disneyToday
    return w.After(n)
}

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

func GetOffers(page string) DiningMap  {
    dining := make(DiningMap)

    // Load the HTML document
    doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
    if err != nil {
        log.Fatal(err)
    }

    sel := doc.Find("div.cardLink.finderCard.hasLink")
    meal := doc.Find("#searchTime-wrapper > div.select-toggle.hoverable > span > span").Eq(0).Contents().Text()

    for i := range sel.Nodes {
        single := sel.Eq(i)
        location := single.Find("h2.cardName").Contents().Text()

        u := single.Find("a")
        url, _ := u.Eq(-1).Attr("href")
        id, _ := u.Eq(-1).Attr("id")

        tId := strings.Split(id, ";")
        idNum, _ := strconv.Atoi(tId[0])
        if tId[1] != "entityType=restaurant" {
            continue
        }

        t := single.Find("div.groupedOffers.show > span > span > a")
        if t.Length() == 0 {
            continue
        }

        if _, ok := dining[idNum]; !ok {
            dining[idNum] = &DiningStruct{
                Name: location,
                ID:   idNum,
                URL:  url,
                Meal: meal,
                Loc:  (strings.Split(url, "/"))[4],
            }
        }
        t.Each(func(i int, s *goquery.Selection) {
            tempTime, _ := s.Attr("data-servicedatetime")
            w, _ :=  time.Parse("2006-01-02T15:04:05-07:00", tempTime)
            tempLink, _ := s.Attr("data-bookinglink")
            dining[idNum].Avail = append(dining[idNum].Avail,
            AvailStruct {
                When: w,
                Link: tempLink,
            })
        })
    }

   return dining
}

func SaveOffers(n string, d []DiningMap) {
    data, _ := json.MarshalIndent(d, "", " ")
    ioutil.WriteFile(n, data, 0644)
}

func LoadOffers(n string) []DiningMap {
    j, _ := ioutil.ReadFile(n)
    var rtn []DiningMap
    json.Unmarshal(j, &rtn)
    return rtn
}

func init() {
    tz, err := time.LoadLocation("America/New_York")
    if err != nil {
        log.Fatal("Can not load US/Eastern Time Zone")
    }
    disneyTZ = tz
    n := time.Now().In(disneyTZ)
    disneyToday = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, disneyTZ)
}

// vim: noai:ts=4:sw=4:set expandtab:
