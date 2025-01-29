package log

import (
	"log/slog"
	"os"
)

func newLocalSLogger(level slog.Level) *slog.Logger {
	replaceAttrFunc := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			level := a.Value.Any().(slog.Level) //nolint:errcheck,forcetypeassert
			levelLabel, exists := LevelNames[level]
			if !exists {
				levelLabel = level.String()
			}

			a.Value = slog.StringValue(levelLabel)
		}

		return a
	}

	options := &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: replaceAttrFunc,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, options))

	return logger
}
