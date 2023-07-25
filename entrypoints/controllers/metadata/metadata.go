package filedata

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"service-worker-sqs-s3-postgres/core/domain/exceptions"
	cases "service-worker-sqs-s3-postgres/core/usecases/metadata"
	env "service-worker-sqs-s3-postgres/dataproviders/utils"
)

// MetaDataController encapsulates all the data necessary for the implementation of the MetaDataService.
type MetaDataController struct {
	metadataUseCases cases.IMetaDataCaseUses
}

// NewMetaDataController instantiate a new metadata controller.
func NewMetaDataController(es cases.IMetaDataCaseUses) *MetaDataController {
	return &MetaDataController{
		metadataUseCases: es,
	}
}

// GetID return a metadata by track ID [metadataUseCases.GetID].
func (ec *MetaDataController) GetID(c echo.Context) error {
	ID, err := env.GetParam(c, "trackid")
	if err != nil {
		return exceptions.NewError(http.StatusBadRequest, err)
	}
	metadata, err := ec.metadataUseCases.GetID(ID)
	if err != nil {
		return exceptions.HandleServiceError(err)
	}
	return c.JSON(http.StatusOK, metadata)
}
