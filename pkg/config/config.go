package config

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/spf13/viper"
	"gitlab.com/startupbuilder/startupbuilder/pkg/config/replacer"
)

const (
	defaultKeyDelimiter          = "_"
	defaultSecondaryKeyDelimiter = "0"
	defaultKeyTag                = "mapstructure"
	defaultConfigType            = "env"
)

type ConfigMap interface {
	Defaults()
}

// ReadConfig reads the environment variables from file or system and set the default ones before starting.
func ReadConfig(cfg ConfigMap, opts ...Option) error {
	stngs := newSetting()

	cfg.Defaults()
	stngs.apply(opts)

	replacer := replacer.NewReplacer(stngs.prefix, stngs.keyDelimiter, stngs.secondaryKeyDelimiter)

	configReader := viper.NewWithOptions(
		viper.KeyDelimiter(stngs.keyDelimiter),
		viper.EnvKeyReplacer(replacer),
	)

	if err := stngs.bindEnvs(configReader, cfg, ""); err != nil {
		return fmt.Errorf("could not bind envs: %w", err)
	}

	configReader.AutomaticEnv()

	// read the envs from the system.
	if stngs.filePath == "" {
		if err := configReader.ReadConfig(&io.LimitedReader{}); err != nil {
			return fmt.Errorf("could not read configurations: %w", err)
		}
	} else {
		configReader.SetConfigType(defaultConfigType)
		configReader.SetConfigFile(stngs.filePath)

		if err := configReader.ReadInConfig(); err != nil {
			return fmt.Errorf("could not read configurations: %w", err)
		}
	}

	if err := configReader.Unmarshal(cfg); err != nil {
		return fmt.Errorf("could not unmarshal configurations: %w", err)
	}

	return nil
}

type setting struct {
	filePath              string
	prefix                string
	keyDelimiter          string
	secondaryKeyDelimiter string
	keyTag                string
	configType            string
}

func newSetting() *setting {
	return &setting{
		keyDelimiter:          defaultKeyDelimiter,
		secondaryKeyDelimiter: defaultSecondaryKeyDelimiter,
		keyTag:                defaultKeyTag,
		configType:            defaultConfigType,
	}
}

func (s *setting) apply(opts []Option) error {
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return fmt.Errorf("could not apply setting: %w", err)
		}
	}

	return nil
}

func (stngs *setting) bindEnvs( //nolint:funlen
	configReader *viper.Viper,
	cfg interface{},
	prefix string,
) error {
	values := reflect.ValueOf(cfg)
	if values.Kind() == reflect.Interface && !values.IsNil() {
		elm := values.Elem()
		if elm.Kind() == reflect.Pointer && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			values = elm
		}
	}

	if values.Kind() == reflect.Pointer {
		values = values.Elem()
	}

	for i := range values.NumField() {
		fieldValue := values.Field(i)
		fieldValue.CanInterface()
		tag := values.Type().Field(i).Tag.Get(stngs.keyTag)

		if tag == ",squash" {
			tag = ""
		}

		tag = strings.ToUpper(tag)

		if strings.Contains(tag, defaultKeyDelimiter) {
			err := fmt.Errorf(
				"can not have same keyDelimiter (%s) as environment name in %s, instead use (%s)",
				stngs.keyDelimiter,
				prefix+defaultKeyDelimiter+tag,
				defaultSecondaryKeyDelimiter,
			)

			return err
		}

		if prefix != "" {
			if tag != "" {
				tag = prefix + defaultKeyDelimiter + tag
			} else {
				tag = prefix
			}
		}

		if fieldValue.Kind() == reflect.Pointer && !fieldValue.IsNil() {
			elm := fieldValue.Elem()
			if elm.Kind() == reflect.Pointer && !elm.IsNil() &&
				elm.Elem().Kind() == reflect.Pointer {
				fieldValue = elm
			}

			if fieldValue.Kind() == reflect.Ptr {
				fieldValue = fieldValue.Elem()
			}

			if fieldValue.CanInterface() {
				stngs.bindEnvs(configReader, fieldValue.Interface(), tag)
			}
			continue
		}

		if fieldValue.Kind() == reflect.Struct {
			stngs.bindEnvs(configReader, fieldValue.Interface(), tag)

			continue
		}

		configReader.BindEnv(tag)
	}

	return nil
}

type Option func(*setting) error

// WithFilePath allow to read envs saved in the file.
func WithFilePath(filePath string) Option {
	return func(s *setting) error {
		if filePath == "" {
			return nil
		}

		s.filePath = filePath

		return nil
	}
}

// WithPrefix adding a prefix to all keys, when searching in environment variables.
func WithPrefix(prefix string) Option {
	return func(s *setting) error {
		if prefix == "" {
			return nil
		}

		s.prefix = prefix

		return nil
	}
}

// WithKeyDelimiter is the default subStructs seprator, default: `_`.
func WithKeyDelimiter(delimiter string) Option {
	return func(s *setting) error {
		if delimiter == "" {
			return nil
		}

		s.keyDelimiter = delimiter

		return nil
	}
}

// WithSecondaryKeyDelimiter change the secondary delimiter to replace with primary delimiter, default: `0`.
func WithSecondaryKeyDelimiter(delimiter string) Option {
	return func(s *setting) error {
		if delimiter == "" {
			return nil
		}

		s.secondaryKeyDelimiter = delimiter

		return nil
	}
}

// WithKeyTag set different key than `mapstructure` as struct tag to find.
func WithKeyTag(tag string) Option {
	return func(s *setting) error {
		if tag == "" {
			return nil
		}

		s.keyTag = tag

		return nil
	}
}

// WithConfigType change the config file type, default to `env`.
func WithConfigType(configType string) Option {
	return func(s *setting) error {
		if configType == "" {
			return nil
		}

		s.configType = configType

		return nil
	}
}
