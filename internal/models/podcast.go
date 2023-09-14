package models

import "time"

type Podcast struct {
	Feed        string
	Title       string
	Description string
	Duration    string
	PublishedAt time.Time
}
