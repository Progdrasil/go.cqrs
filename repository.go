package cqrs

type Repository interface {
	Load(aggregate Aggregate, aggregateId string) (Aggregate, error)
	Save(aggregate Aggregate, expectedVersion int) error
}