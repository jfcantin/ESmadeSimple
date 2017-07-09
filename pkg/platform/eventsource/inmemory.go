package eventsource

import (
	"fmt"
	"time"
)

// InMemory represents an in memory event store
type InMemory struct {
	store map[string][]RecordedEvent
}

// NewEventStore creates an in memory event store
func NewEventStore() *InMemory {
	var es InMemory
	es.store = make(map[string][]RecordedEvent, 0)
	return &es
}

func (es *InMemory) AppendToStream(streamName string, expectedVersion int, events []EventData) error {
	fmt.Printf("AppendToStream: %+v\n", events)
	currentVersion := len(es.store[streamName])
	for _, e := range events {
		currentVersion++
		es.store[streamName] = append(es.store[streamName], convertEventToRecorded(streamName, currentVersion, e))
	}
	fmt.Println("internal length: ", len(es.store[streamName]))
	return nil
}

func (es *InMemory) ReadAllStreamEvents(streamName string) []RecordedEvent {
	fmt.Printf("storage size: %d\n", len(es.store))
	fmt.Printf("requested id: %v\n", streamName)
	for k, v := range es.store {
		fmt.Printf("%v -> %+v\n", k, v)
	}

	return es.store[streamName]
}

func convertEventToRecorded(streamName string, version int, e EventData) RecordedEvent {
	return RecordedEvent{
		StreamID:    streamName,
		EventID:     e.ID,
		EventNumber: version,
		EventType:   e.Type,
		Data:        e.Data,
		MetaData:    e.MetaData,
		CreatedAt:   time.Now(),
	}
}
