package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	namespace = "mytv"
	log       = &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			FullTimestamp: true,
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.DebugLevel,
	}
)

func main() {

	viper.SetConfigName("mytv.config") // name of config file (without extension)
	viper.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/mytv/")
	viper.AddConfigPath("$HOME/.mytv")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix(namespace)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		return
	}

	level, parseLevelErr := logrus.ParseLevel(config.Log.Level)
	if parseLevelErr != nil {
		log.WithError(parseLevelErr).Panicln("error setting log level!")
	}

	log.SetLevel(level)

	lineup := newLineup(config)
	lineup.Scan()
	serve(lineup)

}
