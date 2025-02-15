package log

import (
	"errors"
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	logOutputFormatSimple = "simple"
	logOutputFormatJSON   = "json"
)

type Context = map[string]interface{}

func SetLevel(level string) error {
	switch level {
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		return errors.New("Invalid logging.level configuration")
	}
	return nil
}

func GetLevel() string {
	switch zerolog.GlobalLevel() {
	case zerolog.ErrorLevel:
		return "error"
	case zerolog.WarnLevel:
		return "warn"
	case zerolog.InfoLevel:
		return "info"
	case zerolog.DebugLevel:
		return "debug"
	default:
		return "error"
	}
}

func SetTimestampFieldName(name string) {
	zerolog.TimestampFieldName = name
}

func SetLevelFieldName(name string) {
	zerolog.LevelFieldName = name
}

func SetMessageFieldName(name string) {
	zerolog.MessageFieldName = name
}

func SetTimeFieldFormat(format string) {
	zerolog.TimeFieldFormat = format
}

func Initialize(format string, domain string, service string) error {
	if format != logOutputFormatSimple && format != logOutputFormatJSON {
		return errors.New("Invalid logging.output.format configuration")
	}
	if format == logOutputFormatJSON {
		log.Logger = zerolog.New(os.Stdout).With().Str("domain", domain).Str("service", service).Timestamp().Logger()
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Str("domain", domain).Str("service", service).Timestamp().Logger()
	}
	return nil
}

func Logger() zerolog.Logger {
	return log.Logger
}

func Info(format string, v ...interface{}) {
	withCaller(log.Info()).Msgf(format, v...)
}

func Warn(format string, v ...interface{}) {
	withCaller(log.Warn()).Msgf(format, v...)
}

func Error(format string, v ...interface{}) {
	withCaller(log.Error()).Msgf(format, v...)
}

func Debug(format string, v ...interface{}) {
	withCaller(log.Debug()).Msgf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	withCaller(log.Fatal()).Msgf(format, v...)
}

func InfoC(ctx Context, format string, v ...interface{}) {
	withCaller(log.Info()).Fields(ctx).Msgf(format, v...)
}

func WarnC(ctx Context, format string, v ...interface{}) {
	withCaller(log.Warn()).Fields(ctx).Msgf(format, v...)
}

func ErrorC(ctx Context, format string, v ...interface{}) {
	withCaller(log.Error()).Fields(ctx).Msgf(format, v...)
}

func DebugC(ctx Context, format string, v ...interface{}) {
	withCaller(log.Debug()).Fields(ctx).Msgf(format, v...)
}

func FatalC(ctx Context, format string, v ...interface{}) {
	withCaller(log.Fatal()).Fields(ctx).Msgf(format, v...)
}

func withCaller(event *zerolog.Event) *zerolog.Event {
	var (
		funcName string
	)
	pc, _, _, ok := runtime.Caller(2)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
		event = event.Str("caller", funcName)
	}
	return event
}
