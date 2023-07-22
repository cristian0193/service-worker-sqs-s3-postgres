package mapper

import (
	"service-worker-sqs-s3-postgres/core/domain"
	"service-worker-sqs-s3-postgres/core/domain/entity"
	"time"
)

// ToDomainFileData convert domain filedata to model the postgres filedata .
func ToDomainFileData(f *entity.FileData) *domain.FileData {
	return &domain.FileData{
		ID:      f.ID,
		Message: f.Message,
		Owner:   f.Owner,
		Date:    f.Date,
	}
}

func ToEntityFileData(f *domain.FileData) *entity.FileData {
	return &entity.FileData{
		ID:      f.ID,
		Message: f.Message,
		Owner:   f.Owner,
		Date:    time.Now().Format(time.RFC3339),
	}
}
