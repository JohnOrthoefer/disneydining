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

func checkError(e error) {
    if e != nil {
        log.Printf("Error: %s", e)
    }
}

func loadTemplates() {
    rtn, err := template.ParseGlob("templates/*.tmpl")
    checkError(err)

    for _, i := range rtn.Templates() {
        log.Printf("Loading- %s", i.Name())
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

func getJson(s string) (string, error) {
    var tmpData []Offers

    j := offers.LoadOffers(offersFile)
    log.Printf("Data Size %d", len(j))
    for _, i:=range j {
        for _, offer:=range i {
            var t Offers
            t.Location = offer.Loc
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
	urlPath := r.URL.String()
	urlQuery := r.URL.Query()
    if urlPath == "/favicon.ico" {
        return
    }
	log.Printf("Request from %s: url=%v query=%v\n", r.Host, urlPath, urlQuery)
    if v, found := urlQuery["page"]; found {
        t, err := getTemplate(v[0])
        if err == nil {
            tmpls.ExecuteTemplate(w, t, nil)
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
    if _, found := urlQuery["reload"]; found {
        loadTemplates()
        return
    }
}
   
func main() {
    loadTemplates()

    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8099", nil))
}

// vim: noai:ts=4:sw=4:set expandtab:
