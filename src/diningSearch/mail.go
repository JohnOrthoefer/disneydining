package main
  
import (
   "time"
   "bytes"
   "log"
   "fmt"
   "strings"
   "net/smtp"
)

/*
type mailNotify {
   server   string
   from     string
   to       []string
}

var notifyAddrs map[string]mailNotify
*/

func MakeMsg(loc, u string, when []AvailStruct) string {
   date := when[0].When.Format("Mon, Jan _2 2006")
   var times []string

   for _, t := range when {
      times = append(times, t.When.Format(time.Kitchen))
   }

   return fmt.Sprintf("Found %s on %s @%s.\n%s", loc, date, strings.Join(times, " "), u)
}

func Notify(msg string) {
   smtpHost := "smtp.orthoefer.org"
   c, err := smtp.Dial(smtpHost+":smtp")
   if err != nil {
      log.Fatalf( "Dial:! %q", err)
   }
   defer c.Close()
   
   c.Mail("john@orthoefer.org")
   c.Rcpt("6177214121@vzwpix.com")
   wc, err := c.Data()
   if err != nil {
      log.Fatalf( "Data: %q", err)
   }
   defer wc.Close()

   buf := bytes.NewBufferString(msg)
   if _, err = buf.WriteTo(wc); err != nil {
      log.Fatalf( "Send: %q", err)
   }

   log.Printf("Success!")
}
