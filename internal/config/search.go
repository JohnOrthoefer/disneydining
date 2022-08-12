package config

import (
   "log"
   "gopkg.in/ini.v1"
   "github.com/google/uuid"
   "strings"
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

func (s SearchCursor) SearchAfter() string {
   return s.section.Key("after").MustString("6:00")
}

func (s SearchCursor) SearchBefore() string {
   return s.section.Key("before").MustString("23:59")
}

func (s SearchCursor) SearchSize() string {
   return s.section.Key("size").String()
}

func (s SearchCursor) RestaurantList() []string {
   var rtn []string
   for _, i := range s.section.Key("restaurants").Strings(",") {
      rtn = append(rtn, strings.TrimSpace(i))
   }
   return rtn
}

func (s SearchCursor) UserToken() string {
   return s.section.Key("usertoken").String()
}

func (s SearchCursor) KeyString(k string) string {
   return s.section.Key(k).String()
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
   defEnabled := cfg.Section("DEFAULT").Key("enabled").MustString("true")
   defLocs := ""
   if cfg.Section("DEFAULT").HasKey("restaurants") {
      defLocs = cfg.Section("DEFAULT").Key("restaurants").String()
   }
   defNotify := ""
   if cfg.Section("DEFAULT").HasKey("notify") {
      defNotify = cfg.Section("DEFAULT").Key("notify").String()
   }
   defUsertoken := ""
   if cfg.Section("DEFAULT").HasKey("usertoken") {
      defUsertoken = cfg.Section("DEFAULT").Key("usertoken").String()
   }

   // update if entries need a UUID
   updatedFile := false
   for _, j := range cfg.Sections() {
      if j.Name() == "DEFAULT" {
         continue
      }
      if !j.HasKey("uuid") {
         j.NewKey("uuid", uuid.New().String())
         updatedFile = true
      }
   }
   // Updated?
   if updatedFile {
      if err := cfg.SaveTo(sf); err != nil {
         log.Printf("%s: error saving %s", sf, err)
      }
   }

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
      if !j.HasKey("restaurants") &&
         defLocs != "" {
         j.NewKey("restaurants", defLocs)
      }
      if !j.HasKey("notify") &&
         defNotify != "" {
         j.NewKey("notify", defNotify)
      }
      if !j.HasKey("usertoken") &&
         defUsertoken != "" {
         j.NewKey("usertoken", defUsertoken)
      }
      if !j.HasKey("enabled") {
         j.NewKey("enabled", defEnabled)
      }
   }

   // store the file
   searchFile = append(searchFile, cfg)
   log.Printf("Schedule file: %s", sf)
}


// vim: noai:ts=3:sw=3:set expandtab:
