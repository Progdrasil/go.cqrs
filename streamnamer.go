// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cqrs

import (
	"fmt"
)

// StreamNamer is the interface that stream name delegates should implement.
// Takes an aggregate and returns a streamName as string
type StreamNamer interface {
	GetStreamName(aggregate Aggregate) (string, error)
}

// DelegateStreamNamer stores delegates per aggregate type allowing fine grained
// control of stream names for event streams.
type DelegateStreamNamer struct {
	delegates map[string]func(aggregate Aggregate) string
}

// NewDelegateStreamNamer constructs a delegate stream namer
func NewDelegateStreamNamer() *DelegateStreamNamer {
	return &DelegateStreamNamer{
		delegates: make(map[string]func(aggregate Aggregate) string),
	}
}

// RegisterDelegate allows registration of a stream name delegate function for
// the aggregates specified in the variadic aggregates argument.
func (r *DelegateStreamNamer) RegisterDelegate(delegate func(aggregate Aggregate) string, aggregates ...Aggregate) error {
	for _, aggregate := range aggregates {
		t := typeOf(aggregate)
		if _, ok := r.delegates[t]; ok {
			return fmt.Errorf("The stream name delegate for \"%s\" is already registered with the stream namer.", typeOf(aggregate))
		}
		r.delegates[t] = delegate
	}
	return nil
}

// GetStreamName gets the result of the stream name delgate registered for the aggregate type.
func (r *DelegateStreamNamer) GetStreamName(aggregate Aggregate) (string, error) {
	if f, ok := r.delegates[typeOf(aggregate)]; ok {
		return f(aggregate), nil
	}
	return "", fmt.Errorf("There is no stream name delegate for aggregate of type \"%s\"", typeOf(aggregate))
}
