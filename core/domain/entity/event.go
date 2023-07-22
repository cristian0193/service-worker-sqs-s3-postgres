package entity

import (
	"service-worker-sqs-s3-postgres/core/domain"
)

// Events represents the entity.
type Events struct {
	ID      string `gorm:"NULL;TYPE:VARCHAR(200);COLUMN:id" json:"id"`
	Message string `gorm:"NULL;TYPE:VARCHAR(200);COLUMN:message" json:"message"`
	Date    string `gorm:"NULL;TYPE:VARCHAR(200);COLUMN:date" json:"date"`
}

// TableName definition name for table .
func (Events) TableName() string {
	return "events"
}

// ToDomainEvents convert domain event to model the postgres events .
func (e Events) ToDomainEvents() *domain.Events {
	return &domain.Events{
		ID:      e.ID,
		Message: e.Message,
		Date:    e.Date,
	}
}
