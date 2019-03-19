package cqrs

import "fmt"

// ErrCommandExecution is the error returned in response to a failed command.
type ErrCommandExecution struct {
	Command CommandMessage
	Reason  string
}

// Error fulfills the error interface.
func (e *ErrCommandExecution) Error() string {
	return fmt.Sprintf("Invalid Operation. Command: %s Reason: %s", e.Command.CommandType(), e.Reason)
}

// ErrConcurrencyViolation is returned when a concurrency error is raised by the event store
// when events are persisted to a stream and the version of the stream does not match
// the expected version.
type ErrConcurrencyViolation struct {
	Aggregate       Aggregate
	ExpectedVersion int
	StreamName      string
}

func (e *ErrConcurrencyViolation) Error() string {
	return fmt.Sprintf("ConcurrencyError: StreamName: %s ExpectedVersion: %d StreamName: %s", e.Aggregate.AggregateId(), e.ExpectedVersion, e.StreamName)
}

// ErrUnauthorized is returned when a request to the repository is not authorized
type ErrUnauthorized struct{}

func (e *ErrUnauthorized) Error() string {
	return "Not authorized."
}

// ErrUnexpected is returned for all errors that are not otherwise represented
// explicitly.
//
// The original error is available for inspection in the Err field.
type ErrUnexpected struct {
	Err error
}

func (e *ErrUnexpected) Error() string {
	return fmt.Sprintf("An unepected error occurred. %s", e.Err)
}

// ErrRepositoryUnavailable is returned when the domain is temporarily unavailable
type ErrRepositoryUnavailable struct{}

func (e *ErrRepositoryUnavailable) Error() string {
	return "The repository is temporarily unavailable."
}

// ErrAggregateNotFound error returned when an aggregate was not found in the repository.
type ErrAggregateNotFound struct {
	AggregateID   string
	AggregateType string
}

func (e *ErrAggregateNotFound) Error() string {
	return fmt.Sprintf("Could not find any aggregate of type %s with AggregateId %s",
		e.AggregateType,
		e.AggregateID)
}

// ErrNoMoreEvents is returned when there are no events to return
// from a request to a stream.
type ErrNoMoreEvents struct{}

func (e ErrNoMoreEvents) Error() string {
	return "There are no more events to load."
}

// ErrNotFound is returned when a stream is not found.
type ErrNotFound struct {
	ErrorResponse string
}

func (e ErrNotFound) Error() string {
	return "The stream does not exist."
}

// ErrDeleted is returned when a request is made to a stream that
// has been hard deleted.
type ErrDeleted struct {
	ErrorResponse string
}

func (e ErrDeleted) Error() string {
	return "The stream has was deleted."
}

// ErrTemporarilyUnavailable is returned when the eventstore is temporarily unavailable.
type ErrTemporarilyUnavailable struct {
	ErrorResponse string
}

func (e ErrTemporarilyUnavailable) Error() string {
	return "Server Is Not Ready"
}

// ErrBadRequest is returned when the server returns a bad request error
type ErrBadRequest struct {
	ErrorResponse string
}

func (e ErrBadRequest) Error() string {
	return "Bad request."
}

type ErrNoConfiguredEventHandler struct {
	eventType string
	handler   string
}

func (e ErrNoConfiguredEventHandler) Error() string {
	return fmt.Sprintf("No configuration for: %s in handler: %s", e.eventType, e.handler)
}
