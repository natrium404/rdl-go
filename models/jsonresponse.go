package models

type Root struct {
	Require [][][]RootBBox `json:"require"`
}

type RootBBox struct {
	Bbox BBox `json:"__bbox"`
}

type BBox struct {
	Require [][][]BBoxResult `json:"require"`
}

type BBoxResult struct {
	Bbox BBoxFinal `json:"__bbox"`
}

type MediaNode struct {
	Node struct {
		Media Media `json:"media"`
	} `json:"node"`
}

type Media struct {
	Caption struct {
		Text string `json:"text"`
	} `json:"caption"`
	Code         string `json:"code"`
	CommentCount int    `json:"comment_count"`
	HasAudio     bool   `json:"has_audio"`
	LikeCount    int    `json:"like_count"`
	ID           string `json:"id"`
	User         struct {
		Username      string `json:"username"`
		ProfilePicURL string `json:"profile_pic_url"`
	} `json:"user"`
	ImageVersions struct {
		Candidates []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"candidates"`
	} `json:"image_versions2"`
	VideoVersions []struct {
		Type int    `json:"type"`
		URL  string `json:"url"`
	} `json:"video_versions"`
	VideoDash string `json:"video_dash_manifest"`
}

type BBoxFinal struct {
	Complete bool `json:"complete"`
	Result   struct {
		Data struct {
			Api struct {
				Edges []MediaNode `json:"edges"`
			} `json:"xdt_api__v1__clips__clips_on_logged_out_connection_v2"`
		} `json:"data"`
	} `json:"result"`
}

type ChromeJSONResponse struct {
	Timestamp string
	Channels  struct {
		Stable struct {
			Downloads struct {
				Headless []struct {
					Platform string `json:"platform"`
					URL      string `json:"url"`
				} `json:"chrome-headless-shell"`
			} `json:"downloads"`
		} `json:"Stable"`
	} `json:"channels"`
}
