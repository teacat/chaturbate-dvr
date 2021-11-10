package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/teacat/pathx"

	"github.com/grafov/m3u8"
	"github.com/parnurzeal/gorequest"
	"github.com/urfave/cli/v2"
)

// chaturbateURL is the base url of the website.
const chaturbateURL = "https://chaturbate.com/"

// retriesAfterOnlined tells the retries for stream when disconnected but not really offlined.
var retriesAfterOnlined = 0

// lastCheckOnline logs the last check time.
var lastCheckOnline = time.Now()

// bucket stores the used segment to prevent fetched the duplicates.
var bucket []string

// segmentIndex is current stored segment index.
var segmentIndex int

//
var (
	errInternal   = errors.New("err")
	errNoUsername = errors.New("chaturbate-dvr: channel username required with `-u [username]` argument")
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
	// Get the room data from the page body.
	r := regexp.MustCompile(`window\.initialRoomDossier = "(.*?)"`)
	matches := r.FindAllStringSubmatch(body, -1)

	// Extract the data and get the HLS source URL.
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

	// Decode the HLS table.
	p, _, _ := m3u8.DecodeFrom(strings.NewReader(body), true)
	master := p.(*m3u8.MasterPlaylist)
	return fmt.Sprintf("%s%s", baseURL, master.Variants[len(master.Variants)-1].URI)
}

// parseM3U8Source gets the current segment list, the channel might goes offline if 403 was returned.
func parseM3U8Source(url string) (chunks []*m3u8.MediaSegment, wait float64, err error) {
	resp, body, errs := gorequest.New().Get(url).End()
	// Retry after 3 seconds if the connection lost or status code returns 403 (the channel might went offline).
	if len(errs) > 0 || resp.StatusCode == http.StatusForbidden {
		return nil, 3, errInternal
	}

	// Decode the segment table.
	p, _, _ := m3u8.DecodeFrom(strings.NewReader(body), true)
	media := p.(*m3u8.MediaPlaylist)
	wait = media.TargetDuration / 1.5

	// Ignore the empty segments.
	for _, v := range media.Segments {
		if v != nil {
			chunks = append(chunks, v)
		}
	}
	return
}

// capture captures the specified channel streaming.
func capture(username string) {
	// Define the video filename by current time.
	filename := time.Now().Format("2006-01-02_15-04-05")
	// Get the channel page content body.
	body := getBody(username)
	// Get the master playlist URL from extracting the channel body.
	hlsSource, baseURL := getHLSSource(body)
	// Get the best resolution m3u8 by parsing the HLS source table.
	m3u8Source := parseHLSSource(hlsSource, baseURL)
	// Create the master video file.
	masterFile, _ := os.OpenFile(filename+".ts", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	//
	log.Printf("the video will be saved as \"%s\".", filename+".ts")

	go combineSegment(masterFile, filename)
	watchStream(m3u8Source, username, masterFile, filename, baseURL)
}

// watchStream watches the stream and ends if the channel went offline.
func watchStream(m3u8Source string, username string, masterFile *os.File, filename string, baseURL string) {
	// Keep fetching the stream chunks until the playlist cannot be accessed after retried x times.
	for {
		// Get the chunks.
		chunks, wait, err := parseM3U8Source(m3u8Source)
		// Exit the fetching loop if the channel went offline.
		if err != nil {
			if retriesAfterOnlined > 10 {
				log.Printf("failed to fetch the video segments after retried, %s might went offline.", username)
				break
			} else {
				log.Printf("failed to fetch the video segments, will try again. (%d/10)", retriesAfterOnlined)
				retriesAfterOnlined++
				// Wait to fetch the next playlist.
				<-time.After(time.Duration(wait*1000) * time.Millisecond)
				continue
			}
		}
		if retriesAfterOnlined != 0 {
			log.Printf("%s is back online!", username)
			retriesAfterOnlined = 0
		}
		for _, v := range chunks {
			// Ignore the duplicated chunks.
			if isDuplicateSegment(v.URI) {
				continue
			}
			segmentIndex++
			go fetchSegment(masterFile, v, baseURL, filename, segmentIndex)
		}
		<-time.After(time.Duration(wait*1000) * time.Millisecond)
	}
}

// isDuplicateSegment returns true if the segment is already been fetched.
func isDuplicateSegment(URI string) bool {
	for _, v := range bucket {
		if URI[len(URI)-10:] == v {
			return true
		}
	}
	bucket = append(bucket, URI[len(URI)-10:])
	return false
}

// combineSegment combines the segments to the master video file in the background.
func combineSegment(master *os.File, filename string) {
	index := 1
	var retry int
	<-time.After(4 * time.Second)

	for {
		<-time.After(300 * time.Millisecond)

		if index >= segmentIndex {
			<-time.After(1 * time.Second)
			continue
		}

		if !pathx.Exists(fmt.Sprintf("%s~%d.ts", filename, index)) {
			if retry >= 5 {
				index++
				retry = 0
				continue
			}
			if retry != 0 {
				log.Printf("cannot find segment %d, will try again. (%d/5)", index, retry)
			}
			retry++
			<-time.After(time.Duration(1*retry) * time.Second)
			continue
		}
		if retry != 0 {
			retry = 0
		}
		//
		b, _ := ioutil.ReadFile(fmt.Sprintf("%s~%d.ts", filename, index))
		master.Write(b)
		log.Printf("inserting %d segment to the master file. (total: %d)", index, segmentIndex)
		//
		os.Remove(fmt.Sprintf("%s~%d.ts", filename, index))
		//
		index++
	}
}

// fetchSegment fetches the segment and append to the master file.
func fetchSegment(master *os.File, segment *m3u8.MediaSegment, baseURL string, filename string, index int) {
	_, body, _ := gorequest.New().Get(fmt.Sprintf("%s%s", baseURL, segment.URI)).EndBytes()
	log.Printf("fetching %s (size: %d)\n", segment.URI, len(body))
	if len(body) == 0 {
		log.Printf("skipped %s due to the empty body!\n", segment.URI)
		return
	}
	//
	f, err := os.OpenFile(fmt.Sprintf("%s~%d.ts", filename, index), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	if _, err := f.Write(body); err != nil {
		panic(err)
	}
}

// endpoint implements the application main function endpoint.
func endpoint(c *cli.Context) error {
	if c.String("username") == "" {
		log.Fatal(errNoUsername)
	}

	fmt.Println(" .o88b. db   db  .d8b.  d888888b db    db d8888b. d8888b.  .d8b.  d888888b d88888b")
	fmt.Println("d8P  Y8 88   88 d8' `8b `~~88~~' 88    88 88  `8D 88  `8D d8' `8b `~~88~~' 88'")
	fmt.Println("8P      88ooo88 88ooo88    88    88    88 88oobY' 88oooY' 88ooo88    88    88ooooo")
	fmt.Println("8b      88~~~88 88~~~88    88    88    88 88`8b   88~~~b. 88~~~88    88    88~~~~~")
	fmt.Println("Y8b  d8 88   88 88   88    88    88b  d88 88 `88. 88   8D 88   88    88    88.")
	fmt.Println(" `Y88P' YP   YP YP   YP    YP    ~Y8888P' 88   YD Y8888P' YP   YP    YP    Y88888P")
	fmt.Println("d8888b. db    db d8888b.")
	fmt.Println("88  `8D 88    88 88  `8D")
	fmt.Println("88   88 Y8    8P 88oobY'")
	fmt.Println("88   88 `8b  d8' 88`8b")
	fmt.Println("88  .8D  `8bd8'  88 `88.")
	fmt.Println("Y8888D'    YP    88   YD")
	fmt.Println("---")

	for {
		// Capture the stream if the user is currently online.
		if getOnlineStatus(c.String("username")) {
			log.Printf("%s is online! fetching...", c.String("username"))
			capture(c.String("username"))
			segmentIndex = 0
			bucket = []string{}
			retriesAfterOnlined = 0
			continue
		}
		// Otherwise we keep checking the channel status until the user is online.
		log.Printf("%s is not online, check again after %d minute(s)...", c.String("username"), c.Int("interval"))
		<-time.After(time.Minute * time.Duration(c.Int("interval")))
	}
	return nil
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Value:   "",
				Usage:   "channel username to watching",
			},
			&cli.IntFlag{
				Name:    "interval",
				Aliases: []string{"i"},
				Value:   1,
				Usage:   "minutes to check if a channel goes online or not",
			},
		},
		Name:   "chaturbate-dvr",
		Usage:  "watching a specified chaturbate channel and auto saves the stream as local file",
		Action: endpoint,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
