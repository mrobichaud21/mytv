package main

// https://zhwt.github.io/yaml-to-go/
type Configuration struct {
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	Playlist struct {
		URL interface{} `yaml:"url"`
	} `yaml:"playlist"`
	Filters []Filter `yaml:"filters"`
}

type Filter struct {
	GroupTitle string    `yaml:"groupTitle"`
	Channels   []Channel `yaml:"channels"`
}

type Channel struct {
	Channel       string `yaml:"channel"`
	GuideName     string `yaml:"guideName"`
	ChannelNumber int    `yaml:"channelNumber"`
}
