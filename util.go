package cqrs

import (
	"reflect"

	"github.com/jetbasrawi/go.cqrs/internal/uuid"
)

// typeOf is a convenience function that returns the name of a type
func typeOf(i interface{}) string {
	return reflect.TypeOf(i).Elem().Name()
}

// NewUUID returns a new v4 uuid as a string
func NewUUID() string {
	return uuid.NewUUID()
}

// Int returns a pointer to int.
//
// There are a number of places where a pointer to int
// is required such as expectedVersion argument on the repository
// and this helper function makes keeps the code cleaner in these
// cases.
//func Int(i int) *int {
//	return &i
//}
