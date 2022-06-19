package main

import (
	"context"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"log"
	"time"
   "io/ioutil"
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

func dumpPage(c context.Context, fn string) {
   var res []byte
   quality := 90
   chromedp.Run(c, chromedp.FullScreenshot(&res, quality))
 	if err := ioutil.WriteFile(fn, res, 0o644); err != nil {
		log.Fatal(err)
	}
   log.Printf("Writing %s\n", fn)
}

func wait(w string) {
	log.Printf("Waiting for %s", w)
   chromedp.Run(CTX, chromedp.WaitReady(w))
/*
	for {
		toCTX, toCancel := context.WithTimeout(CTX, 3*time.Minute)
		defer toCancel()

		err := chromedp.Run(toCTX, chromedp.WaitVisible(w, chromedp.ByID))
		if err != nil {
         dumpPage(CTX);
			log.Printf("Timeout Fired, %s", err)
			chromedp.Run(CTX, chromedp.Reload())
			continue
		}
		break
	}
*/
}

func setSearch(date, meal, size string) {
	log.Printf("Setting Search date=%s meal=%s size=%s", date, meal, size)
	chromedp.Run(CTX, chromedp.SetValue(`div.date-input > finder-input > input`, date, chromedp.BySearch))
   dumpPage(CTX, "02a-searchSet.png")

	chromedp.Run(CTX, 
      chromedp.Click(`div.custom-dropdown-container`, chromedp.ByQuery),
   )
   dumpPage(CTX, "02b-searchSet.png")
	chromedp.Run(CTX, 
      chromedp.SendKeys(`#custom-dropdown-button > div.button-text-container`, meal+`\r`, chromedp.ByID),
   )

//	chromedp.Run(CTX, chromedp.SendKeys(`#partySizeCounter`, size+`\r`, chromedp.BySearch))
   dumpPage(CTX, "02c-searchSet.png")

}

func runSearch() {
	if CTX == nil {
		log.Fatal("Oops")
	}
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
			chromedp.Run(CTX, chromedp.Reload())
			wait(`finder-dine-availability-filter`)
			continue
		}
		break
	}
}

func GetPage(u string, d time.Time, t string, p string) string {
	var res string

	if CTX == nil {
		log.Fatal("Oops")
	}
	navigate(u)
	wait("finder-dine-availability-filter")
   dumpPage(CTX, "01-Navigate.png");

	setSearch(d.Format("Monday, January 2, 2006"), t, p)
   dumpPage(CTX, "02-searchSet.png")

	runSearch()
   dumpPage(CTX, "03-runSearch.png")

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

// vim: noai:ts=3:sw=3:set expandtab:
