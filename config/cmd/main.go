package main

import (
	"os"
	"os/signal"
	"service-worker-sqs-s3-postgres/config/cmd/builder"
	"service-worker-sqs-s3-postgres/core/domain"
	cases "service-worker-sqs-s3-postgres/core/usecases/events"
	repository "service-worker-sqs-s3-postgres/dataproviders/postgres/repository/events"
	"service-worker-sqs-s3-postgres/dataproviders/server"
	"service-worker-sqs-s3-postgres/entrypoints/controllers/events"
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
	eventRepository := repository.NewEventRepository(db)

	// usecases are initialized
	eventUseCases := cases.NewEventUseCases(eventRepository)

	// controllers are initialized
	eventController := events.NewEventController(eventUseCases)

	// consumer is initialized
	sqs, err := builder.NewConsumer(logger, config, sessionSQS, eventRepository)
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
	srv := server.NewServer(config.Port, eventController)
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
