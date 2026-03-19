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
			// 1. Format the Timestamp
			if a.Key == slog.TimeKey {
				return slog.String("time", a.Value.Time().Format("2006-01-02 15:04:05"))
			}

			// 2. Inject Service Name right before/with the Level
			if a.Key == slog.LevelKey {
				// We return a "Group" with no name.
				// TextHandler will print these attributes in order.
				return slog.Group("",
					slog.String("srv", serviceName),
					slog.Attr(a), // This is the original Level attribute
				)
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
