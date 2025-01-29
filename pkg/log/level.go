package log

import "log/slog"

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelFatal Level = "fatal"
)

func (l Level) String() string {
	return string(l)
}

// slog default levels: DEBUG (-4), INFO (0), WARN (4), and ERROR (8).
// Below defines custom levels
const SlogLevelFatal = slog.Level(12)

var LevelNames = map[slog.Leveler]string{
	SlogLevelFatal: "FATAL",
}
