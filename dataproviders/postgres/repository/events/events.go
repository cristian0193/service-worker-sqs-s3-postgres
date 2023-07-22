package repository

import (
	"gorm.io/gorm/clause"
	"service-worker-sqs-s3-postgres/core/domain/entity"
	"service-worker-sqs-s3-postgres/core/domain/exceptions"
	"service-worker-sqs-s3-postgres/dataproviders/postgres"
)

type IEventRepository interface {
	GetID(ID string) (*entity.Events, error)
	Insert(events *entity.Events) error
}

// EventRepository encapsulates all the data needed to the persistence in the event table.
type EventRepository struct {
	db *postgres.ClientDB
}

// NewEventRepository instance the connection to the postgres.
func NewEventRepository(db *postgres.ClientDB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

// GetID return the event by ID.
func (er *EventRepository) GetID(ID string) (*entity.Events, error) {
	event := &entity.Events{}

	err := er.db.DB.Model(&event).Where("id = ?", ID).Scan(&event).Error
	if err != nil {
		return nil, exceptions.ErrInternalError
	}

	return event, nil
}

// Insert records an event in the database.
func (er *EventRepository) Insert(events *entity.Events) error {
	r := er.db.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&events)
	if r.Error != nil {
		r.Rollback()
		return r.Error
	}
	return nil
}
