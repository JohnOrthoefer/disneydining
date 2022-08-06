package main

import (
	"encoding/json"
	"errors"
	"gopkg.in/ini.v1"
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
var offersFile string = "/tmp/dining/offers.json"
var scheduleFile string

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

func getSearches() ([]byte, error) {
    type KeyValue map[string]string
    type Section struct {
        Name string
        Key     KeyValue
    }
    var tmpData struct {
        Data []*Section `json:"data"`
    }
    

    dining, _ := ini.Load(scheduleFile)
    for _, section := range dining.Sections() {
        sec := new(Section)
        sec.Name = section.Name()
        sec.Key = make(KeyValue)
        for _, key := range section.Keys() {
            kn := key.Name()
            kv := key.String()
            sec.Key[kn] = kv
        }
        tmpData.Data = append(tmpData.Data, sec)
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
    case "search":
        j, err = getSearches()
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
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal("Failed to read config.ini file")
	}

	iniDefaults := cfg.Section("DEFAULT")

	// Enable/Disable Timestamps
	if !iniDefaults.Key("timestamps").MustBool(true) {
		log.SetFlags(0)
	}

	// info
	displayBuildInfo()

	offersFile = iniDefaults.Key("saveoffers").MustString(offersFile)
	scheduleFile = cfg.Section("DEFAULT").Key("searchfile").MustString("./dining.ini")

	webcfg := cfg.Section("webserver")
	listen := webcfg.Key("listen").MustString(":8099")
	loadTemplates(webcfg.Key("tmplfiles").MustString("templates/*.tmpl"))

	tmplcfg := cfg.Section("templates")
	tmplIndex = make(map[string]URLs)

	for _, p := range tmplcfg.Keys() {
		tp, err := cfg.GetSection(p.String())
		if err != nil {
			continue
		}
		tmplEnt := URLs{
			TmplName: tp.Key("template").String(),
			Vars:     make(map[string]string),
		}
		for _, t := range tp.Keys() {
			tmplEnt.Vars[t.Name()] = t.String()
		}
		tmplIndex[p.Name()] = tmplEnt
	}

	xlatLoc = make(map[string]string)
	for _, i := range cfg.Section("locations").Keys() {
		xlatLoc[i.Name()] = i.String()
	}

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
