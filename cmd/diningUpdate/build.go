package main

import (
	"log"
)

var (
	sha1ver   string
	buildTime string
	repoName  string
	goVersion  string
	goArch  string
)

func displayBuildInfo() {
   log.Printf("Complier: %s %s", goVersion, goArch)
	log.Printf("%s: Build %s, Time %s", repoName, sha1ver, buildTime)
}
