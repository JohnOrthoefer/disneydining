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

// vim: noai:ts=3:sw=3:set expandtab:
