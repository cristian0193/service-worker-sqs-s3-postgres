package builder

import (
	"service-worker-sqs-s3-postgres/dataproviders/postgres"
)

// NewDB defines all configurations to instantiate a postgres client.
func NewDB(config *Configuration) (*postgres.ClientDB, error) {
	db := postgres.NewDBClient(config.DBHost, config.DBUsername, config.DBPassword, config.DBName, config.DBPort)
	err := db.Open()

	return db, err
}
