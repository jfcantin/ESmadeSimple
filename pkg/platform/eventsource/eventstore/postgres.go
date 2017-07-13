package eventstore

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // required postgres db
)

const (
	insertStreamFormatQuery string = `INSERT INTO streams(StreamID, EventID, EventNumber, EventType, MetaData, Data) 
			VALUES ($1, $2, $3, $4, $5, $6);`

	selectStreamFormatQuery string = `SELECT streamid, eventid, eventnumber, eventtype, metadata, data 
								from streams 
								where streamid = $1 and eventnumber > $2 
								order by eventnumber;`
	versionStreamFormatQuery string = `SELECT count(0) as version 
	  							from streams 
								  where streamid = $1 
								  group by streamid;`
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
	if expectedVersion != ExpectedAny && expectedVersion < currentVersion {
		versionDiff := currentVersion - expectedVersion
		if len(events) != versionDiff {
			return fmt.Errorf("version mismatch. Expected: %d, but was: %d and the number of events differ", expectedVersion, currentVersion)
		}
		overlap := es.readStartingAt(streamName, expectedVersion)

		// for each position they need to match incoming to be good
		for i, evt := range overlap {
			if evt.EventID != events[i].ID {
				return fmt.Errorf("version mismatch. Expected: %d, but was: %d and not all events were already commited", expectedVersion, currentVersion)
			}
		}
		return nil
	}

	for _, e := range events {
		currentVersion++
		_, err := es.DB.Exec(insertStreamFormatQuery, streamName, e.ID, currentVersion, e.Type, e.MetaData, e.Data)
		if err != nil {
			log.Fatalf("failed to insert event: %v - %v", e, err)
		}
	}
	log.Printf("Current version is now: %d", es.getCurrentVersionForStream(streamName))
	return nil
}

func (es *Postgres) getCurrentVersionForStream(streamName string) int {
	var current int
	err := es.DB.QueryRow(versionStreamFormatQuery, streamName).Scan(&current)
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
	return es.readStartingAt(streamName, 0)
}

// ReadAll read all events from start to finish for a given stream.
func (es *Postgres) readStartingAt(streamName string, start int) []RecordedEvent {
	rows, err := es.DB.Query(selectStreamFormatQuery, streamName, start)
	if err != nil {
		log.Fatalf("Could not recover result from query %s - %v", selectStreamFormatQuery, err)
	}
	defer rows.Close()
	var events []RecordedEvent
	for rows.Next() {
		var e RecordedEvent
		var eventID []uint8
		err := rows.Scan(&e.StreamID, &eventID, &e.EventNumber,
			&e.EventType, &e.MetaData, &e.Data)
		if err != nil {
			log.Fatalf("Could not read rows: %v", err)
		}
		e.EventID = GUID(eventID)
		events = append(events, e)
	}
	return events
}
