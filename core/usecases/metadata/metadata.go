package filedata

import (
	"service-worker-sqs-s3-postgres/core/domain"
	repository "service-worker-sqs-s3-postgres/dataproviders/postgres/repository/metadata"
)

type IMetaDataCaseUses interface {
	GetID(trackID string) (*domain.MetaData, error)
}

// MetaDataCaseUses encapsulates all the data necessary for the implementation of the MetaDataRepository.
type MetaDataCaseUses struct {
	metadataRepository repository.IMetaDataRepository
}

// NewMetaDataUseCases instance the repository usecases.
func NewMetaDataUseCases(md repository.IMetaDataRepository) *MetaDataCaseUses {
	return &MetaDataCaseUses{
		metadataRepository: md,
	}
}

// GetID return the metadata by ID.
func (md *MetaDataCaseUses) GetID(trackID string) (*domain.MetaData, error) {
	return md.metadataRepository.GetID(trackID)
}
