package services

import (
	"bunker-notifier/internal/infra/cache"
	"bunker-notifier/internal/models"
	"bunker-notifier/internal/utils"
	"context"
	"fmt"

	"github.com/coocood/freecache"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"go.uber.org/zap"
)

type TwitterService interface {
	ShouldNotify(user string, onlyVideos bool) (bool, models.Tweet, error)
}

type twitterServiceImpl struct {
	cache   cache.FreeCache
	logger  *zap.Logger
	scraper *twitterscraper.Scraper
}

func NewTwitterService(cache cache.FreeCache, logger *zap.Logger) TwitterService {
	scraper := twitterscraper.New()
	scraper.SetSearchMode(twitterscraper.SearchVideos)
	return &twitterServiceImpl{
		cache:   cache,
		logger:  logger,
		scraper: scraper,
	}
}

func (p *twitterServiceImpl) ShouldNotify(user string, onlyVideos bool) (bool, models.Tweet, error) {
	tweet := models.Tweet{}
	tweetKey := utils.GetMD5Hash(fmt.Sprintf("twitter:%s", user))

	err := p.cache.Get(tweetKey, &tweet)
	if err != nil && err != freecache.ErrNotFound {
		p.logger.Error("error getting podcast from cache", zap.Error(err))
		return false, models.Tweet{}, err
	}

	fetchedTweet, err := p.fetchLatestTweet(user)
	if err != nil {
		p.logger.Error("error fetching tweet", zap.Error(err))
		return false, models.Tweet{}, err
	}

	if tweet.Author == "" {
		p.cache.Set(tweetKey, fetchedTweet, -1)
		return false, fetchedTweet, nil
	}

	if fetchedTweet.PublishedAt.Unix() > tweet.PublishedAt.Unix() {
		if onlyVideos && fetchedTweet.VideoURI == "" {
			p.cache.Set(tweetKey, fetchedTweet, -1)
			return false, fetchedTweet, nil
		}
		p.cache.Set(tweetKey, fetchedTweet, -1)
		return true, fetchedTweet, nil
	}

	return false, fetchedTweet, nil
}

func (p *twitterServiceImpl) fetchLatestTweet(user string) (models.Tweet, error) {
	tweet := models.Tweet{}
	for t := range p.scraper.GetTweets(context.Background(), user, 1) {
		if t.Error != nil {
			return models.Tweet{}, fmt.Errorf("error fetching last tweet")
		}
		tweet = models.Tweet{
			Author:      t.Username,
			Content:     t.Text,
			PublishedAt: t.TimeParsed,
		}
		if len(t.Videos) > 0 {
			tweet.VideoURI = t.Videos[0].URL
		}
		if len(t.Photos) > 0 {
			tweet.ImageURI = t.Photos[0].URL
		}
	}

	return tweet, nil
}
