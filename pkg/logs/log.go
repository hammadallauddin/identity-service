package logs

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/hammadallauddin/identity-service/pkg/config"
)

const LevelFatal slog.Level = 12

type OutputFormat string

const (
	OutputFormatText = OutputFormat("text")
	OutputFormatJSON = OutputFormat("json")
)

var globalLogLevel = new(slog.LevelVar)

func setLevel(level slog.Level) {
	globalLogLevel.Set(level)
}

func GetLevel() slog.Level {
	return globalLogLevel.Level()
}

var timestampFieldName = slog.TimeKey
var levelFieldName = slog.LevelKey
var messageFieldName = slog.MessageKey
var timestampFormat = time.RFC3339

func setTimestampFieldName(name string) {
	timestampFieldName = name
}

func setLevelFieldName(name string) {
	levelFieldName = name
}

func setMessageFieldName(name string) {
	messageFieldName = name
}

func setTimeFieldFormat(format string) {
	timestampFormat = format
}

func Initialize() (*slog.Logger, error) {
	level, err := config.GetString("logging.level")
	if err != nil {
		return nil, fmt.Errorf("initializeLogging(): invalid 'logging.level' configuration: %w", err)
	}
	var logLevel slog.Level
	switch level {
	case "info":
		logLevel = slog.LevelInfo
	case "error":
		logLevel = slog.LevelError
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	default:
		return nil, fmt.Errorf("initializeLogging(): invalid 'logging.level' configuration: %s", level)
	}
	setLevel(logLevel)

	timestampKey, _ := config.GetString("logging.output.timestamp-key", "timestamp")
	setTimestampFieldName(timestampKey)

	levelKey, _ := config.GetString("logging.output.level-key", "severity")
	setLevelFieldName(levelKey)

	messageKey, _ := config.GetString("logging.output.message-key", "message")
	setMessageFieldName(messageKey)

	timeFieldFormat, _ := config.GetString("logging.output.time-field-format", time.RFC3339)
	setTimeFieldFormat(timeFieldFormat)

	serviceName, err := config.GetString("service.name")
	if err != nil {
		return nil, fmt.Errorf("initializeLogging(): invalid 'service.name' configuration: %w", err)
	}

	domainName, err := config.GetString("logging.domain", "default")
	if err != nil {
		return nil, fmt.Errorf("initializeLogging(): invalid 'logging.domain' configuration: %w", err)
	}

	var outputFormat OutputFormat
	optFmt, err := config.GetString("logging.output.format", "json")
	if err != nil {
		return nil, fmt.Errorf("initializeLogging(): invalid 'logging.output.format' configuration: %w", err)
	}
	switch optFmt {
	case "text":
		outputFormat = OutputFormatText
	default:
		outputFormat = OutputFormatJSON
	}

	return initializeWithOutput(outputFormat, domainName, serviceName, os.Stdout)
}

func initializeWithOutput(format OutputFormat, domain string, service string, output io.Writer) (*slog.Logger, error) {
	logger, err := newWithOptions(
		format,
		domain,
		service,
		options{
			HandlerOptions: &slog.HandlerOptions{
				Level:       globalLogLevel,
				ReplaceAttr: DefaultReplaceAttrs(),
			},
			Output: output,
		},
	)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func initializeWithHandler(handler slog.Handler, domain string, service string) {
	logger := newWithHandler(handler, domain, service)
	slog.SetDefault(logger)
}

func New(format OutputFormat, domain string, service string) (*slog.Logger, error) {
	return newWithOptions(format, domain, service, options{
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

type options struct {
	*slog.HandlerOptions
	Output io.Writer
}

func newWithOptions(
	format OutputFormat,
	domain string,
	service string,
	options options,
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

	return newWithHandler(&ContextHandler{handler}, domain, service), nil
}

func newWithHandler(handler slog.Handler, domain string, service string) *slog.Logger {
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
