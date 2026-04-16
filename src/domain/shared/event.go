package shared

// CloudEventEnvelope defines the standard metadata wrapper for all domain events.
type CloudEventEnvelope struct {
	// ID uniquely identifies the event.
	ID string `json:"id"`
	// Source identifies the context in which an event happened.
	Source string `json:"source"`
	// SpecVersion identifies the CloudEvents specification version.
	SpecVersion string `json:"specversion"`
	// Type describes the type of event related to the originating occurrence.
	Type string `json:"type"`
	// DataContentType describes the content type of the data value.
	DataContentType string `json:"datacontenttype"`
	// Subject identifies the specific subject of the event.
	Subject string `json:"subject,omitempty"`
	// Time represents when the occurrence happened.
	Time string `json:"time,omitempty"`
}

// DomainEvent is the interface that all domain events must satisfy.
type DomainEvent interface {
	// Type returns the CloudEvent type (e.g., com.carddemo.account.opened).
	Type() string
	// AggregateID returns the ID of the aggregate that generated this event.
	AggregateID() string
}
