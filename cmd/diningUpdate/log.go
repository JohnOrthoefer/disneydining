package main

import (
	"log"
)

func setTimestamps() {
	log.SetFlags(1)
}

func clearTimestamps() {
	log.SetFlags(0)
}

// vim: noai:ts=4:sw=4:set expandtab:
