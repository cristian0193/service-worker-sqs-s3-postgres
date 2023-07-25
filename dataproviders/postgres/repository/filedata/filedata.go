package repository

import (
	"gorm.io/gorm/clause"
	"service-worker-sqs-s3-postgres/core/domain"
	"service-worker-sqs-s3-postgres/core/domain/entity"
	"service-worker-sqs-s3-postgres/core/domain/exceptions"
	"service-worker-sqs-s3-postgres/dataproviders/mapper"
	"service-worker-sqs-s3-postgres/dataproviders/postgres"
)

type IFileDataRepository interface {
	GetID(ID string) (*domain.FileData, error)
	Insert(filedata []*domain.FileData) error
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
func (er *FileDataRepository) GetID(ID string) (*domain.FileData, error) {
	filedata := &entity.FileData{}

	err := er.db.DB.Model(&filedata).Where("id = ?", ID).Scan(&filedata).Error
	if err != nil {
		return nil, exceptions.ErrInternalError
	}

	return mapper.ToDomainFileData(filedata), nil
}

// Insert records a filedata in the database.
func (er *FileDataRepository) Insert(filedata []*domain.FileData) error {

	files := make([]*entity.FileData, 0)

	for _, v := range filedata {
		files = append(files, mapper.ToEntityFileData(v))
	}

	r := er.db.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&files)
	if r.Error != nil {
		r.Rollback()
		return r.Error
	}
	return nil
}
