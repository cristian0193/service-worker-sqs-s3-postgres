package csvreader

import (
	"encoding/csv"
	"errors"
	"io"
	"service-worker-sqs-s3-postgres/core/domain"
	"service-worker-sqs-s3-postgres/dataproviders/utils"
	"strings"

	"github.com/spf13/afero"
	"go.uber.org/zap"
)

const numberColumns = 3

// Read returns a list of filedata to be persisted in the database .
func Read(fileName string, logger *zap.SugaredLogger) ([]domain.FileData, error) {
	filedata := make([]domain.FileData, 0)

	a := afero.Afero{
		Fs: afero.NewOsFs(),
	}

	contents, err := a.ReadFile(fileName)
	if err != nil {
		return filedata, err
	}

	for rec := range ProcessCSV(strings.NewReader(string(contents)), logger) {
		if len(rec) == numberColumns {
			filedata = columnsToFileData(rec, filedata)
		}
	}
	return filedata, nil
}

// ProcessCSV returns a channel for reading data associated with the csv .
func ProcessCSV(rc io.Reader, logger *zap.SugaredLogger) (ch chan []string) {
	ch = make(chan []string)

	go func() {
		r := csv.NewReader(rc)
		r.LazyQuotes = true

		if _, err := r.Read(); err != nil { // read header
			logger.Fatal(err)
		}
		defer close(ch)
		for {
			rec, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				logger.Fatal(err)
			}
			ch <- rec
		}
	}()
	return
}

// ---------- Helpers ------------ //

func columnsToFileData(rec []string, info []domain.FileData) []domain.FileData {
	filedata := domain.FileData{
		ID:      utils.StringToInt64(rec[0]),
		Message: rec[1],
		Owner:   rec[2],
	}
	return append(info, filedata)
}
