package filedata

import (
	"service-worker-sqs-s3-postgres/core/domain"
	repository "service-worker-sqs-s3-postgres/dataproviders/postgres/repository/filedata"
)

type IFileDataCaseUses interface {
	GetID(ID string) (*domain.FileData, error)
}

// FileDataCaseUses encapsulates all the data necessary for the implementation of the FileDataRepository.
type FileDataCaseUses struct {
	filedataRepository repository.IFileDataRepository
}

// NewFileDataUseCases instance the repository usecases.
func NewFileDataUseCases(fr repository.IFileDataRepository) *FileDataCaseUses {
	return &FileDataCaseUses{
		filedataRepository: fr,
	}
}

// GetID return the filedata by ID.
func (fd *FileDataCaseUses) GetID(ID string) (*domain.FileData, error) {
	return fd.filedataRepository.GetID(ID)
}
