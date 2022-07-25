package offers

import (
   "log"
)

func checkErr(err error) {
   if err != nil {
      log.Fatalln(err)
   }
}

// vim: noai:ts=3:sw=3:set expandtab:

