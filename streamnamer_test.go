// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cqrs

import (
	"fmt"

	. "gopkg.in/check.v1"
)

var _ = Suite(&DelegateStreamNamerSuite{})

type DelegateStreamNamerSuite struct {
	namer *DelegateStreamNamer
}

func (s *DelegateStreamNamerSuite) SetUpTest(c *C) {
	s.namer = NewDelegateStreamNamer()
}

func (s *DelegateStreamNamerSuite) TestNewDelegateStreamNamer(c *C) {
	namer := NewDelegateStreamNamer()
	c.Assert(namer.delegates, NotNil)
}

func (s *DelegateStreamNamerSuite) TestCanRegisterStreamNameDelegate(c *C) {

	err := s.namer.RegisterDelegate(func(a Aggregate) string { return typeOf(a) + a.AggregateID().String() },
		&SomeAggregate{},
	)
	c.Assert(err, IsNil)
	agg := NewSomeAggregate(NewAggregateId())
	f, _ := s.namer.delegates[typeOf(agg)]
	stream := f(agg)
	c.Assert(stream, Equals, typeOf(agg)+agg.AggregateID().String())
}

func (s *DelegateStreamNamerSuite) TestCanRegisterStreamNameDelegateWithMultipleAggregateRoots(c *C) {
	err := s.namer.RegisterDelegate(func(a Aggregate) string { return typeOf(a) + a.AggregateID().String() },
		&SomeAggregate{},
		&SomeOtherAggregate{},
	)
	c.Assert(err, IsNil)

	ar1 := NewSomeAggregate(NewAggregateId())
	f, _ := s.namer.delegates[typeOf(ar1)]
	stream := f(ar1)
	c.Assert(stream, Equals, typeOf(ar1)+ar1.AggregateID().String())

	ar2 := NewSomeOtherAggregate(NewAggregateId())
	f2, _ := s.namer.delegates[typeOf(ar2)]
	stream2 := f2(ar2)
	c.Assert(stream2, Equals, typeOf(ar2)+ar2.AggregateID().String())
}

func (s *DelegateStreamNamerSuite) TestCanGetStreamName(c *C) {
	err := s.namer.RegisterDelegate(func(a Aggregate) string { return typeOf(a) + a.AggregateID().String() },
		&SomeAggregate{},
	)
	c.Assert(err, IsNil)
	agg := NewSomeAggregate(NewAggregateId())
	stream, err := s.namer.GetStreamName(agg)
	c.Assert(err, IsNil)
	c.Assert(stream, Equals, typeOf(agg)+agg.AggregateID().String())
}

func (s *DelegateStreamNamerSuite) TestGetStreamNameReturnsAnErrorIfNoDelegateRegisteredForAggregate(c *C) {
	err := s.namer.RegisterDelegate(func(a Aggregate) string { return typeOf(a) + a.AggregateID().String() },
		&SomeAggregate{},
	)
	agg := NewSomeOtherAggregate(NewAggregateId())
	stream, err := s.namer.GetStreamName(agg)
	c.Assert(err, NotNil)
	c.Assert(stream, Equals, "")
	c.Assert(err, DeepEquals, fmt.Errorf("There is no stream name delegate for aggregate of type \"%s\"",
		typeOf(agg)))
}

func (s *DelegateStreamNamerSuite) TestRegisteringAStreamNameDelegateMoreThanOnceReturnsAndError(c *C) {

	err := s.namer.RegisterDelegate(func(a Aggregate) string { return typeOf(a) + a.AggregateID().String() },
		&SomeAggregate{},
	)
	c.Assert(err, IsNil)

	err = s.namer.RegisterDelegate(func(a Aggregate) string { return a.AggregateID().String() },
		&SomeAggregate{},
	)
	c.Assert(err, DeepEquals,
		fmt.Errorf("The stream name delegate for \"%s\" is already registered with the stream namer.",
			typeOf(NewSomeAggregate(NewAggregateId()))))
}