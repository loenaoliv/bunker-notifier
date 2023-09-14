package main

import (
	"bunker-notifier/internal/channels"
	"bunker-notifier/internal/infra/cache"
	"bunker-notifier/internal/infra/config"
	"bunker-notifier/internal/infra/logs"
	"bunker-notifier/internal/notification"
	"bunker-notifier/internal/services"
	"time"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger := logs.NewLogger(zapcore.InfoLevel)

	config := config.Config{}
	err := config.Load()
	if err != nil {
		logger.Panic("failed to read config", zap.Error(err))
	}

	cache := cache.NewCache(50000)
	podcastService := services.NewPodcastService(cache, logger)
	twitterService := services.NewTwitterService(cache, logger)
	s := gocron.NewScheduler(time.UTC)

	for _, n := range config.Notifications {
		chans := []channels.Channel{}

		for _, c := range n.Channels {
			switch c.Type {
			case "telegram":
				channel, err := channels.NewTelegramChannel(c.Token, c.ChatID, logger)
				if err != nil {
					logger.Panic("error creating telegram channel", zap.Error(err))
				}
				chans = append(chans, channel)
			case "discord":
				channel := channels.NewDiscordChannel(c.Token, c.ChannelID, logger)
				chans = append(chans, channel)
			}
		}

		noti := notification.NewNotification()
		switch n.Type {
		case "podcast":
			noti = notification.NewPodcast(n.Cron, chans, &podcastService, n.PodcastURL, logger, s)
		case "twitter":
			noti = notification.NewTwitter(n.Cron, chans, &twitterService, n.TwitterUser, n.OnlyVideos, logger, s)
		}

		err := noti.Run()
		if err != nil {
			logger.Panic("error creating job", zap.Error(err))
		}
	}

	logger.Info("started")
	s.StartBlocking()
}
