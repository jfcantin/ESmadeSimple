package eventsource_test

import (
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

func TestAppendingToExistingStream(t *testing.T) {
	store := es.NewEventStore()

	eventID := es.NewGUID()
	testEvents := []es.EventData{
		es.EventData{eventID, "testEvent", false, []byte("some data"), []byte("some metadata")},
		es.EventData{eventID, "testEvent", false, []byte("some data 2"), []byte("some metadata 2")},
		es.EventData{eventID, "testEvent", false, []byte("some data 3"), []byte("some metadata 3")},
	}

	err := store.AppendToStream("test-stream", es.ExpectedAny, testEvents)
	if err != nil {
		t.Fatal(err)
	}

	events := store.ReadAllStreamEvents("test-stream")

	if len(events) != 3 {
		t.Fatalf("Expected 3 events stored but was: %v\n", len(events))
	}

	anotherEvent := es.EventData{eventID, "testEvent", false, []byte("some data 4"), []byte("some metadata 4")}

	err = store.AppendToStream("test-stream", 3, []es.EventData{anotherEvent})
	if err != nil {
		t.Fatal(err)
	}
	events = store.ReadAllStreamEvents("test-stream")

	if len(events) != 4 {
		t.Fatalf("Expected 4 events stored but was: %v\n", len(events))
	}
	event := events[3]
	t.Logf("event: %+v\n", event)
	if event.StreamID != "test-stream" ||
		event.EventID != eventID ||
		event.EventNumber != 4 ||
		event.EventType != "testEvent" ||
		string(event.Data) != "some data 4" ||
		string(event.MetaData) != "some metadata 4" ||
		time.Now().Sub(event.CreatedAt) < 0 {
		t.Errorf("RecordedEvent doesn't match EventData, expected: %+v\nBut was %+v", anotherEvent, event)
	}
}

func TestRetrieveingFromMultipleStream(t *testing.T) {
	store := es.NewEventStore()

	eventID := es.NewGUID()
	stream1 := []es.EventData{
		es.EventData{eventID, "testEvent", false, []byte("some data"), []byte("some metadata")},
	}
	stream2 := []es.EventData{
		es.EventData{eventID, "testEvent", false, []byte("some data"), []byte("some metadata")},
	}

	err := store.AppendToStream("test-stream", es.ExpectedAny, stream1)
	if err != nil {
		t.Fatal(err)
	}
	err = store.AppendToStream("test-stream2", es.ExpectedAny, stream2)
	if err != nil {
		t.Fatal(err)
	}

	eventStream1 := store.ReadAllStreamEvents("test-stream")
	eventStream2 := store.ReadAllStreamEvents("test-stream2")

	if len(eventStream1) != 1 {
		t.Fatalf("Expected 1 events stored but was: %v\n", len(eventStream1))
	}
	if len(eventStream1) != 1 {
		t.Fatalf("Expected 1 events stored but was: %v\n", len(eventStream2))
	}

	event := eventStream1[0]
	if event.StreamID != "test-stream" ||
		event.EventID != eventID ||
		event.EventNumber != 1 ||
		event.EventType != "testEvent" ||
		string(event.Data) != "some data" ||
		string(event.MetaData) != "some metadata" ||
		time.Now().Sub(event.CreatedAt) < 0 {
		t.Errorf("RecordedEvent doesn't match EventData, expected: %+v\nBut was %+v", stream1[0], event)
	}
	event = eventStream2[0]
	if event.StreamID != "test-stream2" ||
		event.EventID != eventID ||
		event.EventNumber != 1 ||
		event.EventType != "testEvent" ||
		string(event.Data) != "some data" ||
		string(event.MetaData) != "some metadata" ||
		time.Now().Sub(event.CreatedAt) < 0 {
		t.Errorf("RecordedEvent doesn't match EventData, expected: %+v\nBut was %+v", stream2[0], event)
	}
}

func TestAppendingANewEventWithSmallerExpextedVersionThanCurrentVersion(t *testing.T) {
	store := es.NewEventStore()

	eventID := es.NewGUID()
	testEvents := []es.EventData{
		es.EventData{eventID, "testEvent", false, []byte("some data"), []byte("some metadata")},
		es.EventData{eventID, "testEvent", false, []byte("some data 2"), []byte("some metadata 2")},
		es.EventData{eventID, "testEvent", false, []byte("some data 3"), []byte("some metadata 3")},
	}

	err := store.AppendToStream("test-stream", es.ExpectedAny, testEvents)
	if err != nil {
		t.Fatal(err)
	}

	anotherEvent := es.EventData{es.NewGUID(), "testEvent", false, []byte("some data 4"), []byte("some metadata 4")}

	err = store.AppendToStream("test-stream", 2, []es.EventData{anotherEvent})
	if err == nil {
		t.Fatalf("Expected concurrency error but returned nil")
	}
}

func TestAddingDuplicatedEvent(t *testing.T) {
	t.Skip("Not implemented")
}

func TestInParallelShouldNotFailEventIfEventStoreIsGlobalVariable(t *testing.T) {
	t.Skip("Not implemented")
}
