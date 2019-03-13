package cqrs

type Repository interface {
	Load(aggregateType string, aggregateId string) (Aggregate, error)
	Save(aggregate Aggregate, expectedVersion int) error
}