package eventsource_test

import (
	"fmt"
	"testing"

	es "github.com/jfcantin/ESmadeSimple/pkg/platform/eventsource"
)

func TestCanAppendAndRetrieveAnEvent(t *testing.T) {
	store := es.NewEventStore()

	eventID := es.NewGUID()

	err := store.AppendToStream("test-stream", es.ExpectedAny, []es.EventData{newTestEventWithGuid(eventID)})
	if err != nil {
		t.Error(err)
	}

	events := store.ReadAllStreamEvents("test-stream")

	if len(events) != 1 {
		t.Fatalf("Expected only 1 event stored but was: %v\n", len(events))
	}
	event := events[0]
	checkEvent(t, newTestEventWithGuid(eventID), event)
}

func TestAppendToExistingStream(t *testing.T) {
	store, err := configureStoreWithTestStream(3)
	if err != nil {
		t.Fatal(err)
	}

	err = store.AppendToStream("test-stream", 3, []es.EventData{newTestEventWithGuid(es.NewGUID())})
	if err != nil {
		t.Fatal(err)
	}
	events := store.ReadAllStreamEvents("test-stream")

	if len(events) != 4 {
		t.Fatalf("Expected 4 events stored but was: %v\n", len(events))
	}
}

func TestRetrieveingFromMultipleStream(t *testing.T) {
	store := es.NewEventStore()

	eventID := es.NewGUID()
	stream1 := []es.EventData{
		newTestEventWithGuid(eventID),
	}
	stream2 := []es.EventData{
		newTestEventWithGuid(eventID),
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
}

func TestAppendingANewEventWithSmallerExpextedVersionThanCurrentVersion(t *testing.T) {
	store, err := configureStoreWithTestStream(3)

	if err != nil {
		t.Fatal(err)
	}

	err = store.AppendToStream("test-stream", 2, []es.EventData{newTestEventWithGuid("")})
	if err == nil {
		t.Fatalf("Expected concurrency error but returned nil")
	}
}

func TestAppendADuplicatedEventWithSmallerExpextedVersionThanCurrentVersion(t *testing.T) {
	store, err := configureStoreWithTestStream(3)

	if err != nil {
		t.Fatal(err)
	}
	recoredEvents := store.ReadAllStreamEvents("test-stream")
	events := []es.EventData{es.EventData{
		ID:       recoredEvents[1].EventID,
		Type:     recoredEvents[1].EventType,
		IsJSON:   false,
		Data:     recoredEvents[1].Data,
		MetaData: recoredEvents[1].MetaData},
		es.EventData{
			ID:       recoredEvents[2].EventID,
			Type:     recoredEvents[2].EventType,
			IsJSON:   false,
			Data:     recoredEvents[2].Data,
			MetaData: recoredEvents[2].MetaData}}

	err = store.AppendToStream("test-stream", 1, events)
	if err != nil {
		t.Fatalf("Expected stream to be written without error: %v", err)
	}

	recoredEvents = store.ReadAllStreamEvents("test-stream")
	if len(recoredEvents) != 3 {
		t.Fatalf("Expected 3 events stored but was: %v\n", len(recoredEvents))
	}
}

func TestAppendSomeDuplicatedEventWithSmallerExpextedVersionThanCurrentVersion(t *testing.T) {
	store, err := configureStoreWithTestStream(3)

	if err != nil {
		t.Fatal(err)
	}
	recoredEvents := store.ReadAllStreamEvents("test-stream")
	events := []es.EventData{es.EventData{
		ID:       recoredEvents[2].EventID,
		Type:     recoredEvents[2].EventType,
		IsJSON:   false,
		Data:     recoredEvents[2].Data,
		MetaData: recoredEvents[2].MetaData},
		newTestEventWithGuid(""),
		newTestEventWithGuid("")}

	t.Logf("events: %+v\n", events)
	err = store.AppendToStream("test-stream", 1, events)
	if err == nil {
		t.Fatalf("Expected a version error, but got nil", err)
	}

	recoredEvents = store.ReadAllStreamEvents("test-stream")
	if len(recoredEvents) != 3 {
		t.Fatalf("Expected 3 events stored but was: %v\n", len(recoredEvents))
	}
}

func configureStoreWithTestStream(numberOfEvent int) (*es.InMemory, error) {

	store := es.NewEventStore()

	var testEvents = make([]es.EventData, 0)
	for i := 0; i < numberOfEvent; i++ {
		testEvents = append(testEvents, newTestEventWithGuid(es.NewGUID()))
	}

	err := store.AppendToStream("test-stream", es.ExpectedAny, testEvents)
	if err != nil {
		return nil, fmt.Errorf("Could not append to stream")
	}
	return store, nil //*InMemory
}

func TestAddingDuplicatedEvent(t *testing.T) {
	t.Skip("Not implemented")
}

func TestInParallelShouldNotFailEventIfEventStoreIsGlobalVariable(t *testing.T) {
	t.Skip("Not implemented")
}

func newTestEventWithGuid(eventID es.GUID) es.EventData {
	if eventID == "" {
		eventID = es.NewGUID()
	}
	return newTestEventWithData(eventID, "", "")
}
func newTestEventWithData(eventID es.GUID, eventData, metaData string) es.EventData {
	if eventData == "" {
		eventData = "some data"
	}
	if metaData == "" {
		metaData = "some metadata"
	}
	return es.EventData{eventID, "testEvent", false, []byte(eventData), []byte(metaData)}
}

func checkEvent(t *testing.T, expected es.EventData, got es.RecordedEvent) {
	if got.StreamID != "test-stream" {
		t.Errorf("Expected %v, but was %v", "test-stream", got.StreamID)
	}
	if got.EventID != expected.ID {
		t.Errorf("Expected %v, but was %v", expected.ID, got.EventID)
	}
	// if got.EventNumber != 1 { t.Errorf("Expected %v, but was %v", "test-stream", event.StreamID)}
	if got.EventType != expected.Type {
		t.Errorf("Expected %v, but was %v", expected.Type, got.EventType)
	}
	if string(got.Data) != string(expected.Data) {
		t.Errorf("Expected %v, but was %v", string(expected.Data), string(got.Data))
	}
	if string(got.MetaData) != string(expected.MetaData) {
		t.Errorf("Expected %v, but was %v", string(expected.MetaData), string(got.MetaData))
	}
	if got.EventType != expected.Type {
		t.Errorf("Expected %v, but was %v", expected.Type, got.EventType)
	}
	// time.Now().Sub(event.CreatedAt) < 0 {
}
