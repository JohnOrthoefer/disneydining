package main

import (
   "net/http"
   "net/url"
   "log"
)

func pushover(user, msg string, api apikey) {
/*
   log.Printf("Url: %s", api.Url)
   log.Printf("Token: %s", api.Token)
   log.Printf("User: %s", user)
   log.Printf("Msg: %s", msg)
*/

   resp, err := http.PostForm(api.Url, url.Values{
      "token":   {api.Token},
      "user":    {user},
      "message": {msg},
   })

   if err != nil {
      log.Printf("pushover error: %s", err)
      return
   }
   log.Printf("pushover resp: %q", resp.Status)
}

