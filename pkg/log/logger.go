package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type option struct {
	Encoder  Encoder
	MinLevel Level
	Service  string
	Env      string
	Version  string
	Commit   string
}

type Option func(*option)

func WithEncoder(e Encoder) Option {
	return func(o *option) {
		o.Encoder = e
	}
}

func WithLocalEncoder() Option {
	return func(o *option) {
		o.Encoder = EncoderLocal
	}
}

func WithDatadogEncoder() Option {
	return func(o *option) {
		o.Encoder = EncoderDatadog
	}
}

func WithMinLevel(l Level) Option {
	return func(o *option) {
		o.MinLevel = l
	}
}

func WithService(service string) Option {
	return func(o *option) {
		o.Service = service
	}
}

func WithVersion(version string) Option {
	return func(o *option) {
		o.Version = version
	}
}

func WithEnv(e string) Option {
	return func(o *option) {
		o.Env = e
	}
}

func WithCommit(commit string) Option {
	return func(o *option) {
		o.Commit = commit
	}
}

func NewLogger(opts ...Option) (Logger, error) {
	opt := option{
		MinLevel: LevelInfo,
		Encoder:  EncoderLocal,
	}

	for _, setter := range opts {
		setter(&opt)
	}

	minLevel, err := parseLevel(opt.MinLevel.String())
	if err != nil {
		return nil, err
	}

	var sl *slog.Logger
	switch opt.Encoder {
	case EncoderLocal:
		sl = newLocalSLogger(minLevel)
	case EncoderDatadog:
		sl = newDatadogSLogger(minLevel, opt.Service, opt.Version, opt.Env, opt.Commit)
	default:
		return nil, fmt.Errorf("invalid log encoder: %q", opt.Encoder)
	}

	return &sLogger{
		option: opt,
		sl:     sl,
	}, nil
}

func NewNop() Logger {
	return &sLogger{
		option: option{
			MinLevel: LevelInfo,
			Encoder:  EncoderLocal,
		},
		sl: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func parseLevel(l string) (slog.Level, error) {
	var sLevel slog.Level
	switch l {
	case LevelDebug.String():
		sLevel = slog.LevelDebug
	case LevelInfo.String():
		sLevel = slog.LevelInfo
	case LevelWarn.String():
		sLevel = slog.LevelWarn
	case LevelError.String():
		sLevel = slog.LevelError
	case LevelFatal.String():
		sLevel = SlogLevelFatal
	default:
		return sLevel, fmt.Errorf("unknown level %q", l)
	}

	return sLevel, nil
}

type sLogger struct {
	option option
	sl     *slog.Logger
}

func (l *sLogger) Debug(ctx context.Context, msg string, fields ...Attr) {
	l.sl.LogAttrs(ctx, slog.LevelDebug, msg, fields...)
}

func (l *sLogger) Info(ctx context.Context, msg string, fields ...Attr) {
	l.sl.LogAttrs(ctx, slog.LevelInfo, msg, fields...)
}

func (l *sLogger) Warn(ctx context.Context, msg string, fields ...Attr) {
	l.sl.LogAttrs(ctx, slog.LevelWarn, msg, fields...)
}

func (l *sLogger) Error(ctx context.Context, msg string, fields ...Attr) {
	l.sl.LogAttrs(ctx, slog.LevelError, msg, fields...)
}

func (l *sLogger) Fatal(ctx context.Context, msg string, fields ...Attr) {
	l.sl.LogAttrs(ctx, SlogLevelFatal, msg, fields...)
}

func (l *sLogger) Level() Level {
	return l.option.MinLevel
}

func (l *sLogger) With(fields ...Attr) Logger {
	converted := make([]any, 0, len(fields))
	for _, f := range fields {
		converted = append(converted, f)
	}

	return &sLogger{
		option: l.option,
		sl:     l.sl.With(converted...),
	}
}
