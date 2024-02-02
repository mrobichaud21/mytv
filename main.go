package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	Playlist struct {
		URL interface{} `yaml:"url"`
	} `yaml:"playlist"`
	Filters []struct {
		GroupTitle string `yaml:"group-title"`
		Channels   []struct {
			Channel string `yaml:"channel"`
		} `yaml:"channels"`
	} `yaml:"filters"`
}

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

	playlist, err := getPlaylist(config.Playlist.URL.(string))

	if err != nil {
		panic(err)
	}

	//file, err := m3uplus.Decode(playlist)
	scanner := bufio.NewScanner(playlist)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	filterChannels := []string{}
	lines := ""

	for scanner.Scan() {

		lineText := scanner.Text()
		//fmt.Println(lineText)

		if strings.HasPrefix(lineText, "#EXTINF") && strings.Contains(lineText, "group-title=\"US| ENTERTAINMENT HD/4K\"") {
			lines = lineText
		} else {
			if len(lines) > 0 {
				filterChannels = append(filterChannels, lines+" url=\""+lineText+"\"")
				lines = ""
			}
		}

	}

	fmt.Printf("Total channels = %d", len(filterChannels))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for i, s := range filterChannels {
		fmt.Println(i, s)
	}

	// rawPlaylist, err := m3uplus.Decode(playlist)
	// if err != nil {
	// 	log.WithError(err).Errorln("unable to parse m3u file")
	// 	panic(err)
	// }

	// fmt.Print(len(rawPlaylist.Tracks))
	// playlist.Close()

}

func extractChannelName(line string, key string) string {
	var ch string

	if strings.Contains(line, key+"=\"") {
		ch = strings.SplitN(strings.SplitN(line, key+"=\"", 2)[1], "\" ", 2)[0]
		// } else if strings.Contains(line, "tvg-name=") {
		// 	ch = strings.SplitN(strings.SplitN(line, "tvg-name=", 2)[1], " tvg", 2)[0]
	} else {
		ch = strings.TrimSpace(strings.SplitN(line, ",", 2)[1])
	}
	if ch == "" {
		ch = "No Name"
	}
	ch = strings.ReplaceAll(ch, "\"", "")
	return ch
}
