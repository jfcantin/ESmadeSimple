package eventstore_test

import (
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"

	es "github.com/jfcantin/ESmadeSimple/pkg/platform/eventsource/eventstore"
)

type pgTest struct {
	*es.Postgres
}

func (pg *pgTest) ResetDB() {
	// stmt, err := pg.db.Prepare("truncate table streams")
	_, err := pg.DB.Exec("truncate table streams")
	if err != nil {
		log.Fatal(err)
	}
}

// NewPostgresStore creates a postgres event store
func newPostgresStore() es.ReadAppender {
	p, err := es.NewPostgresStore("esmadesimple", "esmadesimple")

	pg := pgTest{p.(*es.Postgres)}
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
		return nil //, fmt.Errorf("Could not open database: %v", err)
	}
	log.Print("reset db")
	pg.ResetDB()

	return &pg
}
func TestInMemoryEventStore(t *testing.T) {
	t.Parallel()
	stores := []struct {
		name  string
		store func() es.ReadAppender
	}{
		{"In Memory", es.NewInMemoryStore},
		{"Postgres", newPostgresStore},
	}
	for _, s := range stores {
		t.Run("Test can append and retrieve events for store: "+s.name, func(t *testing.T) { testCanAppendAndRetrieveAnEvent(t, s.store) })
		t.Run("Test can append to existing stream "+s.name, func(t *testing.T) { testAppendToExistingStream(t, s.store) })
		t.Run("Test can retrieve from multiple stream "+s.name, func(t *testing.T) { testRetrieveingFromMultipleStream(t, s.store) })
		t.Run("Test appending a new event with lower than expected version "+s.name, func(t *testing.T) { testAppendingANewEventWithSmallerExpextedVersionThanCurrentVersion(t, s.store) })
		t.Run("Test appending duplicated events with lower than expected version "+s.name, func(t *testing.T) { testAppendADuplicatedEventWithSmallerExpextedVersionThanCurrentVersion(t, s.store) })
		t.Run("Test appending some duplicate with lower than expeted version "+s.name, func(t *testing.T) {
			testAppendSomeDuplicatedEventWithSmallerExpextedVersionThanCurrentVersion(t, s.store)
		})
		t.Run("Test append with empty stream name "+s.name, func(t *testing.T) { testAppendWithoutStreamNameShouldError(t, s.store) })
	}
}

func testCanAppendAndRetrieveAnEvent(t *testing.T, storefunc func() es.ReadAppender) {
	store := storefunc()
	eventID := es.NewGUID()

	log.Println("Appending")
	err := store.Append("test-stream", es.ExpectedAny, []es.EventData{newTestEventWithGuid(eventID)})
	if err != nil {
		t.Error(err)
	}

	log.Println("Reading")
	events := store.ReadAll("test-stream")
	log.Println("Verifying")

	if len(events) != 1 {
		t.Fatalf("Expected only 1 event stored but was: %v\n", len(events))
	}
	event := events[0]
	checkEvent(t, newTestEventWithGuid(eventID), event)
}

func testAppendToExistingStream(t *testing.T, storefunc func() es.ReadAppender) {
	store := storefunc()
	err := configureStoreWithTestStream(store, 3)
	if err != nil {
		t.Fatal(err)
	}

	err = store.Append("test-stream", 3, []es.EventData{newTestEventWithGuid(es.NewGUID())})
	if err != nil {
		t.Fatal(err)
	}
	events := store.ReadAll("test-stream")

	if len(events) != 4 {
		t.Fatalf("Expected 4 events stored but was: %v\n", len(events))
	}
}

func testRetrieveingFromMultipleStream(t *testing.T, storefunc func() es.ReadAppender) {
	store := storefunc()
	eventID := es.NewGUID()
	stream1 := []es.EventData{
		newTestEventWithGuid(eventID),
	}
	stream2 := []es.EventData{
		newTestEventWithGuid(eventID),
	}

	err := store.Append("test-stream", es.ExpectedAny, stream1)
	if err != nil {
		t.Fatal(err)
	}
	err = store.Append("test-stream2", es.ExpectedAny, stream2)
	if err != nil {
		t.Fatal(err)
	}

	eventStream1 := store.ReadAll("test-stream")
	eventStream2 := store.ReadAll("test-stream2")

	if len(eventStream1) != 1 {
		t.Fatalf("Expected 1 events stored but was: %v\n", len(eventStream1))
	}
	if len(eventStream1) != 1 {
		t.Fatalf("Expected 1 events stored but was: %v\n", len(eventStream2))
	}
}

func testAppendingANewEventWithSmallerExpextedVersionThanCurrentVersion(t *testing.T, storefunc func() es.ReadAppender) {
	store := storefunc()
	err := configureStoreWithTestStream(store, 3)

	if err != nil {
		t.Fatal(err)
	}

	err = store.Append("test-stream", 2, []es.EventData{newTestEventWithGuid("")})
	if err == nil {
		t.Fatalf("Expected concurrency error but returned nil")
	}
}

func testAppendADuplicatedEventWithSmallerExpextedVersionThanCurrentVersion(t *testing.T, storefunc func() es.ReadAppender) {
	store := storefunc()
	err := configureStoreWithTestStream(store, 3)

	if err != nil {
		t.Fatal(err)
	}
	recoredEvents := store.ReadAll("test-stream")
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

	err = store.Append("test-stream", 1, events)
	if err != nil {
		t.Fatalf("Expected stream to be written without error: %v", err)
	}

	recoredEvents = store.ReadAll("test-stream")
	if len(recoredEvents) != 3 {
		t.Fatalf("Expected 3 events stored but was: %v\n", len(recoredEvents))
	}
}

func testAppendSomeDuplicatedEventWithSmallerExpextedVersionThanCurrentVersion(t *testing.T, storefunc func() es.ReadAppender) {
	store := storefunc()
	err := configureStoreWithTestStream(store, 3)

	if err != nil {
		t.Fatal(err)
	}
	recoredEvents := store.ReadAll("test-stream")
	events := []es.EventData{es.EventData{
		ID:       recoredEvents[2].EventID,
		Type:     recoredEvents[2].EventType,
		IsJSON:   false,
		Data:     recoredEvents[2].Data,
		MetaData: recoredEvents[2].MetaData},
		newTestEventWithGuid(""),
		newTestEventWithGuid("")}

	err = store.Append("test-stream", 1, events)
	if err == nil {
		t.Fatalf("Expected a version error, but got nil: %v", err)
	}

	recoredEvents = store.ReadAll("test-stream")
	if len(recoredEvents) != 3 {
		t.Fatalf("Expected 3 events stored but was: %v\n", len(recoredEvents))
	}
}

func testAppendWithoutStreamNameShouldError(t *testing.T, storefunc func() es.ReadAppender) {
	store := storefunc()
	err := store.Append("", es.ExpectedAny, []es.EventData{newTestEventWithGuid("")})
	if err == nil {
		t.Error("Expected missing stream name error")
	}
}

func configureStoreWithTestStream(store es.ReadAppender, numberOfEvent int) error {
	var testEvents = make([]es.EventData, 0)
	for i := 0; i < numberOfEvent; i++ {
		testEvents = append(testEvents, newTestEventWithGuid(es.NewGUID()))
	}

	err := store.Append("test-stream", es.ExpectedAny, testEvents)
	if err != nil {
		return fmt.Errorf("Could not append to stream")
	}
	return nil //*InMemory
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
