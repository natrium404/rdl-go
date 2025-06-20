package scraper

import (
	"fmt"
	// "os"
	"rdl/models"
)

var (
	Logger         func(message string)
	ProgressLogger func(message string)
)

func Log(message string) {
	// Log to file
	// file, err := os.OpenFile("scraper.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	panic(err)
	// }
	// file.WriteString(message + "\n")

	if Logger != nil {
		Logger(message)
	}
}

func ProgressLog(message string) {
	if ProgressLogger != nil {
		ProgressLogger(message)
	}
}

type Response struct {
	Data    models.VideoData
	Success bool
	Message string
}

func Scraper(url string) Response {
	reelID, valid := isValidURL(url)

	if !valid {
		ProgressLog("Not a valid url.")
		return Response{
			Data:    models.VideoData{},
			Success: false,
			Message: "Not a valid url.",
		}
	}

	reelURL := fmt.Sprintf("https://instagram.com/reels/%s", reelID)
	pageContent, success := scrapeFromURL(reelURL)
	if !success {
		return Response{
			Data:    models.VideoData{},
			Success: false,
			Message: "Couldn't fetch the reel.",
		}
	}

	videoData := parsePage(pageContent, reelID)
	if !videoData.Success {
		return Response{
			Data:    models.VideoData{},
			Success: false,
			Message: videoData.Message,
		}
	}

	Log("DONE")
	return Response{
		Data:    videoData.Data,
		Success: true,
		Message: "Got the video data.",
	}
}
