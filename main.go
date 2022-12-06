package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/TwiN/go-color"

	"github.com/grafov/m3u8"
	"github.com/parnurzeal/gorequest"
	"github.com/urfave/cli/v2"
)

// chaturbateURL is the base url of the website.
const chaturbateURL = "https://chaturbate.com/"

// retriesAfterOnlined tells the retries for stream when disconnected but not really offlined.
var retriesAfterOnlined = 0

// temp stores the used segment to prevent fetched the duplicates.
var temp []string

// segmentIndex is current stored segment index.
var segmentIndex int

// segmentMap is the map stores temporary video segments, it will be merged into master video file then got deleted.
var segmentMap map[string][]byte = make(map[string][]byte)

// stripLimit reprsents the maximum Bytes sizes to split the video into chunks.
var stripLimit int

// stripQuota represents how many Bytes left til the next video chunk stripping.
var stripQuota int

// path save video
const savePath = "video"

// error/message handler
var (
	errInternal         = errors.New("err")
	errNoUsername       = errors.New("recording: channel username required `-u [USERNAME]` option")
	errSegRetFail       = color.Colorize(color.Red, ("[FAILED] to fetch the video segments after retried, %s might went offline or is in ticket/privat show."))
	errSegRetFailOnline = color.Colorize(color.Red, ("[FAILED] to fetch the video segments, will try again. [%d/10]"))
	infoIsOnline        = color.Colorize(color.Green, ("[RECORDING] %s is online! start fetching.."))
	infoBackOnline      = color.Colorize(color.Green, ("[INFO] %s is back online!"))
	infoMergeSegment    = color.Colorize(color.Green, ("[INFO] inserting %d segment to the master file. [total: %d]"))
	infoSkipped         = color.Colorize(color.Blue, ("[INFO] skipped %s due to the empty body!\n"))
	infoNotOnline       = color.Colorize(color.Gray, ("[INFO] %s is not online, check again in %d minute(s)"))
	warningSegment      = color.Colorize(color.Yellow, ("[WARNING] cannot find segment %d, will try again. [%d/5]"))
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
	resp, body, _ := gorequest.New().Get(getChannelURL(username)).End()
	if resp.StatusCode != 200 {
		return ""
	}
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
	resp, body, _ := gorequest.New().Get(url).End()
	if resp.StatusCode == 403 {
		return ""
	}
	p, _, _ := m3u8.DecodeFrom(strings.NewReader(body), true)
	master, ok := p.(*m3u8.MasterPlaylist)
	if !ok {
		return ""
	}
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
	media, ok := p.(*m3u8.MediaPlaylist)
	if !ok {
		return nil, 3, errInternal
	}
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
	// Define the video filename by current time //04.09.22 added username into filename mK33y.
	filename := username + "_" + time.Now().Format("2006-01-02_15-04-05")
	var m3u8Source, baseURL, hlsSource string
	var tried int
	for {
		tried++
		//
		if tried > 10 {
			panic(errors.New("cannot fetch the Playlist correctly after 10 tries"))
		}
		// Get the channel page content body.
		body := getBody(username)
		//
		if body == "" {
			continue
		}
		// Get the master playlist URL from extracting the channel body.
		hlsSource, baseURL = getHLSSource(body)
		// Get the best resolution m3u8 by parsing the HLS source table.
		m3u8Source = parseHLSSource(hlsSource, baseURL)
		//
		if m3u8Source != "" {
			break
		}
		<-time.After(time.Millisecond * 500)
	}
	// Create the master video file.
	masterFile, _ := os.OpenFile("./"+savePath+"/"+filename+".ts", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	//
	log.Printf("the video will be saved as \"./"+savePath+"/%s\".", filename+".ts")

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
				log.Printf(errSegRetFail, username)
				break
			} else {
				log.Printf(errSegRetFailOnline, retriesAfterOnlined)
				retriesAfterOnlined++
				// Wait to fetch the next playlist.
				<-time.After(time.Duration(wait*1000) * time.Millisecond)
				continue
			}
		}
		if retriesAfterOnlined != 0 {
			log.Printf(infoBackOnline, username)
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
	for _, v := range temp {
		if URI[len(URI)-10:] == v {
			return true
		}
	}
	temp = append(temp, URI[len(URI)-10:])
	return false
}

// combineSegment combines the segments to the master video file in the background.
// fixed segment problems mK33y.
// still needs some attention here
func combineSegment(master *os.File, filename string) {
	index := 1
	stripIndex := 1
	var retry int
	<-time.After(4 * time.Second)

	for {
		<-time.After(300 * time.Millisecond)

		if index >= segmentIndex {
			<-time.After(1 * time.Second)
			continue
		}

		if _, ok := segmentMap[fmt.Sprintf("./%s/%s~%d.ts", savePath, filename, index)]; !ok {
			if retry >= 5 {
				index++
				retry = 0
				continue
			}
			if retry != 0 {
				log.Printf(warningSegment, index, retry)
			}
			retry++
			<-time.After(time.Duration(1*retry) * time.Second)
			continue
		}
		if retry != 0 {
			retry = 0
		}
		//
		b := segmentMap[fmt.Sprintf("./%s/%s~%d.ts", savePath, filename, index)]
		//
		if stripLimit != 0 && stripQuota <= 0 {
			newMasterFilename := "./" + savePath + "/" + filename + "_" + strconv.Itoa(stripIndex) + ".ts"
			master, _ = os.OpenFile(newMasterFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
			log.Printf("exceeded the specified stripping limit, creating new video file. (file: %s)", newMasterFilename)
			stripQuota = stripLimit
			stripIndex++
		}
		master.Write(b)
		//
		log.Printf(infoMergeSegment, index, segmentIndex)

		delete(segmentMap, fmt.Sprintf("./%s/%s~%d.ts", savePath, filename, index))
		index++
	}
}

// fetchSegment fetches the segment and append to the master file.
func fetchSegment(master *os.File, segment *m3u8.MediaSegment, baseURL string, filename string, index int) {
	_, body, _ := gorequest.New().Get(fmt.Sprintf("%s%s", baseURL, segment.URI)).EndBytes()
	log.Printf("fetching %s (size: %d)\n", segment.URI, len(body))
	if len(body) == 0 {
		log.Printf(infoSkipped, segment.URI)
		return
	}
	stripQuota -= len(body)
	segmentMap[fmt.Sprintf("./%s/%s~%d.ts", savePath, filename, index)] = body
}

// endpoint implements the application main function endpoint.
func endpoint(c *cli.Context) error {
	if c.String("username") == "" {
		log.Fatal(errNoUsername)
	}
	// Converts `strip` from MiB to Bytes
	stripLimit = c.Int("strip") * 1024 * 1024
	stripQuota = c.Int("strip") * 1024 * 1024
	//

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

	// Mkdir video folder
	if _, err := os.Stat("./" + savePath); os.IsNotExist(err) {
		os.Mkdir("./"+savePath, 0777)
	}
	//
	if c.Int("strip") != 0 {
		log.Printf("specifying stripping limit as %d MiB(s)", c.Int("strip"))
	}

	for {
		// Capture the stream if the user is currently online.
		if getOnlineStatus(c.String("username")) {
			log.Printf(infoIsOnline, c.String("username"))
			capture(c.String("username"))
			segmentIndex = 0
			temp = []string{}
			retriesAfterOnlined = 0
			continue
		}
		// Otherwise we keep checking the channel status until the user is online.
		log.Printf(infoNotOnline, c.String("username"), c.Int("interval"))
		<-time.After(time.Minute * time.Duration(c.Int("interval")))
	}
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
			&cli.IntFlag{
				Name:    "strip",
				Aliases: []string{"s"},
				Value:   0,
				Usage:   "MB sizes to split the video into chunks",
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
