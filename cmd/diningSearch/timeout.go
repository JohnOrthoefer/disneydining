package main

import (
   "time"
   "log"
)

var cancelTimer *time.Timer

func StartTimer(t string) {
   wait, err := time.ParseDuration(t)
   if err != nil {
      log.Printf("Timer error, %s", err)
      return
   }
   cancelTimer = time.AfterFunc(wait, func () {
      log.Fatalf("Reset Timer Fired")
   })
   log.Printf("Timer Started, %s", t)
}

func StopTimer() {
   if cancelTimer != nil {
      cancelTimer.Stop()
      log.Printf("Timer Stopped")
      return
   }
   log.Printf("No Timer to Stop")
   
}
