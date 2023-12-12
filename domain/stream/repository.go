package stream

type Manager interface {
	ListStreams() ([]*StreamDTO, error)
	AddStream(username string, resFallback ResolutionFallback, resolution, framerate, splitByFilesize, splitByDuration int, isPaused bool) error
	PauseStream(username string) error
	StopStream(username string) error
	ResumeStream(username string) error
	SubscribeStreams(chUpd chan<- *StreamUpdateDTO, chOut chan<- *StreamOutputDTO) error
}
