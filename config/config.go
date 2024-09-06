package config

import (
	"errors"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

const ProfileEnvKey string = "ACTIVE_PROFILE"

var configs map[string]any

func init() {

	configFileName := "config"
	env, profileSet := os.LookupEnv(ProfileEnvKey)
	if profileSet {
		log.Info().Str("Setting active profile to ", env).Send()
		configFileName = configFileName + "-" + env
	}
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		//set default to trace
		lvl = "TRACE"

	}

	switch lvl {
	case "PANIC":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "FATAL":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	log.Info().Msg("Completed setting log level to " + lvl)

	if len(env) == 0 {
		log.Info().Msg("Setting Active Profile to NONE")
	}

	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Error().AnErr("Error reading config file", err).Send()
		panic(-1)
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			val := getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}"))
			viper.Set(k, val)
		}
	}

	err := viper.Unmarshal(&configs)
	if err != nil {
		log.Error().AnErr("Error creating configs", err).Send()
		panic(-1)
	}

	log.Info().Msg("Configurations loaded successfully")

}

func getEnvOrPanic(env string) string {
	res, set := os.LookupEnv(env)
	if !set {
		log.Panic().Msg("Config load failed")
		panic("Mandatory env variable not found:" + env)
	}
	return res
}

func GetConfig[T any](key string) (*T, error) {

	raw := configs[key]
	if typed, ok := raw.(T); ok {
		return &typed, nil
	}

	return nil, errors.New("cannot convert config to tpye")

}
