package filedata

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"service-worker-sqs-s3-postgres/core/domain/exceptions"
	cases "service-worker-sqs-s3-postgres/core/usecases/filedata"
	env "service-worker-sqs-s3-postgres/dataproviders/utils"
)

// FileDataController encapsulates all the data necessary for the implementation of the FileDataService.
type FileDataController struct {
	filedataUseCases cases.IFileDataCaseUses
}

// NewFileDataController instantiate a new filedata controller.
func NewFileDataController(es cases.IFileDataCaseUses) *FileDataController {
	return &FileDataController{
		filedataUseCases: es,
	}
}

// GetID return a filedata by ID [filedataUseCases.GetID].
func (ec *FileDataController) GetID(c echo.Context) error {
	ID, err := env.GetParam(c, "id")
	if err != nil {
		return exceptions.NewError(http.StatusBadRequest, err)
	}
	filedata, err := ec.filedataUseCases.GetID(ID)
	if err != nil {
		return exceptions.HandleServiceError(err)
	}
	return c.JSON(http.StatusOK, filedata)
}
