package config

import (
	"encoding/json"
	"os"
)

const configFile = "config/config.json"

type Config struct {
	Notifications []struct {
		Cron        string `json:"cron"`
		Type        string `json:"type"`
		PodcastURL  string `json:"podcastUrl"`
		TwitterUser string `json:"twitterUser"`
		OnlyVideos  bool   `json:"onlyVideos"`
		Channels    []struct {
			Type      string `json:"type"`
			Token     string `json:"token"`
			ChannelID string `json:"channelId,omitempty"`
			ChatID    int64  `json:"chatId,omitempty"`
		} `json:"channels"`
	} `json:"notifications"`
}

func (c *Config) Load() error {
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(c)
	if err != nil {
		return err
	}

	return nil
}
