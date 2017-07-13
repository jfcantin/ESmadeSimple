package eventstore

import uuid "github.com/satori/go.uuid"

// GUID represent a GUID
type GUID string

// NewGUID creates a new guid with uuid.NewV4
func NewGUID() GUID {
	return GUID(uuid.NewV4().String())
}
