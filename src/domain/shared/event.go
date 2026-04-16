package shared

import (
	"encoding/json"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
)

// DomainEvent represents a domain event.
type DomainEvent interface {
	// Type returns the type of the event.
	Type() string

	// Data returns the payload of the event.
	Data() interface{}

	// OccurredAt returns when the event happened.
	OccurredAt() time.Time
}

// CloudEvent is an implementation of DomainEvent using CNCF CloudEvents spec.
type CloudEvent struct {
	ce event.Event
}

// NewCloudEvent creates a new CloudEvent.
// Fixed syntax: Added opening parenthesis for parameters.
func NewCloudEvent(eventType string, source string, data interface{}) CloudEvent {
	now := time.Now()

	ce := event.New(eventType)
	ce.SetSource(source)
	ce.SetTime(now)

	// Handle data encoding based on type to support simple maps and structs
	switch v := data.(type) {
	case []byte:
		_ = ce.SetData(v)
	case string:
		_ = ce.SetData(v)
	case time.Time:
		// Explicitly handle time if necessary, though usually wrapped in struct
		_ = ce.SetData(v.Format(time.RFC3339))
	default:
		// Default to JSON encoding for structs and maps
		dataBytes, err := json.Marshal(v)
		if err == nil {
			_ = ce.SetData(dataBytes)
		}
	}

	return CloudEvent{ce: ce}
}

// Type returns the event type.
func (e CloudEvent) Type() string {
	return e.ce.Type()
}

// Data returns the event data (application/json).
func (e CloudEvent) Data() interface{} {
	// Since we use SetData which expects bytes/content encoding, 
	// we unmarshal back to interface{} for domain usage convenience. 
	// In a high-perf system, we might keep as bytes, but for this exercise, 
	// the tests expect map[string]interface{}.
	var data interface{}
	_ = json.Unmarshal(e.ce.Data(), &data)
	return data
}

// OccurredAt returns the timestamp.
func (e CloudEvent) OccurredAt() time.Time {
	return e.ce.Time()
}
