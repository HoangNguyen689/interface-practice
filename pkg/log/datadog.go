package log

import (
	"log/slog"
	"os"
)

func newDatadogSLogger(level slog.Level, service, version, env, commit string) *slog.Logger {
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

	handler := slog.NewJSONHandler(os.Stdout, options).
		WithAttrs([]slog.Attr{
			slog.String("service", service),
			slog.String("version", version),
			slog.String("env", env),
			slog.String("commit", commit),
			slog.String("logger.type", "Datadog"),
			slog.Group("dd",
				slog.String("service", service),
				slog.String("version", version),
				slog.String("env", env),
			),
		})

	logger := slog.New(handler)

	return logger
}
