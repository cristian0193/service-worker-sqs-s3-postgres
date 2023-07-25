package mapper

import (
	"service-worker-sqs-s3-postgres/core/domain"
	"service-worker-sqs-s3-postgres/core/domain/entity"
)

// ToDomainMetaData convert domain metadata to model the postgres metadata .
func ToDomainMetaData(m *entity.MetaData) *domain.MetaData {
	return &domain.MetaData{
		TrackID:  m.TrackID,
		Bucket:   m.Bucket,
		FileName: m.FileName,
		Key:      m.Key,
		Size:     m.Size,
	}
}

func ToEntityMetaData(f *domain.MetaData) *entity.MetaData {
	return &entity.MetaData{
		TrackID:  f.TrackID,
		Bucket:   f.Bucket,
		FileName: f.FileName,
		Key:      f.Key,
		Size:     f.Size,
	}
}
