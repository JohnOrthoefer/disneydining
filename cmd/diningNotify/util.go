package main

import (
   "time"
)

func fmtDate(dt time.Time) string {
   return dt.Format("2 Jan 2006")
}

// vim: noai:ts=3:sw=3:set expandtab:
