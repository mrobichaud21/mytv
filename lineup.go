package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const M3U8_CACHE_FILE string = ".cache.m3u8"

// lineup contains the state of the application.
type lineup struct {
	Sources         []Filter
	Scanning        bool
	HDChannels      map[int]hdHomeRunLineupItem
	providerChannel map[int]ProviderChannel
	playlistUrl     string
	rawFile         string
	Discovery       Discovery
}

// newLineup returns a new l ssineup for the given config struct.
func newLineup(config Configuration) *lineup {

	lineup := &lineup{
		HDChannels:      make(map[int]hdHomeRunLineupItem),
		providerChannel: make(map[int]ProviderChannel),
		playlistUrl:     config.Playlist.URL.(string),
		Sources:         config.Filters,
		Discovery:       config.Discovery,
	}

	return lineup
}

func (i *lineup) GetPlaylist() (*os.File, error) {

	if _, err := os.Stat(M3U8_CACHE_FILE); err == nil {

		file, fileErr := os.Open(M3U8_CACHE_FILE)
		if fileErr == nil {
			log.Info("Playlist using cache data.")
			return file, nil
		}

	}

	req, reqErr := http.NewRequest("GET", i.playlistUrl, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15")

	resp, err := http.Get(i.playlistUrl)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(M3U8_CACHE_FILE, body, 0777)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(M3U8_CACHE_FILE)
	if err == nil {
		log.Info("Playlist using cache data1.")

	}
	return file, nil
}

func (lineup *lineup) GetDiscoveryData() DiscoveryData {

	lineup.Discovery.Manufacturer = "Silicondust"
	lineup.Discovery.ModelNumber = "HDTC-2US"
	lineup.Discovery.FirmwareName = "hdhomeruntc_atsc"
	lineup.Discovery.FirmwareVersion = "20150826"

	return DiscoveryData{
		FriendlyName:    lineup.Discovery.FriendlyName,
		Manufacturer:    lineup.Discovery.Manufacturer,
		ModelNumber:     lineup.Discovery.ModelNumber,
		FirmwareName:    lineup.Discovery.FirmwareName,
		TunerCount:      lineup.Discovery.IptvStreams,
		FirmwareVersion: lineup.Discovery.FirmwareVersion,
		DeviceID:        lineup.Discovery.DeviceID,
		DeviceAuth:      lineup.Discovery.DeviceAuth,
		BaseURL:         fmt.Sprintf("http://%s:%d", lineup.Discovery.BaseAdress, lineup.Discovery.ServicePort),
		LineupURL:       fmt.Sprintf("http://%s:%d/lineup.json", lineup.Discovery.BaseAdress, lineup.Discovery.ServicePort),
	}
}

func (lineup *lineup) Scan() error {

	for _, filter := range lineup.Sources {

		playlist, err := lineup.GetPlaylist()

		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(playlist)

		// this is the straight M3U8 Output
		// we add the line if it begins with "EXTINF" and group-title matches
		// we add the channnel line using modulus
		filterChannels := []string{}
		// we set this value when want the channel so when we read the next line we attach the stream url
		currentChannel := 0
		//fmt.Println(currentChannel)

		var xChannels []Channel

		for scanner.Scan() {

			lineText := scanner.Text()

			//filterText := fmt.Sprintf("group-title=\"%s\"", filter.GroupTitle)
			//fmt.Sprintln(filterText)

			if strings.HasPrefix(lineText, "#EXTINF") && strings.Contains(lineText, fmt.Sprintf("group-title=\"%s\"", filter.GroupTitle)) {

				channelName := extractKeyValue(lineText, "tvg-name")

				if containsChannelName(filter.Channels, channelName) {
					filterChannels = append(filterChannels, lineText)
				}

				channel := getChannelMappingData(filter.Channels, channelName)
				if channel != nil {
					//filterChannels = append(filterChannels, lineText)
					tvgId := extractKeyValue(lineText, "tvg-id")
					tvgLogo := extractKeyValue(lineText, "tvg-logo")

					lineup.providerChannel[channel.ChannelNumber] = newProviderChannel(filter.GroupTitle, channelName, tvgId, tvgLogo, coalese(channel.GuideName, channelName), channel.ChannelNumber)
					currentChannel = channel.ChannelNumber
				}

				c := Channel{
					Channel:       channelName,
					ChannelNumber: 1000,
					GuideName:     "",
				}
				xChannels = append(xChannels, c)

			} else {
				if len(filterChannels)%2 == 1 {
					//filterChannels = append(filterChannels, lines+" url=\""+lineText+"\"")
					// First we get a "copy" of the entry
					if entry, ok := lineup.providerChannel[currentChannel]; ok {

						// Then we modify the copy
						entry.StreamURL = lineText

						// Then we reassign map entry
						lineup.providerChannel[currentChannel] = entry
						lineup.HDChannels[currentChannel] = newHDHRItem(lineup.providerChannel[currentChannel], lineup.Discovery)
					}

					filterChannels = append(filterChannels, lineText)
				}
			}
		}
		fileName := "./" + strings.Replace(filter.GroupTitle, "/", "", -1) + ".yaml"
		writeStructToYaml(fileName, xChannels)
		f, err := writePlaylist(filterChannels, filter.GroupTitle)

		if err == nil {
			filter.rawFile = f
		}

		fmt.Printf("Total channels = %d", len(filterChannels)/2)

	}

	_, _ = writePlaylist2(lineup.providerChannel, "all")

	return nil
}

func getChannelMappingData(channels []Channel, name string) *Channel {

	for _, s := range channels {
		if name == s.Channel {
			return &s
		}
	}
	//https://stackoverflow.com/questions/20240179/nil-detection-in-go
	channel := new(Channel)
	return channel
}

func containsChannelName(channels []Channel, name string) bool {

	for _, s := range channels {
		if name == s.Channel {
			return true
		}
	}
	return false
}

func extractKeyValue(line string, key string) string {
	var ch string

	if strings.Contains(line, key+"=\"") {
		ch = strings.SplitN(strings.SplitN(line, key+"=\"", 2)[1], "\" ", 2)[0]
	} else {
		ch = strings.TrimSpace(strings.SplitN(line, ",", 2)[1])
	}
	if ch == "" {
		ch = "No Name"
	}
	ch = strings.ReplaceAll(ch, "\"", "")
	return ch
}

func writeStructToYaml(fileName string, obj []Channel) {

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("error opening/creating file: %v", err)
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)

	err = enc.Encode(obj)
	if err != nil {
		log.Fatalf("error encoding: %v", err)
	}
}

func writePlaylist(data []string, groupTitle string) (filename string, err error) {

	fn := strings.Replace(groupTitle, "/", "", -1)
	fn = strings.Replace(fn, " ", "_", -1)
	fn = strings.Replace(fn, "|", "", -1)

	fileName := "./.cache/" + strings.ToLower(fn) + ".m3u8"

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Fprintln(f, "#EXTM3U")

	for _, v := range data {
		fmt.Fprintln(f, v)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("file written successfully")
	return fileName, nil
}

func writePlaylist2(data map[int]ProviderChannel, groupTitle string) (filename string, err error) {

	fn := strings.Replace(groupTitle, "/", "", -1)
	fn = strings.Replace(fn, " ", "_", -1)
	fn = strings.Replace(fn, "|", "", -1)

	fileName := "./.cache/" + strings.ToLower(fn) + "-1.m3u8"

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Fprintln(f, "#EXTM3U")

	for _, v := range data {

		lineText := fmt.Sprintf("#EXTINF:-1 tvg-id=\"%s\" tvg-name=\"%s\" tvg-logo=\"%s\" group-title=\"%s\",%s", v.TvgId, stripUnwanted(v.ChannelName), v.TvgLogo, stripUnwanted(v.GroupTitle), stripUnwanted(v.ChannelName))
		fmt.Fprintln(f, lineText)
		fmt.Fprintln(f, v.StreamURL)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("file written successfully")
	return fileName, nil
}

func stripUnwanted(str string) string {
	str = strings.Replace(str, "|", "", -1)
	str = strings.Replace(str, "US: ", "", -1)
	return str
}

func coalese(str string, str2 string) string {
	if str == "" {
		return str2
	}
	return str
}
