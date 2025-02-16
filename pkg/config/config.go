package config

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hammadallauddin/identity-service/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	invalidConfiguration = "invalid configuration '%s'"
	missingConfiguration = "missing configuration '%s'"
)

func initializeLogging() error {
	level, err := GetString("logging.level")
	if err != nil {
		return fmt.Errorf("initializeLogging(): invalid 'logging.level' configuration: %w", err)
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
		return fmt.Errorf("initializeLogging(): invalid 'logging.level' configuration: %s", level)
	}
	logger.SetLevel(logLevel)

	timestampKey, _ := GetString("logging.output.timestamp-key", "@timestamp")
	logger.SetTimestampFieldName(timestampKey)

	levelKey, _ := GetString("logging.output.level-key", "severity")
	logger.SetLevelFieldName(levelKey)

	messageKey, _ := GetString("logging.output.message-key", "message")
	logger.SetMessageFieldName(messageKey)

	timeFieldFormat, _ := GetString("logging.output.time-field-format", time.RFC3339)
	logger.SetTimeFieldFormat(timeFieldFormat)

	serviceName, err := GetString("service.name")
	if err != nil {
		return fmt.Errorf("initializeLogging(): invalid 'service.name' configuration: %w", err)
	}

	domainName, err := GetString("logging.domain", "default")
	if err != nil {
		return fmt.Errorf("initializeLogging(): invalid 'logging.domain' configuration: %w", err)
	}

	var outputFormat logger.OutputFormat
	optFmt, err := GetString("logging.output.format", "json")
	if err != nil {
		return fmt.Errorf("initializeLogging(): invalid 'logging.output.format' configuration: %w", err)
	}
	switch optFmt {
	case "text":
		outputFormat = logger.OutputFormatText
	default:
		outputFormat = logger.OutputFormatJSON
	}

	err = logger.Initialize(outputFormat, domainName, serviceName)
	if err != nil {
		return fmt.Errorf("unable to initialize logger: %w", err)
	}

	return nil
}

func Reset() {
	args := strings.Split(os.Getenv("FLAG_FOR_MAIN"), ",")
	os.Args = append([]string{os.Args[0]}, args...)

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.CommandLine = flags

	viper.Reset()
}

func Initialize() error {
	configDirPtr := flag.String("config", "", "Path to configuration directory")
	flag.Parse()

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	configDir := firstNonEmpty(*configDirPtr, os.Getenv("CONFIG_PATH"))
	if configDir == "" {
		return errors.New("configuration directory not specified")
	}
	env := firstNonEmpty(os.Getenv("ENVIRONMENT"), "development")
	configFile := filepath.Join(configDir, fmt.Sprintf("%s-config.yaml", env))
	configFile = filepath.ToSlash(configFile)

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("config read error: %w", err)
	}

	if err := initializeLogging(); err != nil {
		return fmt.Errorf("logging initialization failed: %w", err)
	}

	return nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func Get(key string, defaults ...interface{}) interface{} {
	value := viper.Get(key)
	if value == nil && len(defaults) > 0 {
		return defaults[0]
	}
	return value
}

func GetInt64(key string, defaults ...int64) (int64, error) {
	value := viper.Get(key)
	if value == 0 && len(defaults) > 0 {
		return defaults[0], nil
	} else if value == nil {
		return 0, fmt.Errorf(missingConfiguration, key)
	}
	v, err := cast.ToInt64E(value)
	if err != nil {
		return 0, fmt.Errorf(invalidConfiguration, key)
	}
	return v, nil
}

func GetString(key string, defaults ...string) (string, error) {
	var value interface{}
	value = viper.Get(key)
	if value == nil {
		if len(defaults) == 0 {
			return "", fmt.Errorf(missingConfiguration, key)
		}
		value = defaults[0]
	}

	sv, ok := value.(string)
	if !ok {
		return "", fmt.Errorf(invalidConfiguration, key)
	}
	return sv, nil
}

func GetBool(key string, defaults ...bool) (bool, error) {
	var value interface{}
	value = viper.Get(key)
	if value == nil {
		if len(defaults) == 0 {
			return false, fmt.Errorf(missingConfiguration, key)
		}
		value = defaults[0]
	}

	bv, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf(invalidConfiguration, key)
	}
	return bv, nil
}

func GetInt(key string, defaults ...int) (int, error) {
	var value interface{}
	value = viper.Get(key)
	if value == nil {
		if len(defaults) == 0 {
			return 0, fmt.Errorf(missingConfiguration, key)
		}
		value = defaults[0]
	}

	iv, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf(invalidConfiguration, key)
	}
	return iv, nil
}

func GetFloat(key string, defaults ...float64) (float64, error) {
	var value interface{}
	value = viper.Get(key)
	if value == nil {
		if len(defaults) == 0 {
			return 0, fmt.Errorf(missingConfiguration, key)
		}
		value = defaults[0]
	}

	fv, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf(invalidConfiguration, key)
	}
	return fv, nil
}

func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
}
