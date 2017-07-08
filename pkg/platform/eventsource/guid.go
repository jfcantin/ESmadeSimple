package eventsource

import uuid "github.com/satori/go.uuid"

type guid string

func newGuid() guid {
	return guid(uuid.NewV4().String())
}
