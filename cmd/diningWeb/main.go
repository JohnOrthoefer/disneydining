package main

import (
    "html/template"
    "log"
    "errors"
	"net/http"
    "path/filepath"
    "strings"
    "encoding/json"
    "disneydining/internal/offers"
    "gopkg.in/ini.v1"
)

const offersFile = "/tmp/dining/offers.json"

type Offers struct {
    Location string
    Name    string
    URL     string
    Date    string
    Time    []string
}
var tmpls *template.Template
var xlatLoc map[string]string

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

func getJson(s string) (string, error) {
    var tmpData []Offers

    j := offers.LoadOffers(offersFile)
    for _, i:=range j {
        for _, offer:=range i {
            var t Offers
            t.Location = offer.Loc
            if v, found := xlatLoc[offer.Loc]; found { t.Location = v }
            t.Name = offer.Name
            t.URL  = offer.URL
            t.Date = offer.Avail[0].When.Format("02 Jan 2006")
            for _, tm := range offer.Avail {
                t.Time = append(t.Time, tm.When.Format("03:04 PM"))
            }
            tmpData = append(tmpData, t)
        }
    }
    
    data, err := json.MarshalIndent(tmpData, "", " ")
	return string(data), err
}

func handler(w http.ResponseWriter, r *http.Request) {
    type webVars struct {
        PageTitle string
    }
	urlPath := r.URL.String()
	urlQuery := r.URL.Query()
    if urlPath == "/favicon.ico" {
        return
    }
	log.Printf("Request from %s: url=%v query=%v\n", r.Host, urlPath, urlQuery)
    if v, found := urlQuery["page"]; found {
        t, err := getTemplate(v[0])
        if err == nil {
            tmpls.ExecuteTemplate(w, t, &webVars {
               PageTitle: "Dining Availablity",
            })
        } else {
            log.Printf("Template %s not found", v[0])
        }
        return
    }
    if v, found := urlQuery["api"]; found {
        j, err := getJson(v[0])
        if err == nil {
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte("{\"data\":" + j + "}"))
        } else {
            log.Printf("No JSON %s", v[0])
        }
        return
    }
}
   
func main() {
    // Read the config file
    cfg, err := ini.Load("config.ini")
    if err != nil {
        log.Fatal("Failed to read config.ini file")
    }

    // Enable/Disable Timestamps
    if !cfg.Section("DEFAULT").Key("timestamps").MustBool(true) {
        log.SetFlags(0)
    }

    // info
    displayBuildInfo()

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
