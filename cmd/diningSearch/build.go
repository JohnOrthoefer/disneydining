package main

import (
   "log"
)

var (
	sha1ver   string
	buildTime string
	repoName  string
)

func displayBuildInfo() {
	log.Printf("%s: Build %s, Time %s", repoName, sha1ver, buildTime)
}
