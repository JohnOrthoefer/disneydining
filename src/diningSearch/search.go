package main

import (
   "log"
   "context"
   "time"
   "github.com/chromedp/chromedp"
   "github.com/chromedp/cdproto/dom" 
)

var CTX context.Context
var Cancel context.CancelFunc

func InitContext(agent string) {
  opts := []chromedp.ExecAllocatorOption{
      chromedp.UserAgent(agent),
      chromedp.NoFirstRun,
      chromedp.NoDefaultBrowserCheck,
      chromedp.Headless,
      chromedp.DisableGPU,
   }

   CTX, Cancel = chromedp.NewExecAllocator(context.Background(), opts...)

   CTX, Cancel = chromedp.NewContext(CTX)
}

func navigate(u string) {
   log.Printf("Navigating to %s", u)
   err := chromedp.Run(CTX, chromedp.Navigate(u))
   if err != nil {
      log.Fatalf("Error, %s", err)
   }
}

func wait(w string) {
   log.Printf("Waiting for %s", w)
   chromedp.Run(CTX, chromedp.WaitVisible(w, chromedp.ByID))
}

func setSearch(date, meal, size string) {
   log.Printf("Setting Search")
   chromedp.Run(CTX, 
      chromedp.SetValue(`#diningAvailabilityForm-searchDate`, date, chromedp.ByID),
      chromedp.SendKeys(`#searchTime-wrapper > div.select-toggle.hoverable`, meal+`\r`, chromedp.ByID),
      chromedp.SendKeys(`#partySize-wrapper > div.select-toggle.hoverable`, size+`\r`, chromedp.ByID),
   )
}

func runSearch() {
   if CTX == nil { log.Fatal("Oops") }
   log.Printf("Running Search")
   for {
      toCTX, toCancel := context.WithTimeout(CTX, 1*time.Minute)
      defer toCancel()

      log.Print("Click")
      chromedp.Run(toCTX, chromedp.Click(`#dineAvailSearchButton > span > span > span`, chromedp.ByID))
   
      log.Print("Wait")
      err := chromedp.Run(toCTX, chromedp.WaitVisible(`#finderList > div > h3:nth-child(2)`, chromedp.ByID))
      if err != nil {
         log.Printf("Timeout Fired, %s", err)
         chromedp.Run(CTX,
            chromedp.Reload(),
            chromedp.WaitVisible(`#pageContainer > div.pepGlobalFooter`, chromedp.ByID),
         )
         continue
      }
      break
   }
}

func GetPage(u, d, t, p string) string {
   var res string

   if CTX == nil { log.Fatal("Oops") }
   navigate(u)
   wait("#pageContainer > div.pepGlobalFooter")
   setSearch(d, t, p)
   runSearch()

   chromedp.Run(CTX,
      chromedp.ActionFunc(func(ctx context.Context) error {
      node, err := dom.GetDocument().Do(ctx)
      if err != nil {
        return err
      }
      res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
      return err
   }))
   return res
}

