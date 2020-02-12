package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/grafov/m3u8"
	"github.com/parnurzeal/gorequest"
)

// chaturbateURL is the base url of the website.
const chaturbateURL = "https://chaturbate.com/"

// retriesAfterOnlined tells the retries for stream when disconnected but not really offlined.
var retriesAfterOnlined = 0

// lastCheckOnline logs the last check time.
var lastCheckOnline = time.Now()

// buffer stores the media segments and wait for comsume.
var buffer = make(chan *m3u8.MediaSegment, 999999)

//
var bucket []string

//
var (
	errInternal = errors.New("err")
)

// roomDossier is the struct to parse the HLS source from the content body.
type roomDossier struct {
	HLSSource string `json:"hls_source"`
}

// unescapeUnicode escapes the unicode from the content body.
func unescapeUnicode(raw string) string {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		panic(err)
	}
	return str
}

// getChannelURL returns the full channel url to the specified user.
func getChannelURL(username string) string {
	return fmt.Sprintf("%s%s", chaturbateURL, username)
}

// getBody gets the channel page content body.
func getBody(username string) string {
	_, body, _ := gorequest.New().Get(getChannelURL(username)).End()
	return body
}

// getOnlineStatus check if the user is currently online by checking the playlist exists in the content body or not.
func getOnlineStatus(username string) bool {
	return strings.Contains(getBody(username), "playlist.m3u8")
}

// getHLSSource extracts the playlist url from the room detail page body.
func getHLSSource(body string) (string, string) {
	//
	r := regexp.MustCompile(`window\.initialRoomDossier = "(.*?)"`)
	matches := r.FindAllStringSubmatch(body, -1)

	//
	var roomData roomDossier
	data := unescapeUnicode(matches[0][1])
	err := json.Unmarshal([]byte(data), &roomData)
	if err != nil {
		panic(err)
	}

	return roomData.HLSSource, strings.TrimRight(roomData.HLSSource, "playlist.m3u8")
}

// parseHLSSource parses the HLS table and return the maximum resolution m3u8 source.
func parseHLSSource(url string, baseURL string) string {
	_, body, _ := gorequest.New().Get(url).End()

	//
	p, listType, _ := m3u8.DecodeFrom(strings.NewReader(body), true)
	if listType != m3u8.MASTER {
		return ""
	}

	master := p.(*m3u8.MasterPlaylist)
	return fmt.Sprintf("%s%s", baseURL, master.Variants[len(master.Variants)-1].URI)
}

//
func parseM3U8Source(url string) (chunks []*m3u8.MediaSegment, wait float64, err error) {
	_, body, _ := gorequest.New().Get(url).End()

	//
	p, listType, _ := m3u8.DecodeFrom(strings.NewReader(body), true)
	if listType != m3u8.MEDIA {
		return nil, 0, errInternal
	}

	media := p.(*m3u8.MediaPlaylist)
	wait = media.TargetDuration

	// Only fill with the real segments.
	for _, v := range media.Segments {
		if v == nil {
			continue
		}
		chunks = append(chunks, v)
	}
	return
}

//
func start(username string) {
	//
	for {
		// Check again after a while if the user is currently not online.
		if !getOnlineStatus(username) {
			log.Printf("%s is not online, check again after 3 minutes...", username)
			<-time.After(time.Minute * 3)
		}

		log.Printf("%s is online! Fetching the stream...", username)

		// Define the video filename by current time.
		filename := time.Now().String() + ".mp4"
		// Get the channel page content body.
		body := getBody(username)
		// Get the master playlist URL from extracting the channel body.
		hlsSource, baseURL := getHLSSource(body)
		// Get the best resolution m3u8 by parsing the HLS source table.
		m3u8Source := parseHLSSource(hlsSource, baseURL)
		//
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		go comsumer(f, baseURL)

		// Keep fetching the stream chunks until the playlist cannot be accessed after retried x times (which means the channel is offlined).
		for {
			// Get the chunks.
			chunks, wait, err := parseM3U8Source(m3u8Source)
			//
			if err != nil {
				if retriesAfterOnlined > 10 {
					log.Printf("Failed to fetch the video segments after retried, %s might be offlined.", username)
					retriesAfterOnlined = 0
					break
				} else {
					log.Printf("Failed to fetch the video segments, will try again. (%d/10)", retriesAfterOnlined)
					//
					retriesAfterOnlined++
					// Wait to fetch the next playlist.
					<-time.After(time.Duration(wait) * time.Second)
					continue
				}
			}
			for _, v := range chunks {
				var ignore bool
				for _, j := range bucket {
					if v.URI[len(v.URI)-10:] == j {
						ignore = true
						break
					}
				}
				if ignore {
					continue
				}
				bucket = append(bucket, v.URI[len(v.URI)-10:])
				log.Printf("%s (%d in buffer)", v.URI, len(buffer))
				buffer <- v
			}
			//
			//buffer = append(buffer, chunks...)
			//log.Printf("Storing chunks.")
			// Append the chunks to the video file.
			//appendChunks(f, baseURL, chunks)
			// Wait to fetch the next playlist.
			<-time.After(time.Duration(wait) * time.Second)
		}
	}
}

func comsumer(file *os.File, baseURL string) {
	for {
		v := <-buffer
		_, body, _ := gorequest.New().Get(fmt.Sprintf("%s%s", baseURL, v.URI)).EndBytes()

		if _, err := file.Write(body); err != nil {
			panic(err)
		}
	}
}

// appendChunks appends the streaming chunks data into a single video file.
//func appendChunks(file *os.File, baseURL string, chunks []*m3u8.MediaSegment) {
//	for _, v := range chunks {
//		_, body, _ := gorequest.New().Get(fmt.Sprintf("%s%s", baseURL, v.URI)).EndBytes()
//
//		log.Println(v.URI)
//
//		if _, err := file.Write(body); err != nil {
//			panic(err)
//		}
//	}
//}

func main() {
	username := "yesonee"

	start(username)
}
