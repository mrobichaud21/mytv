package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

// func getPlaylist(url string) (io.ReadCloser, error) {

// 	if _, err := os.Stat("cache.m3u8"); err == nil {

// 		file, fileErr := os.Open("cache.m3u8")
// 		if fileErr == nil {
// 			log.Info("Playlist using cache data.")
// 			return file, nil
// 		}
// 	}

// 	req, reqErr := http.NewRequest("GET", url, nil)
// 	if reqErr != nil {
// 		return nil, reqErr
// 	}

// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15")

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	//body, err := io.ReadAll(resp.Body)
// 	// if strings.HasSuffix(strings.ToLower(url), ".gz") || resp.Header.Get("Content-Type") == "application/x-gzip" {
// 	// 	log.Infof("File (%s) is gzipp'ed, ungzipping now, this might take a while", url)
// 	// 	gz, gzErr := gzip.NewReader(resp.Body)
// 	// 	if gzErr != nil {
// 	// 		return nil, transport, gzErr
// 	// 	}
// 	// 	return gz, transport, nil
// 	// }

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = os.WriteFile("cache.m3u8", body, 0777)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return resp.Body, nil

// }

func getPlaylist(url string) (*os.File, error) {

	if _, err := os.Stat("cache.m3u8"); err == nil {

		file, fileErr := os.Open("cache.m3u8")
		if fileErr == nil {
			log.Info("Playlist using cache data.")
			return file, nil
		}
	}

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15")

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("cache.m3u8", body, 0777)
	if err != nil {
		panic(err)
	}

	file, err := os.Open("cache.m3u8")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return file, nil

}

func readLineByLine() {

	file, err := os.Open("cache.m3u8")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
