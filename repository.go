// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cqrs

type Repository interface {
	Load(aggregateType AggregateType, aggregateId AggregateId) (Aggregate, error)
	Save(aggregate Aggregate, expectedVersion *AggregateVersion) error
}