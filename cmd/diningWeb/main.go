package main

import (
	"encoding/json"
	"errors"
    "disneydining/internal/config"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type Offers struct {
	Location string
    Section  string
	Name     string
	URL      string
	Date     string
	Meal     string
	Seats    int
    MealSort int
    DateUX   int64
	Time     []string
    LastUpd  int64
}
type URLs struct {
	TmplName string
	Vars     map[string]string
}

var tmpls *template.Template
var tmplIndex map[string]URLs
var xlatLoc map[string]string
var mealVal map[string]int
var offersFile string

func checkError(e error) {
	if e != nil {
		log.Printf("Error: %s", e)
	}
}

func loadTemplates(tmplDir string) {
	rtn, err := template.ParseGlob(tmplDir)
	checkError(err)
	for _, i := range rtn.Templates() {
		log.Printf("Template: %s", i.Name())
	}
	tmpls = rtn
}

func filenameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func getTemplate(s string) (string, error) {
	for _, i := range tmpls.Templates() {
		v := i.Name()
		if strings.TrimSuffix(v, ".tmpl") == s {
			return v, nil
		}
	}
	return "", errors.New("Not found")
}

func getOffers(s string) ([]byte, error) {
	var tmpData struct {
		Data []Offers `json:"data"`
	}

    log.Printf("getOffers(\"%s\")", s)
	for _, offer := range currentOffers() {
        for _, date := range offer.GroupByDate() {
            for _, meal := range offer.MealsByDate(date) {
                for _, seats := range offer.SeatsByMeal(date, meal) {
                    var t Offers
	                t.Location = offer.RestaurantLocation(0)
                    t.Section = offer.RestaurantLocation(1)
		            t.Name = offer.RestaurantName()
		            t.URL = offer.RestaurantURL()
		            t.Date = date.Format("02 Jan 2006")
                    t.DateUX = date.Unix()
	        	    t.Meal = meal
		            t.Seats = seats
                    t.MealSort = mealVal[meal]+seats
			        t.Time = offer.TimesByMealDate(date, meal, seats)
		            tmpData.Data = append(tmpData.Data, t)
                }
            }
        }
	}

	return json.MarshalIndent(tmpData, "", " ")
}

func handleJSON(page string, w http.ResponseWriter, r *http.Request) {
	var j []byte
	var err error

	switch page {
	case "offers.json":
		j, err = getOffers(offersFile)
	case "update":
		j = offersTimestamp()
	default:
		log.Printf("Could not find %s", page)
		return
	}
	if err != nil {
		log.Printf("No JSON %s", page)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(j)
}

func handlePage(request string, w http.ResponseWriter, r *http.Request) {
	req, ok := tmplIndex[request]
	if !ok {
		return
	}
	if tmpls.Lookup(req.TmplName) == nil {
		log.Printf("Template %s not found, (%s)", request, req.TmplName)
		return
	}
	log.Printf("Requested %s maps to %s", request, req.TmplName)
	tmpls.ExecuteTemplate(w, req.TmplName, &req.Vars)
}

func handler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.String()
	urlQuery := r.URL.Query()
	if urlPath == "/favicon.ico" {
		return
	}
	log.Printf("Request from %s: url=%v query=%v\n", r.Host, urlPath, urlQuery)
	if v, found := urlQuery["page"]; found {
		handlePage(v[0], w, r)
		return
	}
	if v, found := urlQuery["api"]; found {
		handleJSON(v[0], w, r)
		return
	}
	handlePage("default", w, r)
}

func main() {
	// Read the config file
    config.ReadConfig("config.ini")

	// Enable/Disable Timestamps
	if !config.TimestampsEnabled() {
		log.SetFlags(0)
	}

	// info
	displayBuildInfo()

	offersFile = config.OffersFilename()
    if offersFile == "" {
        log.Fatal("No offers file in configfile.")
    }

	listen := config.GetWebListen()
	loadTemplates(config.GetWebTmpls())

	tmplIndex = make(map[string]URLs)
	for _, p := range config.GetWebTmplList() {
		tmplIndex[p] = URLs{
			TmplName: config.GetWebTmplFile(p),
			Vars:     config.GetWebTmplVars(p),
		}
        log.Printf("%s: %q", p, tmplIndex[p])
	}

	xlatLoc = config.GetWebLocationTranslate()

    startWatcher(offersFile)

	http.HandleFunc("/", handler)
	log.Printf("Starting Web Server, %s", listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}

func init() {
    mealVal = make(map[string]int)
    mealVal["Breakfast"] = 100
    mealVal["Brunch"] = 200
    mealVal["Lunch"] = 300
    mealVal["Dinner"] = 400
}

// vim: noai:ts=4:sw=4:set expandtab:
