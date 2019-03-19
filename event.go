// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cqrs

// EventMessage is the interface that an event must implement.
type EventMessage interface {

	// StreamName returns the ID of the Aggregate that the event relates to
	AggregateId() *AggregateId

	// GetHeaders returns the key value collection of headers for the event.
	//
	// Headers are metadata about the event that do not form part of the
	// actual event but are still required to be persisted alongside the event.
	GetHeaders() map[string]interface{}

	// SetHeader sets the value of the header specified by the key
	SetHeader(string, interface{})

	// Returns the actual event which is the payload of the event message.
	Event() interface{}

	// EventType returns a string descriptor of the command name
	EventType() string

	// number returns the version of the event
	Version() int

	// If the event has been persisted to the event store will be true otherwise false
	Committed() bool

	// Sets the event committed status to true.
	// committed means it has been recorded in the event store
	// Events can not be uncommitted
	Commit()
}

// EventDescriptor is an implementation of the event message interface.
type EventDescriptor struct {
	id        *AggregateId
	event     interface{}
	headers   map[string]interface{}
	version   int
	committed bool
}

// NewEventMessage returns a new event descriptor
func NewEventMessage(aggregateID *AggregateId, event interface{}, version int) *EventDescriptor {
	return &EventDescriptor{
		id:        aggregateID,
		event:     event,
		headers:   make(map[string]interface{}),
		version:   version,
		committed: false,
	}
}

// EventType returns the name of the event type as a string.
func (t *EventDescriptor) EventType() string {
	return typeOf(t.event)
}

// AggregateId returns the AggregateId
func (t *EventDescriptor) AggregateId() *AggregateId {
	return t.id
}

// GetHeaders returns the headers for the event.
func (t *EventDescriptor) GetHeaders() map[string]interface{} {
	return t.headers
}

// SetHeader sets the value of the header specified by the key
func (t *EventDescriptor) SetHeader(key string, value interface{}) {
	t.headers[key] = value
}

// Event the event payload of the event message
func (t *EventDescriptor) Event() interface{} {
	return t.event
}

// number returns the version of the event
func (t *EventDescriptor) Version() int {
	return t.version
}

func (t *EventDescriptor) Committed() bool {
	return t.committed
}

func (t *EventDescriptor) Commit() {
	t.committed = true
}
