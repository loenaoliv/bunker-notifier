package models

import "time"

type Tweet struct {
	Author      string
	Content     string
	VideoURI    string
	ImageURI    string
	PublishedAt time.Time
}
