package notification

import (
	"bunker-notifier/internal/channels"
	podcast "bunker-notifier/internal/services"
	"fmt"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

type podcastImpl struct {
	cron           string
	channels       []channels.Channel
	podcastService podcast.PodcastService
	podcastUrl     string
	scheduler      *gocron.Scheduler
	logger         *zap.Logger
}

func NewPodcast(cron string, channels []channels.Channel, podcastService *podcast.PodcastService, podcastUrl string, logger *zap.Logger, scheduler *gocron.Scheduler) Notification {
	return &podcastImpl{
		cron:           cron,
		channels:       channels,
		podcastService: *podcastService,
		scheduler:      scheduler,
		podcastUrl:     podcastUrl,
		logger:         logger,
	}
}

func (j *podcastImpl) Run() error {
	if len(j.cron) == 0 {
		return fmt.Errorf("no cron")
	}

	if len(j.channels) == 0 {
		return fmt.Errorf("no channels")
	}

	if j.scheduler == nil {
		return fmt.Errorf("no scheduler")
	}

	if j.podcastService != nil {
		if len(j.podcastUrl) == 0 {
			return fmt.Errorf("no podcast url")
		}

		_, err := j.scheduler.Every(j.cron).Do(func() {
			shouldNotify, podcast, err := j.podcastService.ShouldNotify(j.podcastUrl)
			if err != nil {
				j.logger.Error("error fetching podcast", zap.String("podcastUrl", j.podcastUrl), zap.Error(err))
				return
			}
			if shouldNotify {
				for _, v := range j.channels {
					err := v.NotifyPodcast(podcast)
					if err != nil {
						j.logger.Error("error on notify", zap.String("podcastUrl", j.podcastUrl), zap.Error(err))
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
