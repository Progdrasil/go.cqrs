// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package ycq

type Repository interface {
	Load(aggregateType, id string) (Aggregate, error)
	Save(aggregate Aggregate, expectedVersion *int) error
}