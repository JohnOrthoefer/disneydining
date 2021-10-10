package main

import (
    "html/template"
    "os"
    "strconv"
    "log"
    "errors"
	"net/http"
    "path/filepath"
    "strings"
    "encoding/json"
    "disneydining/internal/offers"
    "gopkg.in/ini.v1"
)


type Offers struct {
    Location string
    Name    string
    URL     string
    Date    string
    Meal    string
    Time    []string
}
var tmpls *template.Template
var xlatLoc map[string]string
var offersFile string = "/tmp/dining/offers.json"

func checkError(e error) {
    if e != nil {
        log.Printf("Error: %s", e)
    }
}

func loadTemplates(tmplDir string) {
    rtn, err := template.ParseGlob(tmplDir)
    checkError(err)
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

    j := offers.LoadOffers(s)
    for _, i:=range j {
        for _, offer:=range i {
            var t Offers
            t.Location = offer.Loc
            if v, found := xlatLoc[offer.Loc]; found { t.Location = v }
            t.Name = offer.Name
            t.URL  = offer.URL
            t.Date = offer.Avail[0].When.Format("02 Jan 2006")
            t.Meal = offer.Meal
            for _, tm := range offer.Avail {
                t.Time = append(t.Time, tm.When.Format("03:04 PM"))
            }
            tmpData.Data = append(tmpData.Data, t)
        }
    }
    
    return json.MarshalIndent(tmpData, "", " ")
}

func offersTimestamp(s string) ([]byte, error) {
    var tmpData struct {
        OffersSize string
        OffersTime int64 
    }

    t, err := os.Stat(s)
    if err != nil {
        return nil, err
    }

    tmpData.OffersSize = strconv.FormatInt(t.Size(), 10)
    tmpData.OffersTime = t.ModTime().Unix()
    return json.MarshalIndent(tmpData, "", " ")
}

func handleJSON(page string, w http.ResponseWriter, r *http.Request) {
    var j []byte
    var err error

    switch page {
        case "offers.json":
            j, err = getOffers(offersFile)
        case "update":
            j, err = offersTimestamp(offersFile)
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

func handlePage(page string, w http.ResponseWriter, r *http.Request) {
    type webVars struct {
        PageTitle string
    }

    t, err := getTemplate(page)
    if err != nil {
        log.Printf("Template %s not found", page)
        return
    }
    tmpls.ExecuteTemplate(w, t, &webVars {
        PageTitle: "Dining Availablity",
    })
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

    webcfg := cfg.Section("webserver")
    loadTemplates(webcfg.Key("templates").MustString("templates/*.tmpl"))

    xlatLoc = make(map[string]string)
    for _, i := range cfg.Section("locations").Keys() {
        xlatLoc[i.Name()] = i.String()
    }

    http.HandleFunc("/", handler)
    listen := webcfg.Key("listen").MustString(":8099")
    log.Printf("Starting Web Server, %s", listen)
    log.Fatal(http.ListenAndServe(listen, nil))
}

// vim: noai:ts=4:sw=4:set expandtab:
