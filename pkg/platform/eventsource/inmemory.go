package eventsource

import "fmt"

// InMemory represents an in memory event store
type InMemory struct {
	current map[guid][]GuidVersionDescriptor
}

type Event struct {
}

// TODO: Investigate if event descriptor is required, why not just event?

// EventDescriptor represents the format being stored in the event store
type EventDescriptor struct {
	ID        guid
	Version   int
	EventData Event
}

// GuidVersionDescriptor represents something that is guidable, versionable and descriptable
type GuidVersionDescriptor interface {
	Guid() guid
	Version() int
	EventDescriptor() EventDescriptor
}

func (es *InMemory) SaveEvent(id guid, events []GuidVersionDescriptor, expectedVersion int) error {
	fmt.Printf("event store: %+v\n", events)
	es.current[id] = events
	return nil
}

func (es *InMemory) GetEventsForAggregate(id guid) []GuidVersionDescriptor {
	fmt.Printf("storage size: %d\n", len(es.current))
	fmt.Printf("requested id: %v\n", id)
	for k, v := range es.current {
		fmt.Printf("%v -> %+v\n", k, v)
	}

	return es.current[id]
}

// NewEventStore creates an in memory event store
func NewEventStore() *InMemory {
	var es InMemory
	es.current = make(map[guid][]GuidVersionDescriptor, 0)
	return &es
}
