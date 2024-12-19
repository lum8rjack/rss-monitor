package utils

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

func GetRssUpdates(rssFile string, timeframe int) ([]Post, error) {
	var posts []Post

	// Check timeframe
	if timeframe < 0 {
		return posts, errors.New("timeframe must be greater than 0")
	}
	rssLogger.Debug("valid timeframe to check for", "timeframe", timeframe)

	// Read RSS links
	r, err := readRssLinks(rssFile)
	if err != nil {
		return posts, err
	}
	rssLogger.Debug("successfully read in RSS links", "rss_file", rssFile)

	// Check each RSS fee
	for _, l := range r {
		p, err := checkRSSFeed(l, timeframe)
		if err != nil {
			rssLogger.Error(err.Error(), "url", l)
		}
		if len(p) > 0 {
			posts = append(posts, p...)
		}
	}
	return posts, nil
}

func readRssLinks(rssFile string) ([]string, error) {
	if rssFile == "" {
		return nil, errors.New("must provide a file containing the RSS links")
	}

	file, err := os.Open(rssFile)
	if err != nil {
		e := fmt.Sprintf("error opening file: %v", err)
		return nil, errors.New(e)
	}
	defer file.Close()

	var links []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "https://") || strings.HasPrefix(line, "http://") {
				links = append(links, line)
			}
		}
	}

	if len(links) == 0 {
		return nil, errors.New("no links pulled from the file")
	}

	return links, nil
}

func checkRSSFeed(url string, timeframe int) ([]Post, error) {
	pastHours := 0 - timeframe
	pastTime := time.Now().Add(time.Duration(pastHours) * time.Hour) // 24 hours ago

	webclient := &http.Client{
		Timeout: time.Second * 15,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

	response, err := webclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	rssLogger.Debug("web request completed", "link", url)

	fp := gofeed.NewParser()
	feed, err := fp.Parse(response.Body)
	if err != nil {
		return nil, err
	}

	var posts []Post

	for _, i := range feed.Items {
		date, err := parseTimeWithFallback(i.Published)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if pastTime.Before(date) {
			p := Post{
				Title:     i.Title,
				Link:      i.Link,
				Published: i.Published,
			}
			posts = append(posts, p)
		}
	}

	return posts, nil
}
