// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cqrs

import "github.com/jetbasrawi/go.cqrs/internal/uuid"

type AggregateType struct {
	name string
}

func NewAggregateType(i interface{}) AggregateType {
	t := typeOf(i)
	return AggregateType{name: t}
}

func (t *AggregateType) AggregateTypeName() string {
	return t.name
}

type AggregateId struct {
	stringId string
}

func NewAggregateId(id string) AggregateId {
	if id == "" {
		return AggregateId{id}
	} else {
		return AggregateId{uuid.NewV4().String()}
	}
}

func (id *AggregateId) AggregateIdString() string {
	return id.stringId
}

type AggregateVersion struct {
	number int
}

func NewAggregateVersion(i int) *AggregateVersion {
	return &AggregateVersion{number: i}
}

func (v *AggregateVersion) AggregateVersionNumber() int {
	return v.number
}

//Aggregate is the interface that all aggregates should implement
type Aggregate interface {
	AggregateID() AggregateId
	OriginalVersion() AggregateVersion
	CurrentVersion() AggregateVersion
	IncrementVersion()
	Apply(events EventMessage, isNew bool)
	TrackChange(EventMessage)
	GetChanges() []EventMessage
	ClearChanges()
}

// AggregateBase is a type that can be embedded in an Aggregate
// implementation to handle common aggragate behaviour
//
// All required methods to implement an aggregate are here, to implement the
// Aggregate root interface your aggregate will need to implement the Apply
// method that will contain behaviour specific to your aggregate.
type AggregateBase struct {
	id      AggregateId
	version AggregateVersion
	changes []EventMessage
}

// NewAggregateBase contructs a new AggregateBase.
func NewAggregateBase(id AggregateId) *AggregateBase {
	return &AggregateBase{
		id:      id,
		changes: []EventMessage{},
		version: AggregateVersion{-1},
	}
}

// StreamName returns the StreamName
func (a *AggregateBase) AggregateID() AggregateId {
	return a.id
}

// OriginalVersion returns the version of the aggregate as it was when it was
// instantiated or loaded from the repository.
//
// Importantly an aggregate with one event applied will be at version 0
// this allows the aggregates to match the version in the domain where
// the first event will be version 0.
func (a *AggregateBase) OriginalVersion() AggregateVersion {
	return a.version
}

// CurrentVersion returns the version of the aggregate as it was when it was
// instantiated or loaded from the repository.
//
// Importantly an aggregate with one event applied will be at version 0
// this allows the aggregates to match the version in the domain where
// the first event will be version 0.
func (a *AggregateBase) CurrentVersion() AggregateVersion {
	return AggregateVersion{a.version.AggregateVersionNumber() + len(a.changes)}
}

// IncrementVersion increments the aggregate version number by one.
func (a *AggregateBase) IncrementVersion() {
	a.version.number++
}

// TrackChange stores the EventMessage in the changes collection.
//
// Changes are new, unpersisted events that have been applied to the aggregate.
func (a *AggregateBase) TrackChange(event EventMessage) {
	a.changes = append(a.changes, event)
}

// GetChanges returns the collection of new unpersisted events that have
// been applied to the aggregate.
func (a *AggregateBase) GetChanges() []EventMessage {
	return a.changes
}

//ClearChanges removes all unpersisted events from the aggregate.
func (a *AggregateBase) ClearChanges() {
	a.changes = []EventMessage{}
}
