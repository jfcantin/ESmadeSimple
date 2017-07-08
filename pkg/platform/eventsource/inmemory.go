package eventsource

import "fmt"

// InMemory represents an in memory event store
type InMemory struct {
	current map[guid]GuidVersionDescriptor
}

type Event struct {
}

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

func (es *InMemory) SaveEvent(id guid, event GuidVersionDescriptor, expectedVersion int) error {
	fmt.Printf("event store: %+v\n", event)
	es.current[id] = event
	return nil
}

func (es *InMemory) GetEventForAggregate(id guid) GuidVersionDescriptor {
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
	es.current = make(map[guid]GuidVersionDescriptor, 0)
	return &es
}
