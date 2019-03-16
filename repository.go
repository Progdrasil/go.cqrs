package cqrs

type Repository interface {
	Load(aggregate Aggregate) error
	Save(aggregate Aggregate) error
}