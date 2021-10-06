package mail
  
import (
   "time"
   "bytes"
   "log"
   "fmt"
   "strings"
   "net/smtp"
   "disneydining/internal/offers"
)

type mailNotify struct {
   server   string
   from     string
   to       map[string]string
}

var mailStore mailNotify

func SetSMTPHost(host string) {
   log.Printf("SMTP host: %s", host)
   mailStore.server = host
}

func SetFromAddr(from string) {
   log.Printf("From: %s", from)
   mailStore.from = from
}

func SetToAddr(n, v string) {
   log.Printf("mail- %s: %s", n, v)
   mailStore.to[n] = v
}

func MakeMsg(loc, u string, when []offers.AvailStruct) string {
   date := when[0].When.Format("Mon, Jan _2 2006")
   var times []string

   for _, t := range when {
      times = append(times, t.When.Format(time.Kitchen))
   }

   return fmt.Sprintf("Found %s on %s @%s.\n%s", loc, date, strings.Join(times, " "), u)
}

func Notify(n, msg string) {
   if mailStore.server == "" { log.Fatal("No SMTP server defined") }
   if mailStore.from == "" { log.Fatal("No From Address defined") }
   if mailStore.to[n] == "" { 
      log.Printf("No To, %s, Address defined", n) 
      return
   }

   c, err := smtp.Dial(mailStore.server+":smtp")
   if err != nil {
      log.Fatalf( "Dial:! %q", err)
   }
   defer c.Close()
   
   c.Mail(mailStore.from)
   c.Rcpt(mailStore.to[n])
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

func init() {
   mailStore.to = make(map[string]string)
}
