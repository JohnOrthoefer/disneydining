package config

import (
   "log"
   "gopkg.in/ini.v1"
)

type SearchCursor struct {
   section *ini.Section
}
var searchFile []*ini.File

func (s SearchCursor) SearchName() string {
   return s.section.Name()
}

func (s SearchCursor) SearchDate() string {
   return s.section.Key("date").String()
}

func (s SearchCursor) SearchTime() string {
   return s.section.Key("time").String()
}

func (s SearchCursor) SearchSize() string {
   return s.section.Key("size").String()
}

func DiningQueries() []SearchCursor {
   var rtn []SearchCursor
   for _, i := range searchFile {
      for _, j := range i.Sections() {
         if j.Name() == "DEFAULT" ||
            !j.Key("enabled").MustBool(defaultEnable) || 
            !j.HasKey("date") {
            continue
         }
         rtn = append(rtn, SearchCursor {
            section: j,
            })  
      }
   }
   return rtn
}

func readSearchFile(sf string) {
   // Load the file
   cfg, err := ini.Load(sf)
   if err != nil {
      log.Fatal("Failed to read %s file", sf)
   }

   // Apply defaults
   defSize := cfg.Section("DEFAULT").Key("size").MustString("2")
   defTime := cfg.Section("DEFAULT").Key("time").MustString("Lunch")

   for _, j := range cfg.Sections() {
      if j.Name() == "DEFAULT" {
         continue
      }
      if !j.HasKey("size") {
         j.NewKey("size", defSize)
      }
      if !j.HasKey("time") {
         j.NewKey("time", defTime)
      }
   }

   // store the file
   searchFile = append(searchFile, cfg)
   log.Printf("Schedule file: %s", sf)
}


// vim: noai:ts=3:sw=3:set expandtab: