package repository

import (
	"gorm.io/gorm/clause"
	"service-worker-sqs-s3-postgres/core/domain/entity"
	"service-worker-sqs-s3-postgres/core/domain/exceptions"
	"service-worker-sqs-s3-postgres/dataproviders/postgres"
)

type IFileDataRepository interface {
	GetID(ID string) (*entity.FileData, error)
	Insert(filedata []entity.FileData) error
}

// FileDataRepository encapsulates all the data needed to the persistence in the filedata table.
type FileDataRepository struct {
	db *postgres.ClientDB
}

// NewFileDataRepository instance the connection to the postgres.
func NewFileDataRepository(db *postgres.ClientDB) *FileDataRepository {
	return &FileDataRepository{
		db: db,
	}
}

// GetID return the filedata by ID.
func (er *FileDataRepository) GetID(ID string) (*entity.FileData, error) {
	filedata := &entity.FileData{}

	err := er.db.DB.Model(&filedata).Where("id = ?", ID).Scan(&filedata).Error
	if err != nil {
		return nil, exceptions.ErrInternalError
	}

	return filedata, nil
}

// Insert records a filedata in the database.
func (er *FileDataRepository) Insert(filedata []entity.FileData) error {
	r := er.db.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&filedata)
	if r.Error != nil {
		r.Rollback()
		return r.Error
	}
	return nil
}
