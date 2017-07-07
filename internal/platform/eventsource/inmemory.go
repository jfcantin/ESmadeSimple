package eventsource

import "fmt"

// InMemory represents an in memory event store
type InMemory struct {
	current map[guid]GuidVersionDescriptor
}

type Event struct {
}

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

func (s *InMemory) SaveEvent(id guid, event GuidVersionDescriptor, expectedVersion int) error {
	fmt.Printf("event store: %+v\n", event)
	s.current[id] = event
	return nil
}

func (s *InMemory) GetEventForAggregate(id guid) GuidVersionDescriptor {
	fmt.Printf("storage size: %d\n", len(s.current))
	fmt.Printf("requested id: %v\n", id)
	for k, v := range s.current {
		fmt.Printf("%v -> %+v\n", k, v)
	}

	return s.current[id]
}

// NewEventStore creates an in memory event store
func NewEventStore() *InMemory {
	var es InMemory
	es.current = make(map[guid]GuidVersionDescriptor, 0)
	return &es
}
