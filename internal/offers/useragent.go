package offers

var userAgent string

func SetUserAgent(ua string) {
   userAgent = ua
}

func getUserAgent() string {
   if userAgent == "" {
      return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"
   }
   return userAgent
}

// vim: noai:ts=3:sw=3:set expandtab:
