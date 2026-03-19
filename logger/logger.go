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
		// ReplaceAttr allows us to customize the order and display of default fields
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				// Use .Time() to get the time.Time object
				t := a.Value.Time()

				// Return a new attribute with your custom string format
				return slog.String("time", t.Format("2006-01-02 15:04:05"))
			}
			return a
		},
	}

	w := io.MultiWriter(os.Stdout, file)

	// Create the base logger
	baseHandler := slog.NewTextHandler(w, handlerOpts)

	// Use .With() so the service name always appears immediately after the level/time
	logger := slog.New(baseHandler).With("src", serviceName)
	slog.SetDefault(logger)
	return logger
}
