package models

type VideoData struct {
	Caption   string
	Thumbnail string
	User      User
	Videos    []Video
	Audio     Audio
	Reel      string
	Code      string
}

type Video struct {
	MimeType string
	Codecs   string
	Quality  string
	URL      string
	Width    string
	Height   string
}

type Audio struct {
	MimeType string
	URL      string
	Codecs   string
}

type User struct {
	Username string
	Profile  string
}
