package manager

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	streamDomain "github.com/teacat/chaturbate-dvr/domain/stream"
)

type stream struct {
	username           string
	channelURL         string
	isPaused           bool
	isOnline           bool
	resolution         int
	resolutionFallback streamDomain.ResolutionFallback
	framerate          int
	splitByDuration    int
	splitByFilesize    int
	session            *streamSession
	chUpdate           chan<- *streamDomain.StreamUpdateDTO
	chOutput           chan<- *streamDomain.StreamOutputDTO
}

type streamSession struct {
	buffer        map[int][]byte
	bufferIndex   int
	file          *os.File
	retries       int
	resolution    int
	framerate     int
	durationTotal int
	durationQuota int
	filesizeTotal int
	filesizeQuota int
}

func newStream(username string) (*stream, chan<- *streamDomain.StreamUpdateDTO, chan<- *streamDomain.StreamOutputDTO) {
	chUpd := make(chan *streamDomain.StreamUpdateDTO)
	chOut := make(chan *streamDomain.StreamOutputDTO)

	return &stream{
		username:   username,
		channelURL: "https://chaturbate.com/" + username,
		// TODO: resolution, framerate split, duration split, filesize split
		isPaused: true,
		chUpdate: chUpd,
		chOutput: chOut,
	}, chUpd, chOut
}

func (s *stream) start() {
	for {
		body, err := s.retrieveChannel()
		if err != nil {
			s.log("Error occurred while retrieving channel webpage: %s", err)
		}
		if s.isOnline {
			s.log("%s is now online.", s.username)
			if err := s.startRecording(body); err != nil { // blocking
				s.log("Error occurred while start recording: %s", err)
			}
			continue
		}
		s.log("%s went offline.", s.username)
		<-time.After(time.Minute * time.Duration(1)) // 1 minute cooldown
	}
}

func (s *stream) startRecording(body string) error {
	folder := fmt.Sprintf("./videos/%s", s.username)
	if err := os.MkdirAll(folder, 0777); err != nil {
		return fmt.Errorf("create folder: %w", err)
	}

	basename := fmt.Sprintf("./videos/%s/%s_%s", s.username, time.Now().Format("2006-01-02_15-04-05"))
	file, err := os.OpenFile(basename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	s.log("The video will be saved as %s.ts", basename)
	s.session = &streamSession{
		buffer:      make(map[int][]byte),
		bufferIndex: 0,
		retries:     0,
		file:        file,
	}

	//
	streamURI, err := parseHLS(s.resolution, s.framerate, s.resolutionFallback, body)
	if err != nil {
		return fmt.Errorf("parse hls: %w", err)
	}

	s.concatStreams()

	s.retrieveStream(streamURI)
}

func (s *stream) retrieveStream(uri string) {

}

func (s *stream) pause() {
}

func (s *stream) stop() {
}

func (s *stream) resume() {
}

func (s *stream) retrieveChannel() (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(s.channelURL)
	if err != nil {
		return "", fmt.Errorf("client get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	s.isOnline = strings.Contains(string(body), "playlist.m3u8")

	return string(body), nil
}

func (s *stream) log(message string, v ...interface{}) {
	s.chOutput <- &streamDomain.StreamOutputDTO{
		Username: s.username,
		Output:   "[" + time.Now().Format("2006-01-02 15:04:05") + "] " + fmt.Sprintf(message, v...),
	}

}
