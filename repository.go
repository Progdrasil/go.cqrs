// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package ycq

import (
	"errors"
	"fmt"
	"net/url"
)

// DomainRepository is the interface that all domain repositories should implement.
type DomainRepository interface {
	//Loads an aggregate of the given type and ID
	Load(aggregateTypeName string, aggregateID string) (AggregateRoot, error)

	//Saves the aggregate.
	Save(aggregate AggregateRoot, expectedVersion *int) error

	// Set the function that will instantiate the for this repo
	SetAggregateFactory(factory AggregateFactory)

	// Set the functions that will instantiate the events
	SetEventFactory(factory EventMaker)

	// Set the function that will provide the name for aggregates in this repo
	SetStreamNameDelegate(delegate StreamNamer)
}

type repository struct {
	eventStore         EventStore
	eventBus           EventBus
	streamNameDelegate StreamNamer
	aggregateFactory   AggregateFactory
	eventFactory       EventMaker
}

func NewRepository(eventStore EventStore) (*repository, error) {
	if (eventStore == nil){
		return nil, errors.New("a valid eventstore is required.")
	}

	r := &repository{eventStore:eventStore}
	return r, nil
}

// SetAggregateFactory sets the aggregate factory that should be used to
// instantiate aggregate instances
//
// Only one AggregateFactory can be registered at any one time.
// Any registration will overwrite the provious registration.
func (r *repository) SetAggregateFactory(factory AggregateFactory) {
	r.aggregateFactory = factory
}

// SetEventFactory sets the event factory that should be used to instantiate event
// instances.
//
// Only one event factory can be set at a time. Any subsequent registration will
// overwrite the previous factory.
func (r *repository) SetEventFactory(factory EventMaker) {
	r.eventFactory = factory
}

// SetStreamNameDelegate sets the stream name delegate
func (r *repository) SetStreamNameDelegate(delegate StreamNamer) {
	r.streamNameDelegate = delegate
}

// Load will load all events from a stream and apply those events to an aggregate
// of the type specified.
//
// The aggregate type and id will be passed to the configured StreamNamer to
// get the stream name.
func (r *repository) Load(aggregateType, id string) (AggregateRoot, error) {

	if r.aggregateFactory == nil {
		return nil, fmt.Errorf("The common domain repository has no Aggregate Factory.")
	}

	if r.streamNameDelegate == nil {
		return nil, fmt.Errorf("The common domain repository has no stream name delegate.")
	}

	if r.eventFactory == nil {
		return nil, fmt.Errorf("The common domain has no Event Factory.")
	}

	aggregate := r.aggregateFactory.GetAggregate(aggregateType, id)
	if aggregate == nil {
		return nil, fmt.Errorf("The repository has no aggregate factory registered for aggregate type: %s", aggregateType)
	}

	streamName, err := r.streamNameDelegate.GetStreamName(aggregateType, id)
	if err != nil {
		return nil, err
	}

	stream := r.eventStore.StreamReader(streamName)
	for stream.Next() {
		switch err := stream.Err().(type) {
		case nil:
			break
		case *url.Error, ErrTemporarilyUnavailable:
			return nil, &ErrRepositoryUnavailable{}
		case ErrNoMoreEvents:
			return aggregate, nil
		case ErrUnauthorized:
			return nil, &ErrUnauthorized{}
		case ErrNotFound:
			return nil, &ErrAggregateNotFound{AggregateType: aggregateType, AggregateID: id}
		default:
			return nil, &ErrUnexpected{Err: err}
		}

		event := r.eventFactory.MakeEvent(stream.EventType())

		//TODO: No test for meta
		meta := make(map[string]string)
		stream.Scan(event, &meta)
		if stream.Err() != nil {
			return nil, stream.Err()
		}
		em := NewEventMessage(id, event, Int(stream.Version()))
		for k, v := range meta {
			em.SetHeader(k, v)
		}
		aggregate.Apply(em, false)
		aggregate.IncrementVersion()
	}

	return aggregate, nil

}

// Save persists an aggregate
func (r *repository) Save(aggregate AggregateRoot, expectedVersion *int) error {

	if r.streamNameDelegate == nil {
		return fmt.Errorf("The common domain repository has no stream name delagate.")
	}

	resultEvents := aggregate.GetChanges()

	streamName, err := r.streamNameDelegate.GetStreamName(typeOf(aggregate), aggregate.AggregateID())
	if err != nil {
		return err
	}

	if len(resultEvents) > 0 {

		evs := make([]*Event, len(resultEvents))

		for k, v := range resultEvents {
			//TODO: There is no test for this code
			v.SetHeader("AggregateID", aggregate.AggregateID())
			evs[k] = NewEvent("", v.EventType(), v.Event(), v.GetHeaders())
		}

		streamWriter := r.eventStore.StreamWriter(streamName)
		err := streamWriter.Append(expectedVersion, evs...)
		switch e := err.(type) {
		case nil:
			break
		case ErrConcurrencyViolation:
			return &ErrConcurrencyViolation{Aggregate: aggregate, ExpectedVersion: expectedVersion, StreamName: streamName}
		case ErrUnauthorized:
			return &ErrUnauthorized{}
		case ErrTemporarilyUnavailable:
			return &ErrRepositoryUnavailable{}
		default:
			return &ErrUnexpected{Err: e}
		}
	}

	aggregate.ClearChanges()

	for k, v := range resultEvents {
		if expectedVersion == nil {
			r.eventBus.PublishEvent(v)
		} else {
			em := NewEventMessage(v.AggregateID(), v.Event(), Int(*expectedVersion+k+1))
			r.eventBus.PublishEvent(em)
		}
	}

	return nil
}

