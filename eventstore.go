package ycq

import (
	"encoding/json"
	"github.com/jetbasrawi/go.cqrs/internal/uuid"
	"reflect"
	"time"
)

type StreamReader interface {
	Next() bool
	Err() error
	Scan(e interface{}, m interface{})
	EventType() string
	Version() int
	NextVersion(version int)
}

type StreamWriter interface {

	// Append writes an event to the head of the stream.
	//
	// If the stream does not exist, it will be created.
	//
	// There are some special version numbers that can be provided.
	// http://docs.geteventstore.com/http-api/3.7.0/writing-to-a-stream/
	//
	// -2 : The write should never conflict with anything and should always succeed.
	//
	// -1 : The stream should not exist at the time of writing. This write will create it.
	//
	// 0 : The stream should exist but it should be empty.
	Append(expectedVersion *int, events ...*Event) error
}


type EventStore interface {
	StreamReader(streamName string) StreamReader
	StreamWriter(streamName string) StreamWriter
}

// Event encapsulates the data of an domain event.
//
// EventStreamID is the id returned in the event atom response.
// EventNumber represents the stream version for this event.
// EventType describes the event type.
// EventID is the guid of the event.
// Data contains the data of the event.
// Links contains the urls of the event on the evenstore
// MetaData contains the metadata for the event.
type Event struct {
	EventStreamID string      `json:"eventStreamId,omitempty"`
	EventNumber   int         `json:"eventNumber,omitempty"`
	EventType     string      `json:"eventType,omitempty"`
	EventID       string      `json:"eventId,omitempty"`
	Data          interface{} `json:"data"`
	Links         []Link      `json:"links,omitempty"`
	MetaData      interface{} `json:"metadata,omitempty"`
}

// PrettyPrint renders an indented json view of the Event object.
func (e *Event) PrettyPrint() string {
	b, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Link encapsulates url data for events.
type Link struct {
	URI      string `json:"uri"`
	Relation string `json:"relation"`
}

// TimeStr is a type used to format feed dates.
type TimeStr string

// Time returns a TimeStr version of the time.Time argument t.
func Time(t time.Time) TimeStr {
	return TimeStr(t.Format("2006-01-02T15:04:05-07:00"))
}

// NewEvent creates a new event object.
//
// If an empty eventId is provided a new uuid will be generated automatically
// and retured in the event.
// If an empty eventType is provided the eventType will be set to the
// name of the type provided.
// data and meta can be nil.
func NewEvent(eventID, eventType string, data interface{}, meta interface{}) *Event {
	e := &Event{}

	e.EventID = eventID
	if eventID == "" {
		e.EventID = NewUUID()
	}

	e.EventType = eventType
	if eventType == "" {
		e.EventType = typeOf(data)
	}

	e.Data = data
	e.MetaData = meta
	return e
}

// NewUUID returns a new V4 uuid as a string.
func NewUUID() string {
	return uuid.NewV4().String()
}

// typeOf is a helper to get the names of types.
func typeOf(i interface{}) string {
	if i == nil {
		return ""
	}
	return reflect.TypeOf(i).Elem().Name()
}

