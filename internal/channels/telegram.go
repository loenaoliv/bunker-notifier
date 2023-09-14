package channels

import (
	"bunker-notifier/internal/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tele "github.com/tucnak/telebot"
	"go.uber.org/zap"
)

func NewTelegramChannel(token string, chatID int64, logger *zap.Logger) (Channel, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return &telegramChannelImpl{}, err
	}

	return &telegramChannelImpl{
		token:  token,
		chatID: chatID,
		logger: logger,
		bot:    bot,
	}, nil
}

type telegramChannelImpl struct {
	token  string
	chatID int64
	logger *zap.Logger
	bot    *tele.Bot
}

func (t *telegramChannelImpl) Notify(string) error {
	return nil
}

func (t *telegramChannelImpl) NotifyPodcast(podcast models.Podcast) error {
	t.logger.Info("creating podcast notification",
		zap.String("channel", "telegram"),
		zap.String("title", podcast.Title),
		zap.String("description", podcast.Description),
		zap.String("duration", podcast.Duration),
		zap.Time("publishedAt", podcast.PublishedAt),
		zap.Int64("channelId", t.chatID),
	)

	_, err := t.bot.Send(&tele.Chat{ID: t.chatID},
		fmt.Sprintf("<i>🎙️%s</i>\n\n<b>%s</b>\n\n%s\n\n<i>%s</i>", podcast.Feed, podcast.Title, podcast.Description, podcast.Duration),
		&tele.SendOptions{
			ParseMode: tele.ModeHTML,
		})
	if err != nil {
		return err
	}
	return nil
}

func (t *telegramChannelImpl) NotifyTweet(tweet models.Tweet) error {
	t.logger.Info("creating notification",
		zap.String("channel", "twitter"),
		zap.String("user", tweet.Author),
		zap.String("content", tweet.Content),
		zap.Time("publishedAt", tweet.PublishedAt),
		zap.Int64("channelId", t.chatID),
	)

	if tweet.VideoURI != "" {
		err := downloadFile("downloads/video.mp4", tweet.VideoURI)
		if err != nil {
			return err
		}

		a := &tele.Video{File: tele.FromDisk("downloads/video.mp4")}
		a.Caption = tweet.Content
		_, err = t.bot.Send(&tele.Chat{ID: t.chatID}, a)
		if err != nil {
			return err
		}
		return nil
	} else if tweet.ImageURI != "" {
		a := &tele.Photo{File: tele.FromURL(tweet.ImageURI)}
		a.Caption = tweet.Content
		_, err := t.bot.Send(&tele.Chat{ID: t.chatID}, a)
		if err != nil {
			return err
		}
		return nil
	}

	_, err := t.bot.Send(&tele.Chat{ID: t.chatID}, tweet.Content)
	if err != nil {
		return err
	}

	return nil
}

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
