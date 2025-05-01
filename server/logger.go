package server

var Logger ILogger

type ILogger interface {
	Write(v string) error
}
