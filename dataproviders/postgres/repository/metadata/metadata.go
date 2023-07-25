package repository

import (
	"service-worker-sqs-s3-postgres/core/domain"
	"service-worker-sqs-s3-postgres/core/domain/entity"
	"service-worker-sqs-s3-postgres/core/domain/exceptions"
	"service-worker-sqs-s3-postgres/dataproviders/mapper"
	"service-worker-sqs-s3-postgres/dataproviders/postgres"
)

type IMetaDataRepository interface {
	GetID(trackID string) (*domain.MetaData, error)
	Insert(metadata domain.MetaData) error
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
func (er *MetaDataRepository) GetID(trackID string) (*domain.MetaData, error) {
	metadata := &entity.MetaData{}

	err := er.db.DB.Model(&metadata).Where("trackid = ?", trackID).Scan(&metadata).Error
	if err != nil {
		return nil, exceptions.ErrInternalError
	}

	return mapper.ToDomainMetaData(metadata), nil
}

// Insert records a metadata in the database.
func (er *MetaDataRepository) Insert(metadata *domain.MetaData) error {

	meta := mapper.ToEntityMetaData(metadata)

	err := er.db.DB.Create(&meta).Error
	if err != nil {
		er.db.DB.Rollback()
		return err
	}
	return nil
}
