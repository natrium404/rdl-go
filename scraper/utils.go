package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"rdl/models"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type BrowserInfo struct {
	Name string
	Path string
}

type ProgressBar struct {
	total int64
	done  int64
}

var timer *time.Timer

func isValidURL(rawURL string) (string, bool) {
	// Verify domain name
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", false
	}

	validDomains := []string{
		"www.instagram.com",
		"instagram.com",
	}

	domainValid := slices.Contains(validDomains, parsedURL.Host)
	if !domainValid {
		return "", false
	}

	// Check path starts with /reel/ or /p/ or /username/ or /reels/
	path := strings.Trim(parsedURL.Path, "/")
	pathParts := strings.Split(path, "/")

	var reelID string

	// Validate and extract from different URL formats
	switch len(pathParts) {
	case 2:
		// Format: /reel/{reelID}/ or /reels/{reelID}/ or /p/{reelID}/
		if !slices.Contains([]string{"reel", "reels", "p"}, pathParts[0]) {
			return "", false
		}
		reelID = pathParts[1]
	case 3:
		// Format: /{username}/reel/{reelID}/
		if pathParts[0] == "" {
			return "", false
		}
		if !validUsername(pathParts[0]) && pathParts[1] != "reel" {
			return "", false
		}
		reelID = pathParts[2]
	default:
		return "", false
	}

	// Validate Reel ID format
	if !validReelID(reelID) {
		return "", false
	}
	return reelID, true
}

func validUsername(username string) bool {
	usernameRegex := regexp.MustCompile(`^[A-Za-z0-9._]{1,30}$`)
	if !usernameRegex.MatchString(username) {
		return false
	}
	return true
}

func validReelID(reelID string) bool {
	reelIDRegex := regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)
	if !reelIDRegex.MatchString(reelID) {
		return false
	}
	return true
}

// Find chromium based browsers
func findBrowser() BrowserInfo {
	browserPaths := []struct {
		name  string
		paths []string
	}{
		{
			"Google Chrome", []string{
				"/usr/bin/google-chrome",
				"/usr/bin/google-chrome-stable",
				`C:\Program Files\Google\Chrome\Application\chrome.exe`,
				`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
			},
		},
	}

	for _, browser := range browserPaths {
		for _, path := range browser.paths {
			if _, err := os.Stat(path); err == nil {
				return BrowserInfo{
					Name: browser.name,
					Path: path,
				}
			}
		}
	}

	// check for executables
	execuables := []string{
		"google-chrome", "google-chrome-stable",
		"chromium", "chromium-browser",
		"chrome",
	}

	for _, name := range execuables {
		if path, err := exec.LookPath(name); err == nil {
			return BrowserInfo{
				Name: name,
				Path: path,
			}
		}
	}

	// check for chrome-headless-shell
	return findDownloadedBrowser()
}

// Find downlaoded browser shell
func findDownloadedBrowser() BrowserInfo {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return BrowserInfo{}
	}

	platform := getPlatform()
	shellPath := getExecutable(filepath.Join(cacheDir, "rdl", fmt.Sprintf("chrome-headless-shell-%s", platform)))
	if _, err = os.Stat(shellPath); err == nil {
		return BrowserInfo{
			Name: "Chrome",
			Path: shellPath,
		}
	}
	return BrowserInfo{}
}

// Monitor download progress
func downloadBrowser(dir string) (string, error) {
	chromeJsonURL := "https://googlechromelabs.github.io/chrome-for-testing/last-known-good-versions-with-downloads.json"
	res, err := http.Get(chromeJsonURL)
	if err != nil {
		Log(fmt.Sprintf("Oops... Something went wrong. Error: %s", err))
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		Log(fmt.Sprintf("Oops... Can't read the response. Error: %s", err))
		return "", err
	}

	var chromeJsonResponse models.ChromeJSONResponse
	err = json.Unmarshal([]byte(body), &chromeJsonResponse)
	if err != nil {
		return "", err
	}

	platform := getPlatform()
	headlessShells := chromeJsonResponse.Channels.Stable.Downloads.Headless

	var downloadURL string
	for _, shell := range headlessShells {
		if strings.ToLower(shell.Platform) == strings.ToLower(platform) {
			downloadURL = shell.URL
		}
	}

	// Downlaod file from url
	filename, err := downloadFile(downloadURL, dir)
	if err != nil {
		Log(fmt.Sprintf("Oops... Something wrong. Error: %s", err))
		return "", err
	}

	Log("DOWNLOAD DONE")

	return filename, nil
}

func gerChromeOptions(browserPath string) []chromedp.ExecAllocatorOption {
	// Chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-translate", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-component-update", true),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("safebrowsing-disable-auto-update", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("incognito", true),
		chromedp.Flag("remote-debugging-port", "0"),
		chromedp.Flag("no-sandbox", false),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("single-process", false),

		chromedp.Flag("disable-logging", true),
		chromedp.Flag("log-level", "3"),
		chromedp.Flag("silent", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-crash-reporter", true),

		chromedp.ExecPath(browserPath),
	)
	return opts
}

// Get the os platform
func getPlatform() string {
	system := runtime.GOOS
	arch := runtime.GOARCH

	switch system {
	case "darwin":
		system = "mac"
	case "windows":
		system = "win"
	}

	switch arch {
	case "amd64":
		arch = "64"
	case "386":
		arch = "32"
	case "arm64":
		arch = "arm64"
	}

	if system == "mac" {
		if arch == "64" {
			return system + "-x" + arch
		} else {
			return system + "-" + arch
		}
	}

	return system + arch
}

// Download file with progress bar
func downloadFile(downloadURL, downloadPath string) (string, error) {
	res, err := http.Get(downloadURL)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Bad Status - %s", res.Status)
	}

	parseURL, err := url.Parse(downloadURL)
	if err != nil {
		return "", err
	}

	filename := path.Base(parseURL.Path)
	err = os.MkdirAll(downloadPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	output, err := os.Create(filepath.Join(downloadPath, filename))
	if err != nil {
		return "", err
	}

	defer output.Close()

	progressWrite := &ProgressBar{total: res.ContentLength}
	_, err = io.Copy(io.MultiWriter(output, progressWrite), res.Body)
	return filename, err
}

// Progress bar
func (pb *ProgressBar) Write(b []byte) (int, error) {
	byteLength := len(b)
	pb.done += int64(byteLength)

	percentage := float64(pb.done) / float64(pb.total) * 100
	ProgressLog(fmt.Sprintf("Downloading chrome-headless-shell... \n[%.2f%%]", percentage))
	return byteLength, nil
}

func getExecutable(srcDir string) string {
	system := runtime.GOOS
	var suffix string
	switch system {
	case "windows":
		suffix = ".exe"
	case "linux":
		suffix = ""
	}

	execPath := filepath.Join(srcDir, "chrome-headless-shell"+suffix)

	return execPath
}
