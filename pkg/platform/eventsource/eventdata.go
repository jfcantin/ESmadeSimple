package eventsource

import (
	"fmt"
	"time"
)

const (
	ExpectedAny int = -1
)

// EventData represent an event that need to be stored
type EventData struct {
	ID       GUID
	Type     string
	IsJSON   bool
	Data     []byte
	MetaData []byte
}

func (e EventData) String() string {
	return fmt.Sprintf("EventData: {ID: %v, Type: %s, IsJSON: %v, Data: '%s', MetaData: '%s'}\n",
		e.ID, e.Type, e.IsJSON, string(e.Data), string(e.MetaData))
}

// RecordedEvent represent an event that has been stored
type RecordedEvent struct {
	StreamID    string
	EventID     GUID
	EventNumber int
	EventType   string
	Data        []byte
	MetaData    []byte
	CreatedAt   time.Time
}

func (e RecordedEvent) String() string {
	return fmt.Sprintf("RecordedEvent: {StreamID: %s, EventID: %v, EventNumber: %d, EventType: %s, Data: '%s', MetaData: '%s', CreatedAt: %v}\n",
		e.StreamID, e.EventID, e.EventNumber, e.EventType, string(e.Data), string(e.MetaData), e.CreatedAt)
}
