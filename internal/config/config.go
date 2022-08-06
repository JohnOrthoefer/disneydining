package config

import (
   "gopkg.in/ini.v1"
   "time"
   "log"
)

var configFile *ini.File

// Enable/Disable Timestamps in log
func TimestampsEnabled() bool {
   return configFile.Section("DEFAULT").Key("timestamps").MustBool(true)
}

// Get UserAgent 
func GetUserAgent() string {
   return configFile.Section("browser").Key("agent").MustString(defaultUserAgent)
}

// Get Max Runtime 
func GetRuntimeLimit() string {
   return configFile.Section("DEFAULT").Key("timeout").MustString("10m")
}

// The URL to obtain the Auth cookie
func GetAuthURL() string {
   return configFile.Section("disney").Key("AuthURL").MustString(defaultAuthURL)
}

// The cookie name that has the auth string
func GetAuthCookie() string {
   return configFile.Section("disney").Key("AuthCookie").MustString(defaultAuthCookie)
}

// the URL to append the query to
func GetQueryURL() string {
   return configFile.Section("disney").Key("url").MustString(defaultQueryURL)
}

func RetentionTime() time.Duration{
   defRet, _ := time.ParseDuration("30m")
   return configFile.Section("DEFAULT").Key("offerretention").MustDuration(defRet)
}

func OffersFilename() string {
   if configFile.Section("DEFAULT").HasKey("saveoffers") {
      return configFile.Section("DEFAULT").Key("saveoffers").String()
   }
   return ""
}

func ReadConfig(cf string) {
   // Read the config file
   cfg, err := ini.LoadSources(ini.LoadOptions{
      IgnoreInlineComment:         true,
      UnescapeValueCommentSymbols: true,
      }, cf)
   if err != nil {
      log.Fatal("ReadConfig (%s): %s", cf, err)
   }
   configFile = cfg
   if cfg.Section("DEFAULT").HasKey("searchfile") {
      readSearchFile(cfg.Section("DEFAULT").Key("searchfile").String())
   } else {
      log.Printf("No search files found.")
   }
}

// vim: noai:ts=3:sw=3:set expandtab:

