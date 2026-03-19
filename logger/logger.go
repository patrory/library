package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

func InitLogger(serviceName string) *slog.Logger {
	// Note: Changed the date format string. "2006-01-02" is the standard Go layout.
	logFile := fmt.Sprintf("logs/%s-%s.log", serviceName, time.Now().Format("2006-01-02"))

	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(fmt.Errorf("failed to create logs directory: %w", err))
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Errorf("failed to open log file: %w", err))
	}

	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 1. Format Time
			if a.Key == slog.TimeKey {
				return slog.String("time", a.Value.Time().Format("2006-01-02 15:04:05"))
			}

			// 2. Remove 'level' key and replace it with 'srv' + 'level'
			// but as a single string to avoid slog.Group recursion
			if a.Key == slog.LevelKey {
				// This makes the output look like: ... srv=AuthSrv level=INFO ...
				// by injecting the srv field into the metadata flow safely
				return slog.Attr{
					Key:   "srv",
					Value: slog.StringValue(fmt.Sprintf("%s level=%s", serviceName, a.Value.String())),
				}
			}
			return a
		},
	}

	w := io.MultiWriter(os.Stdout, file)

	logger := slog.New(slog.NewTextHandler(w, handlerOpts))
	slog.SetDefault(logger)
	return logger
}
