package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

const LevelFatal slog.Level = 12

type OutputFormat string

const (
	OutputFormatText = OutputFormat("text")
	OutputFormatJSON = OutputFormat("json")
)

var globalLogLevel = new(slog.LevelVar)

func SetLevel(level slog.Level) {
	globalLogLevel.Set(level)
}

func GetLevel() slog.Level {
	return globalLogLevel.Level()
}

var timestampFieldName = slog.TimeKey
var levelFieldName = slog.LevelKey
var messageFieldName = slog.MessageKey
var timestampFormat = time.RFC3339

func SetTimestampFieldName(name string) {
	timestampFieldName = name
}

func SetLevelFieldName(name string) {
	levelFieldName = name
}

func SetMessageFieldName(name string) {
	messageFieldName = name
}

func SetTimeFieldFormat(format string) {
	timestampFormat = format
}

func Initialize(format OutputFormat, domain string, service string) error {
	return InitializeWithOutput(format, domain, service, os.Stdout)
}

func InitializeWithOutput(format OutputFormat, domain string, service string, output io.Writer) error {
	logger, err := NewWithOptions(
		format,
		domain,
		service,
		Options{
			HandlerOptions: &slog.HandlerOptions{
				Level:       globalLogLevel,
				ReplaceAttr: DefaultReplaceAttrs(),
			},
			Output: output,
		},
	)
	if err != nil {
		return err
	}

	slog.SetDefault(logger)
	return nil
}

func InitializeWithHandler(handler slog.Handler, domain string, service string) {
	logger := NewWithHandler(handler, domain, service)
	slog.SetDefault(logger)
}

func New(format OutputFormat, domain string, service string) (*slog.Logger, error) {
	return NewWithOptions(format, domain, service, Options{
		HandlerOptions: &slog.HandlerOptions{
			Level: slog.LevelDebug,
			ReplaceAttr: ReplaceAttrs(
				"timestamp",
				"message",
				"level",
				time.RFC3339,
			),
		},
		Output: os.Stdout,
	})
}

type Options struct {
	*slog.HandlerOptions
	Output io.Writer
}

func NewWithOptions(
	format OutputFormat,
	domain string,
	service string,
	options Options,
) (*slog.Logger, error) {
	var handler slog.Handler
	switch format {
	case OutputFormatJSON:
		handler = slog.NewJSONHandler(options.Output, options.HandlerOptions)
	case OutputFormatText:
		handler = slog.NewTextHandler(options.Output, options.HandlerOptions)
	default:
		return nil, fmt.Errorf("invalid 'logging.output.format' configuration: %s", format)
	}

	return NewWithHandler(&ContextHandler{handler}, domain, service), nil
}

func NewWithHandler(handler slog.Handler, domain string, service string) *slog.Logger {
	return slog.New(handler).With("domain", domain, "service", service)
}

func Logger() *slog.Logger {
	return slog.Default()
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Fatal(msg string, args ...any) {
	FatalCtx(context.Background(), msg, args...)
}

func InfoCtx(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

func WarnCtx(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

func ErrorCtx(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

func DebugCtx(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

func FatalCtx(ctx context.Context, msg string, args ...any) {
	slog.Log(ctx, LevelFatal, msg, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	slog.Log(ctx, level, msg, args...)
}
