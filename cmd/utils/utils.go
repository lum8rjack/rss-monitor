package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type BaseArgs struct {
	Crontab      string
	Debug        bool
	RssFile      string
	TemplateFile string
	Timewindow   int
	Logger       *slog.Logger
}

func GetArgsData(cmd *cobra.Command) (BaseArgs, error) {
	ba := BaseArgs{}

	// Setup logger
	d, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return ba, err
	}

	// Get arguments from user
	rss, err := cmd.Flags().GetString("rss")
	if err != nil {
		return ba, err
	}

	template, err := cmd.Flags().GetString("template")
	if err != nil {
		return ba, err
	}

	timeframe, err := cmd.Flags().GetInt("timeframe")
	if err != nil {
		return ba, err
	}

	// Set arguments
	ba.Logger = NewLogger(d)
	ba.RssFile = rss
	ba.TemplateFile = template
	ba.Timewindow = timeframe

	return ba, err
}

func CheckError(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func parseTimeWithFallback(value string) (time.Time, error) {
	var t time.Time
	var e error
	if value == "" {
		return t, errors.New("cannot pass empty string to parse as time")
	}

	t, e = time.Parse(time.ANSIC, value)
	if e != nil {
		t, e = time.Parse(time.RFC1123, value)
		if e != nil {
			t, e = time.Parse(time.RFC1123Z, value)
			if e != nil {
				t, e = time.Parse(time.RFC3339, value)
				if e != nil {
					t, e = time.Parse(time.UnixDate, value)
				}
			}
		}
	}
	return t, e
}
