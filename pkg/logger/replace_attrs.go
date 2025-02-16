package logger

import (
	"log/slog"
	"strings"
)

func DefaultReplaceAttrs() func(groups []string, a slog.Attr) slog.Attr {
	return ReplaceAttrs(
		timestampFieldName,
		messageFieldName,
		levelFieldName,
		timestampFormat,
	)
}

func ReplaceAttrs(timestampFieldName, messageFieldName, levelFieldName, timestampFormat string) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.TimeKey:
			a.Key = timestampFieldName
		case slog.MessageKey:
			a.Key = messageFieldName
		case slog.LevelKey:
			a.Key = levelFieldName
			level, ok := a.Value.Any().(slog.Level)
			if ok {
				switch level {
				case LevelFatal:
					a.Value = slog.StringValue("FATAL")
				default:
					a.Value = slog.StringValue(strings.ToLower(level.String()))
				}
			}
		}

		switch a.Value.Kind() {
		case slog.KindTime:
			slog.StringValue(a.Value.Time().Format(timestampFormat))
		case slog.KindDuration:
			a.Value = slog.StringValue(a.Value.Duration().String())
		}

		return a
	}

}
