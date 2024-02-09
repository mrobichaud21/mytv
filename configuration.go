package main

// https://zhwt.github.io/yaml-to-go/
// https://github.com/go-yaml/yaml/issues/165
type Configuration struct {
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`

	Playlist struct {
		URL interface{} `yaml:"url"`
	} `yaml:"playlist"`
	Filters []Filter `yaml:"filters"`

	Discovery Discovery `yaml:"discovery"`
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

type Discovery struct {
	DeviceID        string `yaml:"deviceId"`
	DeviceUUID      string `yaml:"deviceUUID"`
	FriendlyName    string `yaml:"friendlyName"`
	Manufacturer    string `yaml:"manufacturer"`
	ModelNumber     string `yaml:"modelNumber"`
	FirmwareName    string `yaml:"firmwareName"`
	FirmwareVersion string `yaml:"firmwareVersion"`
	IptvStreams     int    `yaml:"iptvStreams"`
	DeviceAuth      string `yaml:"deviceAuth"`
	BaseAdress      string `yaml:"baseAdress"`
}

func (s *Discovery) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawDiscovery Discovery
	raw := rawDiscovery{
		Manufacturer:    "Silicondust",
		ModelNumber:     "HDTC-2US",
		FirmwareName:    "hdhomeruntc_atsc",
		FirmwareVersion: "20150826",
	} // Put your defaults here
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*s = Discovery(raw)
	return nil
}
