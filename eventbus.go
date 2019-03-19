// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cqrs

// EventBus is the interface that an event bus must implement.
type EventBus interface {
	PublishEvent(EventMessage) error
	AddEventHandler(EventHandler, ...interface{})
}

// InternalEventBus provides a lightweight in process event bus
type InternalEventBus struct {
	eventHandlers map[string]map[EventHandler]struct{}
}

// NewInternalEventBus constructs a new InternalEventBus
func NewInternalEventBus() *InternalEventBus {
	b := &InternalEventBus{
		eventHandlers: make(map[string]map[EventHandler]struct{}),
	}
	return b
}

// PublishEvent publishes events to all registered event handlers
func (t *InternalEventBus) PublishEvent(event EventMessage) error {
	if handlers, ok := t.eventHandlers[event.EventType()]; ok {
		for handler := range handlers {
			handler.Handle(event)
		}
	} else {
		return &ErrNoConfiguredHandler{targetType: event.EventType(), handler: typeOf(t)}
	}
	return nil
}

// AddCommandHandler registers an event handler for all of the events specified in the
// variadic events parameter.
func (t *InternalEventBus) AddEventHandler(handler EventHandler, events ...interface{}) {
	for _, event := range events {
		et := typeOf(event)

		// There can be multiple handlers for any event.
		// Here we check that a map is initialized to hold these handlers
		// for a given type. If not we create one.
		if _, ok := t.eventHandlers[et]; !ok {
			t.eventHandlers[et] = make(map[EventHandler]struct{})
		}

		// Add this handler to the collection of handlers for the type.
		t.eventHandlers[et][handler] = struct{}{}
	}
}
