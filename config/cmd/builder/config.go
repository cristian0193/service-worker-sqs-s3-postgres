package builder

import (
	env "service-worker-sqs-s3-postgres/dataproviders/utils"
)

// Configuration represents parameters of application.
type Configuration struct {
	Port                 int
	ApplicationID        string
	LogLevel             string
	Region               string
	AccessKey            string
	SecretKey            string
	SQSUrl               string
	SQSMaxMessages       int
	SQSVisibilityTimeout int
	S3Bucket             string
	DBPort               string
	DBHost               string
	DBName               string
	DBUsername           string
	DBPassword           string
}

// LoadConfig get all the configuration variables for the implemented usecases.
func LoadConfig() (*Configuration, error) {
	applicationID, err := env.GetString("APPLICATION_ID")
	if err != nil {
		return nil, err
	}

	port, err := env.GetInt("SERVER_PORT")
	if err != nil {
		return nil, err
	}

	loglevel, err := env.GetString("LOG_LEVEL")
	if err != nil {
		return nil, err
	}

	access, err := env.GetString("AWS_ACCESS_KEY")
	if err != nil {
		return nil, err
	}

	secret, err := env.GetString("AWS_SECRET_KEY")
	if err != nil {
		return nil, err
	}

	region, err := env.GetString("AWS_REGION")
	if err != nil {
		return nil, err
	}

	sqsUrl, err := env.GetString("AWS_SQS_URL")
	if err != nil {
		return nil, err
	}

	sqsMaxMessages, err := env.GetInt("AWS_SQS_MAX_MESSAGES")
	if err != nil {
		return nil, err
	}

	sqsVisibilityTimeout, err := env.GetInt("AWS_SQS_VISIBILITY_TIMEOUT")
	if err != nil {
		return nil, err
	}

	s3Bucket, err := env.GetString("AWS_S3_BUCKET")
	if err != nil {
		return nil, err
	}

	dbPort, err := env.GetString("DB_PORT")
	if err != nil {
		return nil, err
	}

	dbHost, err := env.GetString("DB_HOST")
	if err != nil {
		return nil, err
	}

	dbName, err := env.GetString("DB_NAME")
	if err != nil {
		return nil, err
	}

	dbUsername, err := env.GetString("DB_USERNAME")
	if err != nil {
		return nil, err
	}

	dbPassword, err := env.GetString("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	return &Configuration{
		Port:                 port,
		ApplicationID:        applicationID,
		LogLevel:             loglevel,
		AccessKey:            access,
		SecretKey:            secret,
		Region:               region,
		SQSUrl:               sqsUrl,
		SQSMaxMessages:       sqsMaxMessages,
		SQSVisibilityTimeout: sqsVisibilityTimeout,
		S3Bucket:             s3Bucket,
		DBPort:               dbPort,
		DBHost:               dbHost,
		DBName:               dbName,
		DBUsername:           dbUsername,
		DBPassword:           dbPassword,
	}, nil
}
