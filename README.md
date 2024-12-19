# rss-monitor

Monitor RSS feeds and post to different services.

```
./rss-monitor
Monitor RSS feeds based on the timeframe provided and send recent posts to different services

Usage:
  RSS-Monitor [command]

Available Commands:
  discord     Discord webhook notification
  help        Help about any command
  slack       Slack webhook notification

Flags:
  -d, --debug             Enable debug logging
  -h, --help              help for RSS-Monitor
  -r, --rss string        File containing the RSS links to scan
  -t, --template string   Files used as a template for the message to send
  -f, --timeframe int     Only get articles that were posted in the past number of hours (default 24)

Use "RSS-Monitor [command] --help" for more information about a command.
```

## Build

Clone this repo and then run the following commands to build for your current OS. You can also run the `make` command.
```bash
# Move into the truffleproxy directory
cd rss-monitor

# Build the binary
CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath
```

## Usage

You will first need to create a message template and a list of RSS feeds to monitor.

### Message Template

The message template 
- Date - The current date the program is ran
- Payload
    - Link - Hyperlink of the article
    - Published - Date/time the article was published
    - Title - Title of the article


Below is an example template that will include the date and then a bullet list of new feeds that include the page's title and URL link.

```
*:new: Posts from {{ .Date }} :new:*
{{ range .Posts }}
* {{ .Title }} | {{ .Link }}
{{- end }}
```

### RSS Links

You will also need a file containing each of the RSS feeds you would like to monitor. One URL per line and a `#` at the start of the line is a comment. The below example will only check `https://blog.lum8rjack.com/index.xml`.

```
# This is a comment
https://blog.lum8rjack.com/index.xml
```

### Command

The rss-monitor application takes in the following:
- message template
- rss links
- service - which service you would like to post to

The following services are currently the only supported services.
- Discord
- Slack

Below is an example of checking each RSS feed for articles posted in the last 24 hours and sending a message using a Discord webhook.

```bash
./rss-monitor discord -d -t templates/discord-message.txt -r templates/rss-links.txt -w  https://discord.com/api/webhooks/xxxxxxxx/yyyyyyyyyyyyyyyy
```

## Future Improvements

- Build in a cron scheduler to run on a routine basis without setting up crontabs in Linux
- Add additional services:
    - email

## References

- [gofeed](https://github.com/mmcdole/gofeed) - Parse RSS, Atom and JSON feeds in Go
