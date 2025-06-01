package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	env := os.Getenv("APP_ENV")

	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = NewHandler(os.Stdout, slog.LevelInfo)
	}

	return slog.New(handler)
}
