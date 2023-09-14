package channels

import "bunker-notifier/internal/models"

type Channel interface {
	NotifyPodcast(models.Podcast) error
	Notify(string) error
	NotifyTweet(models.Tweet) error
}
