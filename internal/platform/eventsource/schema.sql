DROP TABLE IF EXISTS "events";
DROP TABLE IF EXISTS "aggregates";

CREATE TABLE "events" (
  "id" serial primary key not null,
  "aggregate_id" uuid NOT NULL,
  "event_seq" integer NOT NULL,
  "msg_type" text NOT NULL,
    "msg_ver" smallint NOT NULL,
    "msg_data" jsonb NOT NULL,
  "inserted_at" timestamp(6) NOT NULL DEFAULT statement_timestamp()
    
           [aggregate_id] [uniqueidentifier] NOT NULL, 
       [event_seq] [int] NOT NULL, 
       [msg_type] [nvarchar](256) NOT NULL, 
       [msg_ver] [smallint] NOT NULL, 
       [msg_data] [varbinary](max) NOT NULL, 
);

CREATE TABLE "Aggregates" (
    "id" serial primary key not null,
       "aggregate_id" uuid NOT NULL, 
       "aggregate_type" text NOT NULL, 
       "snapshot_event_seq" integer NOT NULL, 
)

CREATE INDEX "idx_events_type" ON "events" (type ASC);

CREATE INDEX "idx_events_uuid" ON "events" (uuid);

CREATE INDEX "idx_events_inserted_at" ON "events" (inserted_at DESC);
