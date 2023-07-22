package repository

import (
	"service-worker-sqs-s3-postgres/core/domain/entity"
	"service-worker-sqs-s3-postgres/core/domain/exceptions"
	"service-worker-sqs-s3-postgres/dataproviders/postgres"
)

type IMetaDataRepository interface {
	GetID(trackID string) (*entity.MetaData, error)
	Insert(metadata entity.MetaData) error
}

// MetaDataRepository encapsulates all the data needed to the persistence in the filedata table.
type MetaDataRepository struct {
	db *postgres.ClientDB
}

// NewMetaDataRepository instance the connection to the postgres.
func NewMetaDataRepository(db *postgres.ClientDB) *MetaDataRepository {
	return &MetaDataRepository{
		db: db,
	}
}

// GetID return the metadata by ID.
func (er *MetaDataRepository) GetID(trackID string) (*entity.MetaData, error) {
	metadata := &entity.MetaData{}

	err := er.db.DB.Model(&metadata).Where("trackid = ?", trackID).Scan(&metadata).Error
	if err != nil {
		return nil, exceptions.ErrInternalError
	}

	return metadata, nil
}

// Insert records a metadata in the database.
func (er *MetaDataRepository) Insert(metadata entity.MetaData) error {
	err := er.db.DB.Create(&metadata).Error
	if err != nil {
		er.db.DB.Rollback()
		return err
	}
	return nil
}
