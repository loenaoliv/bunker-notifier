package channels

import (
	"bunker-notifier/internal/models"

	"go.uber.org/zap"
)

func NewDiscordChannel(token string, channelID string, logger *zap.Logger) Channel {
	return &discordChannelImpl{
		token:     token,
		channelID: channelID,
		logger:    logger,
	}
}

type discordChannelImpl struct {
	token     string
	channelID string
	logger    *zap.Logger
}

func (t *discordChannelImpl) Notify(string) error {
	return nil
}

func (t *discordChannelImpl) NotifyPodcast(podcast models.Podcast) error {
	t.logger.Info("creating notification",
		zap.String("channel", "discord"),
		zap.String("title", podcast.Title),
		zap.String("description", podcast.Description),
		zap.String("duration", podcast.Duration),
		zap.Time("publishedAt", podcast.PublishedAt),
		zap.String("channelId", t.channelID),
	)
	return nil
}

func (t *discordChannelImpl) NotifyTweet(tweet models.Tweet) error {
	t.logger.Info("creating notification",
		zap.String("channel", "twitter"),
		zap.String("user", tweet.Author),
		zap.String("content", tweet.Content),
		zap.Time("publishedAt", tweet.PublishedAt),
		zap.String("channelId", t.channelID),
	)
	return nil
}
