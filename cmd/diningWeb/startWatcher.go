package main

import (
   "log"
   "github.com/fsnotify/fsnotify"
   "disneydining/internal/offers"
)

var savedOffersFile string
var savedOffers offers.DiningMap

func setOffersFile(f string) {
   savedOffersFile = f
}

func getOffersFile() string {
   return savedOffersFile
}

func saveOffers() {
    log.Printf("saving Offers: %s", getOffersFile())
   savedOffers = offers.NewOffers()
   savedOffers.LoadOffers(getOffersFile())
}

func currentOffers() offers.DiningMap {
   return savedOffers
}

func reloadOffers(w *fsnotify.Watcher) {
   for {
      select {
      case event, ok := <-w.Events:
         if !ok { 
            return
         }
         log.Printf("Event(%d): %s", event.Op, event)
         if (event.Op == fsnotify.Remove) ||
            (event.Op == fsnotify.Rename) {
            log.Printf("Reloading %s saved offers", event.Name)
            // reset watcher to look at the new file
            w.Remove(event.Name)
            w.Add(event.Name)
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
   //defer watcher.Close()

   go reloadOffers(watcher)

   err = watcher.Add(f)
   if err != nil {
      log.Fatal("Watcher Add: ", err)
   }
}

// vim: noai:ts=4:sw=4:set expandtab:
