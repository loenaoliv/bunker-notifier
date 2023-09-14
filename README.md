# bunker-notifier
Bunker Notifier is a simple Go application that checks for updates to a list of podcasts and sends notifications when new episodes are available.

## Usage
To use Podcast Notifier, you'll need to create a configuration file that lists the podcasts you want to track. The configuration file should be a JSON file with the following format:

```
{
    "notifications": [
        {
            "cron": "1m",
            "type": "podcast",
            "podcastUrl": "PODCAST_FEED_URL",
            "channels": [
                {
                    "type": "telegram",
                    "token": "TELEGRAM_BOT_TOKEN",
                    "chatId": -10000000000 // telegram chat id
                },
                {
                    "type": "telegram",
                    "token": "TELEGRAM_BOT_TOKEN",
                    "chatId": -10000000000 // telegram chat id
                }
            ]
        },
        {
            "cron": "1m",
            "type": "podcast",
            "podcastUrl": "PODCAST_FEED_URL",
            "channels": [
                {
                    "type": "telegram",
                    "token": "TELEGRAM_BOT_TOKEN",
                    "chatId": -10000000000 // telegram chat id
                },
                {
                    "type": "telegram",
                    "token": "TELEGRAM_BOT_TOKEN",
                    "chatId": -10000000000 // telegram chat id
                }
            ]
        }
    ]
}
```

Once you have your configuration file in */config* directory, you can run the Podcast Notifier application using the following command:

```
git clone https://github.com/loenaoliv/bunker-notifier.git
cd bunker-notifier
make run
```

The application will check for updates to the podcasts listed in the configuration file and send notifications (if the notify field is set to true) when new episodes are available.

## Contributing
If you'd like to contribute to Podcast Notifier, please fork the repository and create a pull request with your changes. We welcome contributions of all kinds, including bug fixes, new features, and documentation improvements.

## License
Bunker Notifier is released under the MIT License. See LICENSE for details.
