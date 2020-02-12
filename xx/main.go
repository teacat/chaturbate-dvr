package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/grafov/m3u8"
	"github.com/parnurzeal/gorequest"
)

//
const chaturbateURL = "https://chaturbate.com/"

//
var retriesAfterOnlined = 0

//
var lastCheckOnline = time.Now()

//
type roomDossier struct {
	HLSSource string `json:"hls_source"`
}

//
func unescapeUnicode(raw string) string {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		panic(err)
	}
	return str
}

//
func getChannelURL(username string) string {
	return fmt.Sprintf("%s%s", chaturbateURL, username)
}

//
func getPlaylistURL(body string) string {
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

	return roomData.HLSSource
}

//
func getBody(username string) string {
	//
	resp, body, errs := gorequest.New().Get(getChannelURL(username)).End()
	return body
}

//
func getOnlineStatus(username string) bool {
	//
	body := getBody(username)
	return strings.Contains(body, "playlist.m3u8")

	//
	// resp, body, errs = gorequest.New().Get(url).End()
	// if resp.StatusCode == http.StatusForbidden {
	// 	return false
	// }
}

//
func getChunklistURL(playlistURL string) string {
	_, body, _ := gorequest.New().Get(playlistURL).End()

	p, listType, err := m3u8.DecodeFrom(strings.NewReader(body), true)
	if err != nil {
		panic(err)
	}
	if listType != m3u8.MEDIA {
		return
	}
	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		for _, v := range mediapl.Segments.URI {
			fmt.Printf("%s\n", v)
		}
	}

	//
	/*resp, body, errs := gorequest.New().Get(playlistURL).End()

	//
	var lines []string
	reader := bufio.NewReader(strings.NewReader(body))
	for {
		line, err := reader.ReadString('\n')
		lines = append(lines, line)
		if err != nil {
			break
		}
	}

	//
	baseURL := strings.TrimRight(playlistURL, "playlist.m3u8")
	//
	return fmt.Sprintf("%s%s", baseURL, lines[len(lines)-1])*/
}

func main() {
	username := "sexykiska"

	//
	for {
		//
		if !getOnlineStatus(username) {
			<-time.After(time.Minute * 3)
		}

		//
		body := getBody(username)
		//
		playlistURL := getPlaylistURL(body)
		//
		chunklistURL := getChunklistURL(playlistURL)
	}
}
