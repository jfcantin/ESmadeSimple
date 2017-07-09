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
	err := es.SaveEvent(id, []GuidVersionDescriptor{testEvent}, 1)
	if err != nil {
		t.Error(err)
	}

	events := es.GetEventsForAggregate(id)
	if len(events) != 1 {
		t.Error("Expected only 1 event stored")
	}
	event := events[0]
	fmt.Printf("event: %+v\n", event)
	if len(events) != 1 || event.Guid() != testEvent.Guid() || event.Version() != testEvent.Version() {
		t.Errorf("Expected id %v, but was %v", id, event.Guid())
		t.Errorf("Expected version %v, but was %v", testEvent.Version(), event.Version())
	}
}
func TestGetEventForAggregateWithMultipleEvent(t *testing.T) {
	es := NewEventStore()
	guid1 := newGuid()
	err := es.SaveEvent(guid1, []GuidVersionDescriptor{
		newTestEvent(guid1, 1), newTestEvent(guid1, 2)}, 2)
	if err != nil {
		t.Errorf("could not save Events for guid: %v\n - %v", guid1, err)
	}

	guid2 := newGuid()
	err = es.SaveEvent(guid2, []GuidVersionDescriptor{
		newTestEvent(guid2, 1), newTestEvent(guid2, 2), newTestEvent(guid2, 3)}, 3)

	if err != nil {
		t.Errorf("could not save Events for guid: %v\n - %v", guid1, err)
	}

	events := es.GetEventsForAggregate(guid1)
	if len(events) != 2 {
		t.Errorf("Wrong number of events expected %v, but was %v", 2, len(events))
	}

	events = es.GetEventsForAggregate(guid2)
	if len(events) != 3 {
		t.Errorf("Wrong number of events expected %v, but was %v", 3, len(events))
	}
}

func Test(t *testing.T) {

}
