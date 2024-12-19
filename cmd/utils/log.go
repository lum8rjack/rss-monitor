package utils

import (
	"log/slog"
	"os"
)

var (
	rssLogger *slog.Logger
)

func NewLogger(debug bool) *slog.Logger {
	ops := &slog.HandlerOptions{}
	if debug {
		ops.Level = slog.LevelDebug
	}
	rssLogger = slog.New(slog.NewJSONHandler(os.Stdout, ops))

	return rssLogger
}
