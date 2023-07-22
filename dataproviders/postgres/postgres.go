package postgres

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "gorm.io/gorm/logger"
	"service-worker-sqs-s3-postgres/core/domain/entity"
	"time"
)

// ClientDB represents DB client.
type ClientDB struct {
	DB     *gorm.DB
	params Params
}

type Params struct {
	host     string
	userName string
	password string
	name     string
	port     string
}

// NewDBClient instances of a Client to connect postgresql with parameters.
func NewDBClient(host, username, password, name, port string) *ClientDB {
	return &ClientDB{
		params: Params{
			host:     host,
			userName: username,
			password: password,
			name:     name,
			port:     port,
		},
	}
}

// Open the postgres connection only the first time. The next times, it maintains the same connection.
func (client *ClientDB) Open() error {

	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		client.params.host,
		client.params.userName,
		client.params.password,
		client.params.name,
		client.params.port)

	if client.DB == nil {
		db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 logger.Default.LogMode(logger.Silent),
			CreateBatchSize:        1000,
		})
		if err != nil {
			return errors.Wrapf(err, "Error opening postgres file: %v", err.Error())
		}

		dbs := db.Session(&gorm.Session{CreateBatchSize: 1000})
		sqlDB, err := dbs.DB()
		if err != nil {
			return errors.Wrapf(err, "Error instance postgres : %v", err.Error())
		}

		sqlDB.SetConnMaxLifetime(5 * time.Minute)
		sqlDB.SetConnMaxIdleTime(10)
		sqlDB.SetMaxOpenConns(10)

		err = dbs.AutoMigrate(entity.Events{})
		if err != nil {
			return errors.Wrapf(err, "Error migrating postgres : %v", err.Error())
		}

		client.DB = dbs
	}

	return nil
}
