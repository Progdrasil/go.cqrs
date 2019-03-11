// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package ycq

import (
	"errors"
)

type Repository struct {
	eventStore         EventStore
	eventBus           EventBus
	streamNameDelegate StreamNamer
	aggregateFactory   AggregateFactory
	eventFactory       EventMaker
}

func NewRepository(eventStore EventStore) (*Repository, error) {
	if (eventStore == nil){
		return nil, errors.New("a valid eventstore is required.")
	}

	r := &Repository{eventStore: eventStore}
	return r, nil
}

// SetAggregateFactory sets the aggregate factory that should be used to
// instantiate aggregate instances
//
// Only one AggregateFactory can be registered at any one time.
// Any registration will overwrite the provious registration.
func (r *Repository) SetAggregateFactory(factory AggregateFactory) {
	r.aggregateFactory = factory
}

// SetEventFactory sets the event factory that should be used to instantiate event
// instances.
//
// Only one event factory can be set at a time. Any subsequent registration will
// overwrite the previous factory.
func (r *Repository) SetEventFactory(factory EventMaker) {
	r.eventFactory = factory
}

// SetStreamNameDelegate sets the stream name delegate
func (r *Repository) SetStreamNameDelegate(delegate StreamNamer) {
	r.streamNameDelegate = delegate
}


