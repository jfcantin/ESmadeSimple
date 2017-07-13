package eventstore

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // required postgres db
)

// Postgres represents the postgres client
type Postgres struct {
	DB *sql.DB
}

// NewPostgresStore creates a postgres event store
func NewPostgresStore(user, dbname string) (ReadAppender, error) {
	var pg Postgres
	db, err := sql.Open("postgres", "user="+user+" dbname="+dbname+" sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("Could not open database: %v", err)
	}
	pg.DB = db
	return &pg, nil
}

// Append a list of EventData to a stream
func (es *Postgres) Append(streamName string, expectedVersion int, events []EventData) error {
	if streamName == "" {
		return fmt.Errorf("Missingn stream name")
	}
	currentVersion := es.getCurrentVersionForStream(streamName)
	log.Printf("expected %d, current %d\n", expectedVersion, currentVersion)
	if expectedVersion != ExpectedAny && expectedVersion < currentVersion {
		// return handleExpectedLowerThanCurrentVersion(stream, currentVersion, expectedVersion, events)

		versionDiff := currentVersion - expectedVersion
		log.Printf("Version diff: %d\n", versionDiff)
		if len(events) != versionDiff {
			return fmt.Errorf("version mismatch. Expected: %d, but was: %d and the number of events differ", expectedVersion, currentVersion)
		}
		overlap := es.readStartingAt(streamName, expectedVersion)
		// log.Printf("overlap: %+v\n", overlap)
		// log.Printf("newEvents: %+v\n", events)
		// for each position they need to match incoming to be good
		log.Printf("events: %+v\nOverlap: %+v", events, overlap)
		for i, evt := range overlap {
			if evt.EventID != events[i].ID {
				log.Printf("event: %+v\nOverlap: %+v", events[i], evt)
				return fmt.Errorf("version mismatch. Expected: %d, but was: %d and not all events were already commited.", expectedVersion, currentVersion)
			}
		}
		return nil
	}

	// VALUES (:streamID, :eventID, :eventNumber, :eventType, :metaData, :data)
	query := fmt.Sprintf(`INSERT INTO %s(StreamID, EventID, EventNumber, EventType, MetaData, Data) 
			VALUES ($1, $2, $3, $4, $5, $6)
		`, "streams")
	stmt, err := es.DB.Prepare(query)
	if err != nil {
		log.Fatalf("failed preparing insert statement: %s - %v", query, err)
	}
	for _, e := range events {
		currentVersion++
		_, err := stmt.Exec(streamName, e.ID, currentVersion, e.Type, e.MetaData, e.Data)
		if err != nil {
			log.Fatalf("failed to insert event: %v - %v", e, err)
		}
	}
	log.Printf("Current version is now: %d", es.getCurrentVersionForStream(streamName))
	return nil
}

func (es *Postgres) getCurrentVersionForStream(streamName string) int {
	var current int
	query := fmt.Sprintf("SELECT count(0) as version from streams where streamid = '%s' group by streamid", streamName)
	err := es.DB.QueryRow(query).Scan(&current)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No stream found with id: %s", streamName)
	case err != nil:
		log.Fatalf("Could not query row: %v", err)
	default:
		log.Printf("Version found was %v", current)
	}
	return current
}

// ReadAll read all events from start to finish for a given stream.
func (es *Postgres) ReadAll(streamName string) []RecordedEvent {
	query := fmt.Sprint("SELECT * from streams where streamid = $1")
	stmt, err := es.DB.Prepare(query)
	if err != nil {
		log.Fatalf("preparing statement: %s - %v", query, err)
	}
	rows, err := stmt.Query(streamName)
	if err != nil {
		log.Fatalf("Could not recover result from query %s - %v", query, err)
	}
	defer rows.Close()
	var events []RecordedEvent
	for rows.Next() {
		var e RecordedEvent
		var id int
		var eventID []uint8
		var createdAt time.Time
		err := rows.Scan(&id, &e.StreamID, &eventID, &e.EventNumber,
			&e.EventType, &e.MetaData, &e.Data, &createdAt)
		if err != nil {
			log.Fatalf("Could not read rows: %v", err)
		}
		e.EventID = GUID(eventID)
		events = append(events, e)
	}
	return events
}

// ReadAll read all events from start to finish for a given stream.
func (es *Postgres) readStartingAt(streamName string, start int) []RecordedEvent {
	query := fmt.Sprint("SELECT * from streams where streamid = $1 and eventnumber > $2")
	stmt, err := es.DB.Prepare(query)
	if err != nil {
		log.Fatalf("preparing statement: %s - %v", query, err)
	}
	rows, err := stmt.Query(streamName, start)
	if err != nil {
		log.Fatalf("Could not recover result from query %s - %v", query, err)
	}
	defer rows.Close()
	var events []RecordedEvent
	for rows.Next() {
		var e RecordedEvent
		var id int
		var eventID []uint8
		var createdAt time.Time
		err := rows.Scan(&id, &e.StreamID, &eventID, &e.EventNumber,
			&e.EventType, &e.MetaData, &e.Data, &createdAt)
		if err != nil {
			log.Fatalf("Could not read rows: %v", err)
		}
		e.EventID = GUID(eventID)
		events = append(events, e)
	}
	return events
}
