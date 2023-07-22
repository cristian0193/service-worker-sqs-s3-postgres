package domain

import "go.uber.org/zap"

// Event represents a process.
type Event struct {
	TrackID         string
	File            string
	Bucket          string
	FileSize        int64
	OriginalEvent   interface{}
	Log         	*zap.SugaredLogger
	Origin          string
	ReceivedAt      string
	Filename        string

}

// Source represents a source of events.
type Source interface {
	Consume() <-chan *Event
	Processed(e *Event) error
	Close() error
}
