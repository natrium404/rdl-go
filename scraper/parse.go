package scraper

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"rdl/models"
)

type ParseResponse struct {
	Success bool
	Message string
	Data    models.VideoData
}

func parsePage(content string, reelID string) ParseResponse {
	Log("Parsing the data...")
	var scriptData models.Root
	json.Unmarshal([]byte(content), &scriptData)

	var videoData models.VideoData

	medias := scriptData.Require[0][3][0].Bbox.Require[0][3][1].Bbox.Result.Data.Api.Edges
	var mediaInfo models.Media
	for _, video := range medias {
		if video.Node.Media.Code == reelID {
			mediaInfo = video.Node.Media
			break
		}
	}

	if mediaInfo.Code != reelID {
		Log("Reel not found. This reel may be redirecting or censored or need authentication.")
		return ParseResponse{
			Success: false,
			Message: "Reel not found. This reel may be redirecting. Check the reel in another tab.",
			Data:    models.VideoData{},
		}
	}

	mediaDash := mediaInfo.VideoDash
	var mediaDashData models.MPD
	err := xml.Unmarshal([]byte(mediaDash), &mediaDashData)
	if err != nil {
		log.Fatal(err)
	}

	videos := mediaDashData.Period.AdaptationSet[0].Representation
	audio := mediaDashData.Period.AdaptationSet[1].Representation[0]

	videoData.Caption = mediaInfo.Caption.Text

	videoData.Audio = models.Audio{
		MimeType: audio.MimeType,
		Codecs:   audio.Codecs,
		URL:      audio.BaseURL,
	}

	videoData.Videos = make([]models.Video, len(videos))
	for i := range videos {
		videoData.Videos[i] = models.Video{
			URL:      videos[i].BaseURL,
			Codecs:   videos[i].Codecs,
			MimeType: videos[i].MimeType,
			Height:   videos[i].Height,
			Width:    videos[i].Width,
			Quality:  videos[i].FBQualityLabel,
		}
	}

	videoData.Reel = mediaInfo.VideoVersions[0].URL
	videoData.Thumbnail = mediaInfo.ImageVersions.Candidates[0].URL
	videoData.User = models.User{
		Username: mediaInfo.User.Username,
		Profile:  mediaInfo.User.ProfilePicURL,
	}
	videoData.Code = mediaInfo.Code

	return ParseResponse{
		Success: true,
		Message: "Found the reel.",
		Data:    videoData,
	}
}
