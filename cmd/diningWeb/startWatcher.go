package main

import (
    "time"
    "log"
    "os"
    "strconv"
    "github.com/fsnotify/fsnotify"
    "disneydining/internal/offers"
    "path/filepath"
    "encoding/json"
)

var savedOffersFile string
var savedOffers offers.DiningMap
var watchFile string
var timeStamp []byte

func setOffersFile(f string) {
    savedOffersFile = f
}

func getOffersFile() string {
    return savedOffersFile
}

func saveOffers() {
    var tmpData struct {
        OffersSize string
        OffersTime int64
    }

    log.Printf("saving Offers: %s", getOffersFile())
    savedOffers = offers.NewOffers()
    savedOffers.LoadOffers(getOffersFile())

    t, _:= os.Stat(getOffersFile())
    tmpData.OffersSize = strconv.FormatInt(t.Size(), 10)
    tmpData.OffersTime = t.ModTime().Unix()
    timeStamp, _ = json.Marshal(tmpData)
}

func currentOffers() offers.DiningMap {
   return savedOffers
}

func offersTimestamp() []byte {
    return timeStamp
}


func reloadOffers(w *fsnotify.Watcher) {
   for {
      select {
      case event, ok := <-w.Events:
         if !ok { 
            return
         }
         // watching the directory seems to be more stable
         if (event.Op == fsnotify.Create) && 
            (event.Name == getOffersFile()) {
            log.Printf("Event(%d): %s", event.Op, event)
            time.Sleep(50 * time.Microsecond) 
            log.Printf("Reloading %s saved offers", watchFile)
            saveOffers()
         }
      case err, ok := <-w.Errors:
         if !ok {
            return
         }
         log.Printf("Watcher Error: %s", err)
      }
   }
}

func startWatcher(f string) {
   log.Printf("Watching %s", f)
   setOffersFile(f)
   saveOffers()

   watcher, err := fsnotify.NewWatcher()
   if err != nil {
      log.Fatal("Watcher New: ", err)
   }
   watchFile = f

   go reloadOffers(watcher)

   err = watcher.Add(filepath.Dir(f))
//   err = watcher.Add(f)
   if err != nil {
      log.Fatal("Watcher Add: ", err)
   }
}

// vim: noai:ts=4:sw=4:set expandtab:
