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
	Sources     []Filter
	Scanning    bool
	channels    map[int]hdHomeRunLineupItem
	playlistUrl string
}

// newLineup returns a new lineup for the given config struct.
func newLineup(config Configuration) *lineup {

	lineup := &lineup{
		channels:    make(map[int]hdHomeRunLineupItem),
		playlistUrl: config.Playlist.URL.(string),
		Sources:     config.Filters,
	}

	return lineup
}

func containsChannelName(channels []Channel, name string) bool {

	for _, s := range channels {
		if name == s.Channel {
			return true
		}
	}
	return false
}

// Scan processes all sources.
func (l *lineup) Scan() error {

	for _, filter := range l.Sources {

		playlist, err := l.getPlaylist()

		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(playlist)

		// this is the straight M3U8 Output
		// we add the line if it begins with "EXTINF" and group-title matches
		// we add the channnel line using modulus
		filterChannels := []string{}

		var xChannels []Channel

		for scanner.Scan() {

			lineText := scanner.Text()

			filterText := fmt.Sprintf("group-title=\"%s\"", filter.GroupTitle)

			fmt.Sprintln(filterText)

			if strings.HasPrefix(lineText, "#EXTINF") && strings.Contains(lineText, fmt.Sprintf("group-title=\"%s\"", filter.GroupTitle)) {

				channelName := extractChannelName(lineText, "tvg-name")

				if containsChannelName(filter.Channels, channelName) {
					filterChannels = append(filterChannels, lineText)
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
					filterChannels = append(filterChannels, lineText)
				}
			}
		}
		fileName := "./" + strings.Replace(filter.GroupTitle, "/", "", -1) + ".yaml"
		writeStructToYaml(fileName, xChannels)

		fmt.Printf("Total channels = %d", len(filterChannels)/2)

		//playlist.Close()
		//
	}
	return nil
}

func (i *lineup) getPlaylist() (*os.File, error) {

	if _, err := os.Stat(M3U8_CACHE_FILE); err == nil {

		file, fileErr := os.Open("cache.m3u8")
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

	err = os.WriteFile("cache.m3u8", body, 0777)
	if err != nil {
		return nil, err
	}

	file, err := os.Open("cache.m3u8")
	if err == nil {
		log.Info("Playlist using cache data1.")

	}
	return file, nil
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
