package main

import (
	"os"
	"os/signal"
	"service-worker-sqs-s3-postgres/config/cmd/builder"
	"service-worker-sqs-s3-postgres/core/domain"
	cfiledata "service-worker-sqs-s3-postgres/core/usecases/filedata"
	cmetadata "service-worker-sqs-s3-postgres/core/usecases/metadata"
	rfiledata "service-worker-sqs-s3-postgres/dataproviders/postgres/repository/filedata"
	rmetadata "service-worker-sqs-s3-postgres/dataproviders/postgres/repository/metadata"
	"service-worker-sqs-s3-postgres/dataproviders/server"
	hfiledata "service-worker-sqs-s3-postgres/entrypoints/controllers/filedata"
	hmetadata "service-worker-sqs-s3-postgres/entrypoints/controllers/metadata"
	"syscall"
)

func main() {

	// logger is initialized
	logger := builder.NewLogger()
	logger.Info("Starting service-worker-sqs-s3-postgres ...")
	defer builder.Sync(logger)

	// config is initialized
	config, err := builder.LoadConfig()
	if err != nil {
		logger.Fatalf("error in LoadConfig : %v", err)
	}

	// session aws s3 is initialized
	sessionS3, err := builder.NewSession(config, domain.S3)
	if err != nil {
		logger.Fatalf("error in Session : %v", err)
	}

	// session aws sqs is initialized
	sessionSQS, err := builder.NewSession(config, domain.SQS)
	if err != nil {
		logger.Fatalf("error in Session : %v", err)
	}

	// db is initialized
	db, err := builder.NewDB(config)
	if err != nil {
		logger.Fatalf("error in RDS : %v", err)
	}

	// repositories are initialized
	filedataRepository := rfiledata.NewFileDataRepository(db)
	metadataRepository := rmetadata.NewMetaDataRepository(db)

	// use-cases are initialized
	filedataUseCases := cfiledata.NewFileDataUseCases(filedataRepository)
	metadataUseCases := cmetadata.NewMetaDataUseCases(metadataRepository)

	// controllers are initialized
	filedataController := hfiledata.NewFileDataController(filedataUseCases)
	metadataController := hmetadata.NewMetaDataController(metadataUseCases)

	// consumer is initialized
	sqs, err := builder.NewConsumer(logger, config, sessionSQS, sessionS3, filedataRepository, metadataRepository)
	if err != nil {
		logger.Fatalf("error in SQS : %v", err)
	}

	// processor is initialized
	processor, err := builder.NewProcessor(logger, sqs)
	if err != nil {
		logger.Fatalf("error in Processor : %v", err)
	}
	go processor.Start()

	// server is initialized
	srv := server.NewServer(config.Port, filedataController, metadataController)
	if err = srv.Start(); err != nil {
		logger.Fatalf("error Starting Server: %v", err)
	}

	// Graceful shutdown
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	sig := <-sigQuit

	logger.Infof("Shutting down server with signal [%s] ...", sig.String())
	if err = sqs.Close(); err != nil {
		logger.Error("error Closing Consumer SQS: %v", err)
	}

	if err = srv.Stop(); err != nil {
		logger.Error("error Stopping Server: %v", err)
	}

	logger.Info("service-worker-sqs-s3-postgres ended")

}
