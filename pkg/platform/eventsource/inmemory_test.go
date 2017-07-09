package eventsource_test

import (
	"fmt"
	"testing"
	"time"

	es "github.com/jfcantin/ESmadeSimple/pkg/platform/eventsource"
)

func TestCanSaveAndRetrieveAnEvent(t *testing.T) {
	store := es.NewEventStore()

	eventID := es.NewGUID()
	testEvent := es.EventData{eventID, "testEvent", false, []byte("some data"), []byte("some metadata")}

	err := store.AppendToStream("test-stream", es.ExpectedAny, []es.EventData{testEvent})
	if err != nil {
		t.Error(err)
	}

	events := store.ReadAllStreamEvents("test-stream")

	if len(events) != 1 {
		t.Fatalf("Expected only 1 event stored but was: %v\n", len(events))
	}
	event := events[0]
	fmt.Printf("event: %+v\n", event)
	if event.StreamID != "test-stream" ||
		event.EventID != eventID ||
		event.EventNumber != 1 ||
		event.EventType != "testEvent" ||
		string(event.Data) != "some data" ||
		string(event.MetaData) != "some metadata" ||
		time.Now().Sub(event.CreatedAt) < 0 {
		t.Errorf("RecordedEvent doesn't match EventData, expected: %+v\nBut was %+v", testEvent, event)
	}
}

// func TestGetEventForAggregateWithMultipleEvent(t *testing.T) {
// 	es := NewEventStore()
// 	guid1 := newGuid()
// 	err := es.SaveEvent(guid1, []GuidVersionDescriptor{
// 		newTestEvent(guid1, 1), newTestEvent(guid1, 2)}, 2)
// 	if err != nil {
// 		t.Errorf("could not save Events for guid: %v\n - %v", guid1, err)
// 	}

// 	guid2 := newGuid()
// 	err = es.SaveEvent(guid2, []GuidVersionDescriptor{
// 		newTestEvent(guid2, 1), newTestEvent(guid2, 2), newTestEvent(guid2, 3)}, 3)

// 	if err != nil {
// 		t.Errorf("could not save Events for guid: %v\n - %v", guid1, err)
// 	}

// 	events := es.GetEventsForAggregate(guid1)
// 	if len(events) != 2 {
// 		t.Errorf("Wrong number of events expected %v, but was %v", 2, len(events))
// 	}

// 	events = es.GetEventsForAggregate(guid2)
// 	if len(events) != 3 {
// 		t.Errorf("Wrong number of events expected %v, but was %v", 3, len(events))
// 	}
// }
