package log

import (
	"context"
	"log/slog"
	"time"
)

type Attr = slog.Attr

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Attr)
	Info(ctx context.Context, msg string, fields ...Attr)
	Warn(ctx context.Context, msg string, fields ...Attr)
	Error(ctx context.Context, msg string, fields ...Attr)
	Fatal(ctx context.Context, msg string, fields ...Attr)

	With(fields ...Attr) Logger

	Level() Level
}

type Encoder string

const (
	EncoderLocal   Encoder = "local"
	EncoderDatadog Encoder = "datadog"
)

func (e Encoder) String() string {
	return string(e)
}

func String(key, v string) Attr {
	return slog.String(key, v)
}

func Int64(key string, v int64) Attr {
	return slog.Int64(key, v)
}

func Int(key string, v int) Attr {
	return slog.Int(key, v)
}

func Uint64(key string, v uint64) Attr {
	return slog.Uint64(key, v)
}

func Float64(key string, v float64) Attr {
	return slog.Float64(key, v)
}

func Bool(key string, v bool) Attr {
	return slog.Bool(key, v)
}

func Time(key string, v time.Time) Attr {
	return slog.Time(key, v)
}

func Duration(key string, v time.Duration) Attr {
	return slog.Duration(key, v)
}

// TODO: need to handle error when passing to Datadog
func Error(v error) Attr {
	return slog.Any("error", v)
}

func Any(key string, v any) Attr {
	return slog.Any(key, v)
}
