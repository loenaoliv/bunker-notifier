package services

import (
	"bunker-notifier/internal/infra/cache"
	"bunker-notifier/internal/models"
	"bunker-notifier/internal/utils"
	"fmt"

	"github.com/coocood/freecache"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
)

type PodcastService interface {
	ShouldNotify(podcastUrl string) (bool, models.Podcast, error)
}

type podcastServiceImpl struct {
	cache  cache.FreeCache
	logger *zap.Logger
}

func NewPodcastService(cache cache.FreeCache, logger *zap.Logger) PodcastService {
	return &podcastServiceImpl{
		cache:  cache,
		logger: logger,
	}
}

func (p *podcastServiceImpl) ShouldNotify(podcastUrl string) (bool, models.Podcast, error) {
	podcast := models.Podcast{}
	podcastKey := utils.GetMD5Hash(podcastUrl)

	err := p.cache.Get(podcastKey, &podcast)
	if err != nil && err != freecache.ErrNotFound {
		p.logger.Error("error getting podcast from cache", zap.Error(err))
		return false, models.Podcast{}, err
	}

	fetchedPodcast, err := fetchLatest(podcastUrl)
	if err != nil {
		p.logger.Error("error fetching podcast", zap.Error(err))
		return false, models.Podcast{}, err
	}

	if podcast.Title == "" {
		p.cache.Set(podcastKey, fetchedPodcast, -1)
		return false, fetchedPodcast, nil
	}

	if fetchedPodcast.PublishedAt.Compare(podcast.PublishedAt) > 0 {
		p.cache.Set(podcastKey, fetchedPodcast, -1)
		return true, fetchedPodcast, nil
	}

	return false, fetchedPodcast, nil
}

func fetchLatest(podcastUrl string) (models.Podcast, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(podcastUrl)

	if err != nil {
		return models.Podcast{}, err
	}

	if feed == nil {
		return models.Podcast{}, fmt.Errorf("invalid feed")
	}

	if feed.Items == nil || len(feed.Items) == 0 {
		return models.Podcast{}, fmt.Errorf("no items")
	}

	latest := feed.Items[0]
	if latest.PublishedParsed == nil {
		return models.Podcast{}, fmt.Errorf("invalid date")
	}

	duration := "UNKNOW"
	if latest.ITunesExt != nil && len(latest.ITunesExt.Duration) > 0 {
		duration = latest.ITunesExt.Duration
	}

	return models.Podcast{
		Feed:        strip.StripTags(feed.Title),
		Title:       strip.StripTags(latest.Title),
		Description: strip.StripTags(latest.Description),
		PublishedAt: *latest.PublishedParsed,
		Duration:    duration,
	}, nil
}
