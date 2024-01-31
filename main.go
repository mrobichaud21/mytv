package main

import (
	"os"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/version"
)

var (
	namespace = "telly"
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

	// Web flags
	flag.Parse()

	if flag.Lookup("config.file").Changed {
		viper.SetConfigFile(flag.Lookup("config.file").Value.String())
	} else {
		viper.SetConfigName("telly.config")
		viper.AddConfigPath("/etc/telly/")
		viper.AddConfigPath("$HOME/.telly")
		viper.AddConfigPath(".")
		viper.SetEnvPrefix(namespace)
		viper.AutomaticEnv()
	}

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.WithError(err).Panicln("fatal error while reading config file:")
		}
	}

	level, parseLevelErr := logrus.ParseLevel(viper.GetString("log.level"))
	if parseLevelErr != nil {
		log.WithError(parseLevelErr).Panicln("error setting log level!")
	}
	log.SetLevel(level)

	log.Infoln("telly is preparing to go live", version.Info())

	validateConfig()
}

func validateConfig() {
	// if viper.IsSet("filter.regexstr") {
	// 	if _, regexErr := regexp.Compile(viper.GetString("filter.regex")); regexErr != nil {
	// 		log.WithError(regexErr).Panicln("Error when compiling regex, is it valid?")
	// 	}
	// }

	// if !(viper.IsSet("source")) {
	// 	log.Warnln("There is no source element in the configuration, the config file is likely missing.")
	// }

	// var addrErr error
	// if _, addrErr = net.ResolveTCPAddr("tcp", viper.GetString("web.listenaddress")); addrErr != nil {
	// 	log.WithError(addrErr).Panic("Error when parsing Listen address, please check the address and try again.")
	// 	return
	// }

	// if _, addrErr = net.ResolveTCPAddr("tcp", viper.GetString("web.base-address")); addrErr != nil {
	// 	log.WithError(addrErr).Panic("Error when parsing Base addresses, please check the address and try again.")
	// 	return
	// }

	// if getTCPAddr("web.base-address").IP.IsUnspecified() {
	// 	log.Panicln("base URL is set to 0.0.0.0, this will not work. please use the --web.baseaddress option and set it to the (local) ip address telly is running on.")
	// }

	// if getTCPAddr("web.listenaddress").IP.IsUnspecified() && getTCPAddr("web.base-address").IP.IsLoopback() {
	// 	log.Warnln("You are listening on all interfaces but your base URL is localhost (meaning Plex will try and load localhost to access your streams) - is this intended?")
	// }
}
