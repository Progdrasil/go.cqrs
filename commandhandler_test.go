// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cqrs

import (
	. "gopkg.in/check.v1"
)

var _ = Suite(&CommandHandlerSuite{})

type CommandHandlerSuite struct {
}

type TestCommandHandler struct {
	command CommandMessage
}

func (t *TestCommandHandler) Handle(command CommandMessage) error {
	t.command = command
	return nil
}

type MockRepository struct {
	aggregates map[AggregateId]Aggregate
}

func (m *MockRepository) Load(aggregateType string, id AggregateId) (Aggregate, error) {
	return m.aggregates[id], nil
}

func (m *MockRepository) Save(aggregate Aggregate) error {
	m.aggregates[*aggregate.AggregateId()] = aggregate
	return nil
}
