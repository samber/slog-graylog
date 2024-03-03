package sloggraylog

import "log/slog"

var LogLevels = map[slog.Level]int32{
	slog.LevelDebug: 7,
	slog.LevelInfo:  6,
	slog.LevelWarn:  4,
	slog.LevelError: 3,
}
