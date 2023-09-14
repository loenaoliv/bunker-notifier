package notification

import (
	"bunker-notifier/internal/channels"
	twitter "bunker-notifier/internal/services"
	"fmt"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

type twitterImpl struct {
	cron           string
	channels       []channels.Channel
	twitterUser    string
	onlyVideos     bool
	twitterService twitter.TwitterService
	scheduler      *gocron.Scheduler
	logger         *zap.Logger
}

func NewTwitter(cron string, channels []channels.Channel, twitterService *twitter.TwitterService, twitterUser string, onlyVideos bool, logger *zap.Logger, scheduler *gocron.Scheduler) Notification {
	return &twitterImpl{
		cron:           cron,
		channels:       channels,
		twitterService: *twitterService,
		scheduler:      scheduler,
		twitterUser:    twitterUser,
		logger:         logger,
		onlyVideos:     onlyVideos,
	}
}

func (j *twitterImpl) Run() error {
	if len(j.cron) == 0 {
		return fmt.Errorf("no cron")
	}

	if len(j.channels) == 0 {
		return fmt.Errorf("no channels")
	}

	if j.scheduler == nil {
		return fmt.Errorf("no scheduler")
	}

	if j.twitterService != nil {
		if len(j.twitterUser) == 0 {
			return fmt.Errorf("no twittter user")
		}

		_, err := j.scheduler.Every(j.cron).Do(func() {
			shouldNotify, tweet, err := j.twitterService.ShouldNotify(j.twitterUser, j.onlyVideos)
			if err != nil {
				j.logger.Error("error fetching tweet", zap.String("twitterUser", j.twitterUser), zap.Error(err))
				return
			}
			if shouldNotify {
				for _, v := range j.channels {
					err := v.NotifyTweet(tweet)
					if err != nil {
						j.logger.Error("error on notify", zap.String("twitterUser", j.twitterUser), zap.Error(err))
					}
				}
			}
		})

		if err != nil {
			return err
		}
	}

	return nil
}
