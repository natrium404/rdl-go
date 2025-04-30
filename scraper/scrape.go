package scraper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type Scrape struct{}

func scrapeFromURL(url string) (string, bool) {
	ctx, cancle := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancle()

	Log("Initializing the browser...")
	chromeCtx, cancle := chromedp.NewContext(ctx)
	defer cancle()

	Log("Checking the reel...")
	Log("Wait!! This might take a moment...")
	var res string
	err := chromedp.Run(chromeCtx, chromedp.Navigate(url),
		chromedp.WaitReady("script"),
	)
	if err != nil {
		Log(fmt.Sprintf("Ahh... Couldn't watch the reel.\nError: %s", err))
		return "", false
	}

	res, err = getTag(chromeCtx, "video_dash_manifest")
	if err != nil {
		Log(fmt.Sprintf("OH NO... Couldn't watch the reel.\nError: %s", err))
		return "", false
	}

	Log("UMMM... Interesting reel.")
	Log("Closing the browser...")
	return res, true
}

// This will find that exact tag that needed
func getTag(ctx context.Context, target string) (string, error) {
	var tagContents []string
	err := chromedp.Run(ctx, chromedp.Evaluate(`
					Array.from(document.querySelectorAll('script')).map(s => s.innerText)
					`, &tagContents))
	if err != nil {
		return "", err
	}

	for _, content := range tagContents {
		if strings.Contains(content, target) {
			return content, nil
		}
	}

	return "", fmt.Errorf("404 - Reel Not Found :(")
}
