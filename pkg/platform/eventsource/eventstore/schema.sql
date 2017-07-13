DROP TABLE IF EXISTS "stream";
DROP TABLE IF EXISTS "events";
DROP TABLE IF EXISTS "aggregates";

CREATE TABLE streams (
  id serial primary key not null,
  StreamID text NOT NULL,
  EventID uuid NOT NULL,
  EventNumber integer NOT NULL,
  EventType text NOT NULL,
  MetaData bytea NOT NULL,
  Data bytea NOT NULL,
  CreatedAt timestamp with time zone NOT NULL DEFAULT statement_timestamp() 
);

GRANT ALL ON SEQUENCE public.streams_id_seq TO esmadesimple;

GRANT ALL ON TABLE public.streams TO esmadesimple;

CREATE or REPLACE FUNCTION insert_streams(streamid text, eventid uuid, 
eventnumber int, eventtype text, metadata bytea, data bytea) RETURNS void AS $$
  BEGIN
	INSERT INTO streams(StreamID, EventID, EventNumber, EventType, MetaData, Data) 
			VALUES (streamid, eventid, eventnumber, eventtype, metadata, data)
  END;

$$ LANGUAGE plpgsql;

CREATE TABLE "events" (
  "id" serial primary key not null,
  "aggregate_id" uuid NOT NULL,
  "event_seq" integer NOT NULL,
  "msg_type" text NOT NULL,
    "msg_ver" smallint NOT NULL,
    "msg_data" jsonb NOT NULL,
  "inserted_at" timestamp with time zone NOT NULL DEFAULT statement_timestamp() 
);

CREATE TABLE "Aggregates" (
    "id" serial primary key not null,
       "aggregate_id" uuid NOT NULL, 
       "aggregate_type" text NOT NULL, 
       "snapshot_event_seq" integer NOT NULL
);

/*
CREATE INDEX "idx_streams_streamId" ON "streams" (uuid);
CREATE INDEX "idx_streams_uuid" ON "streams" (uuid);

CREATE INDEX "idx_events_type" ON "events" (type ASC);

CREATE INDEX "idx_events_uuid" ON "events" (uuid);

CREATE INDEX "idx_events_inserted_at" ON "events" (inserted_at DESC);
*/