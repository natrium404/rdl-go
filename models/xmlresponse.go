package models

import "encoding/xml"

type MPD struct {
	XMLName xml.Name `xml:"MPD"`
	Period  Period   `xml:"Period"`
}

type Period struct {
	AdaptationSet []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
	Representation []Representation `xml:"Representation"`
}

type Representation struct {
	Codecs          string `xml:"codecs,attr"`
	MimeType        string `xml:"mimeType,attr"`
	FBContentLength string `xml:"FBContentLength,attr,omitempty"`
	Width           string `xml:"width,attr,omitempty"`
	Height          string `xml:"height,attr,omitempty"`
	FBQualityLabel  string `xml:"FBQualityLabel,attr,omitempty"`
	BaseURL         string `xml:"BaseURL"`
}
