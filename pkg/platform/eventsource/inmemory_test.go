package eventsource

import (
	"fmt"
	"testing"
)

type testEvent struct {
	id      guid
	version int
}

func (t testEvent) Guid() guid {
	return t.id

}

func (t testEvent) Version() int {
	return t.version
}

func (t testEvent) EventDescriptor() EventDescriptor {
	return EventDescriptor{t.id, t.version, Event{}}
}

func newTestEvent(id guid, version int) GuidVersionDescriptor {
	return testEvent{id, version}
}

func TestCanSaveAndRetrieveAnEvent(t *testing.T) {
	es := NewEventStore()

	id := newGuid()
	testEvent := newTestEvent(id, 1)
	err := es.SaveEvent(id, testEvent, 1)
	if err != nil {
		t.Error(err)
	}

	event := es.GetEventForAggregate(id)
	fmt.Printf("event: %+v\n", event)
	if event.Guid() != testEvent.Guid() || event.Version() != testEvent.Version() {
		t.Errorf("Expected id %v, but was %v", id, event.Guid())
		t.Errorf("Expected version %v, but was %v", testEvent.Version(), event.Version())
	}
}
