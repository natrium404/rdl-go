package scraper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"rdl/models"
	"strings"

	"github.com/chromedp/chromedp"
)

type Scrape struct{}

func scrapeFromURL(url string) (string, bool) {
	// Find browser
	browser := findBrowser()
	if browser.Path != "" {
		Log(fmt.Sprintf("Chromium-based browser detected!! Using %s", browser.Name))
	} else {
		fmt.Println("No chromium-based browser found!! Downloading one might take time.")
		// Downlaod browser
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			Log(fmt.Sprintf("Can't get the cache diretory. Error: %s", err))
			return "", false
		}
		downloadDir := filepath.Join(cacheDir, models.DOWNLOAD_CACHE_DIR)
		filename, err := downloadBrowser(downloadDir)
		if err != nil {
			Log(fmt.Sprintf("Oops... Download failed. Error: %s", err))
			return "", false
		}
		err = extract(filepath.Join(downloadDir, filename), downloadDir)
		if err != nil {
			Log(fmt.Sprintf("Oops... Unziping failed. Error: %s", err))
			return "", false
		}

		browser = findDownloadedBrowser()
		Log("Download successfull...")

	}

	Log("Initializing the browser...")
	// Use the new context creation function that handles Windows popup hiding
	chromeCtx, cancel, err := createChromeContext(browser.Path)
	if err != nil {
		Log(fmt.Sprintf("Failed to create Chrome context. Error: %s", err))
		return "", false
	}
	defer cancel()

	Log("Checking the reel...")
	Log("Wait!! This might take a moment...")
	var res string
	err = chromedp.Run(chromeCtx, chromedp.Navigate(url),
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
