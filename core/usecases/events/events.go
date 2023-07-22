package events

import (
	"service-worker-sqs-s3-postgres/core/domain/entity"
	repository "service-worker-sqs-s3-postgres/dataproviders/postgres/repository/events"
)

type IEventCaseUses interface {
	GetID(ID string) (*entity.Events, error)
}

// EventCaseUses encapsulates all the data necessary for the implementation of the EventsRepository.
type EventCaseUses struct {
	eventRepository repository.IEventRepository
}

// NewEventUseCases instance the repository usecases.
func NewEventUseCases(er repository.IEventRepository) *EventCaseUses {
	return &EventCaseUses{
		eventRepository: er,
	}
}

// GetID return the event by ID.
func (es *EventCaseUses) GetID(ID string) (*entity.Events, error) {
	return es.eventRepository.GetID(ID)
}
