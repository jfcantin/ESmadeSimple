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

// AppendToStream append a list of EventData to a stream
func (es *InMemory) AppendToStream(streamName string, expectedVersion int, events []EventData) error {
	// TODO: Should synchronise access to the map in case more than one go routine
	// tries to read and write at the same time.
	stream := es.store[streamName]
	currentVersion := len(stream)
	// log.Printf("expected %d, current %d\n", expectedVersion, currentVersion)
	// log.Printf("stream content: %+v\n", stream)
	if expectedVersion != ExpectedAny && expectedVersion < currentVersion {
		return handleExpectedLowerThanCurrentVersion(stream, currentVersion, expectedVersion, events)
	}
	for _, e := range events {
		currentVersion++
		es.store[streamName] = append(es.store[streamName], convertEventToRecorded(streamName, currentVersion, e))
	}
	// fmt.Println("internal length: ", len(es.store[streamName]))
	return nil
}

func handleExpectedLowerThanCurrentVersion(stream []RecordedEvent, currentVersion int, expectedVersion int, events []EventData) error {
	// It is an all or nothing case.
	// 1. if all events have been recorded already we are good.
	// 2. if some but not all events have been recorded its no good.
	// 3. otherwise no good
	// log.Printf("handleExpected\n")
	versionDiff := currentVersion - expectedVersion
	// log.Printf("Version diff: %d\n", versionDiff)
	if len(events) != versionDiff {
		return fmt.Errorf("version mismatch. Expected: %d, but was: %d and not all events were already commited", expectedVersion, currentVersion)
	}
	overlap := stream[currentVersion-versionDiff:]
	// log.Printf("overlap: %+v\n", overlap)
	// log.Printf("newEvents: %+v\n", events)
	// for each position they need to match incoming to be good
	for i, evt := range overlap {
		if evt.EventID != events[i].ID {
			return fmt.Errorf("version mismatch. Expected: %d, but was: %d and not all events were already commited", expectedVersion, currentVersion)
		}
	}
	return nil
}

// ReadAllStreamEvents read all events from start to finish for a given stream.
func (es *InMemory) ReadAllStreamEvents(streamName string) []RecordedEvent {
	// for k, v := range es.store {
	// 	fmt.Printf("%v -> %+v\n", k, v)
	// }

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
