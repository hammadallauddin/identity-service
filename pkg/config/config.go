package config

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/hammadallauddin/identity-service/pkg/log"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const (
	invalidConfiguration = "Invalid configuration '%s'"
	missingConfiguration = "Missing configuration '%s'"
)

func initializeLogging() error {
	serviceName, err := GetString("service.name")
	if err != nil {
		return fmt.Errorf("initializeLogging(): invalid 'service.name' configuration: %w", err)
	}

	level, err := GetString("logging.level")
	if err != nil {
		return fmt.Errorf("initializeLogging(): invalid 'logging.level' configuration: %w", err)
	}

	if err := log.SetLevel(level); err != nil {
		return fmt.Errorf("initializeLogging(): invalid 'logging.level' configuration: level=%s error=%w", level, err)
	}

	domainName, _ := GetString("logging.domain", "default")
	outputFormat, _ := GetString("logging.output.format", "simple")
	timestampKey, _ := GetString("logging.output.timestamp-key", "@timestamp")
	levelKey, _ := GetString("logging.output.level-key", "severity")
	messageKey, _ := GetString("logging.output.message-key", "message")
	timeFieldFormat, _ := GetString("logging.output.time-field-format", time.RFC3339)

	log.SetTimestampFieldName(timestampKey)
	log.SetLevelFieldName(levelKey)
	log.SetMessageFieldName(messageKey)
	log.SetTimeFieldFormat(timeFieldFormat)

	return log.Initialize(outputFormat, domainName, serviceName)
}

func Reset() {
	args := strings.Split(os.Getenv("FLAG_FOR_MAIN"), ",")
	os.Args = append([]string{os.Args[0]}, args...)

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.CommandLine = flags

	viper.Reset()
}

func Initialize() error {
	configFilePtr := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	configFilePath := os.Getenv("CONFIG_PATH")
	if len(configFilePath) == 0 {
		configFilePath = *configFilePtr
		if len(configFilePath) == 0 {
			return fmt.Errorf("Initialize(): missing configuration file")
		}
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return fmt.Errorf("Initialize(): configuration file does not exist: %s", configFilePath)
	}

	configPath, configFile := filepath.Split(configFilePath)
	ext := path.Ext(configFile)
	if ext != ".yml" && ext != ".yaml" {
		return fmt.Errorf("Initialize(): invalid configuration file extension: %s", configFilePath)
	}
	configFile = strings.TrimSuffix(configFile, ext)

	viper.SetConfigType("yaml")
	viper.SetConfigName(configFile)
	viper.AddConfigPath(configPath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Initialize(): %w", err)
	}

	err = initializeLogging()
	if err != nil {
		return fmt.Errorf("Initialize(): %w", err)
	}

	return nil
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
