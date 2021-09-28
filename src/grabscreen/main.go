// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
   "os"
	"context"
   "time"
	"io/ioutil"
	"log"
	"github.com/chromedp/chromedp"
   "github.com/chromedp/cdproto/dom" 
)

func main() {
   date := "11/21/2021"
   meal := "Lunch"
   size := "3"
	opts := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent( "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0"),
		chromedp.WindowSize(1920, 1080),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context
	ctx, cancel = chromedp.NewContext(
      ctx,
		//context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// capture screenshot of an element
	var buf []byte
   var res string

   log.Print("Navigate")
	chromedp.Run(ctx, chromedp.Navigate(`https://disneyworld.disney.go.com/dining`))
   
   log.Print("Wait")
	chromedp.Run(ctx, chromedp.WaitVisible(`#pageContainer > div.pepGlobalFooter`, chromedp.ByID))

   log.Print("Type")
	chromedp.Run(ctx, 
      chromedp.SetValue(`#diningAvailabilityForm-searchDate`, date, chromedp.ByID),
      chromedp.SendKeys(`#searchTime-wrapper > div.select-toggle.hoverable`, meal+`\r`, chromedp.ByID),
      chromedp.SendKeys(`#partySize-wrapper > div.select-toggle.hoverable`, size+`\r`, chromedp.ByID),
   )

   for {
      toCTX, toCancel := context.WithTimeout(ctx, 1*time.Minute)
      defer toCancel()

      log.Print("Click")
   	chromedp.Run(toCTX, chromedp.Click(`#dineAvailSearchButton > span > span > span`, chromedp.ByID))
   
      log.Print("Wait")
	   err := chromedp.Run(toCTX, chromedp.WaitVisible(`#finderList > div > h3:nth-child(2)`, chromedp.ByID))
      if err != nil {
         log.Print("Timeout Fired")
         chromedp.Run(ctx, 
            chromedp.Reload(),
	         chromedp.WaitVisible(`#pageContainer > div.pepGlobalFooter`, chromedp.ByID),
         )
         continue
      }
      break
   }

   log.Print("Write Soruce")
   chromedp.Run(ctx, 
      chromedp.ActionFunc(func(ctx context.Context) error {                          
      node, err := dom.GetDocument().Do(ctx)                                       
      if err != nil {                                                              
        return err                                                                 
      }                                                                            
      res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)                
      return err                                                                   
    }))

    f, _ := os.Create(`out.html`)
    defer f.Close()
    f.WriteString(res)

//   log.Print("ScreenShot")
//	chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 90))

	if err := ioutil.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote elementScreenshot.png and fullScreenshot.png")
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Reset
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}
