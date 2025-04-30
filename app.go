package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"rdl/scraper"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	scraper.Logger = func(message string) {
		runtime.EventsEmit(a.ctx, "log", message)
	}
}

// Scrape the web page from the url
func (a *App) ScrapeWebPage(url string) scraper.Response {
	response := scraper.Scraper(url)
	return response
}

// Save the file
func (a *App) SaveFile(url string, nameID string, ext string, extra string) error {
	now := time.Now()
	timestamp := now.Format("2006_01_02_15_04_05")
	saveDialog := runtime.SaveDialogOptions{
		Title:           "Save File",
		DefaultFilename: fmt.Sprintf("reel_downloader_%s_%s_%s.%s", nameID, timestamp, extra, ext),
	}

	downloadPath, err := runtime.SaveFileDialog(a.ctx, saveDialog)
	if err != nil || downloadPath == "" {
		return err
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	outputFile, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, res.Body)
	return err
}
