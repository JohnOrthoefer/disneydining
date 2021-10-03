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

