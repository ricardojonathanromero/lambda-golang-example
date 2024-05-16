package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardojonathanromero/go-utilities/db/dynamodb"
	"github.com/ricardojonathanromero/go-utilities/environment"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-document-lambda/internal/handler"
	"github.com/ricardojonathanromero/lambda-golang-example/get-document-lambda/pkg/repository"
	"github.com/ricardojonathanromero/lambda-golang-example/get-document-lambda/pkg/service"
	dbInfra "github.com/ricardojonathanromero/lambda-golang-example/internal/db"
)

const (
	logLevelEnv        = "LOG_LEVEL"
	defaultLogLevelEnv = "info"
	envTableName       = "DYNAMODB_TABLE_NAME"
	defaultEmpty       = ""
	appName            = "get-document-lambda"
)

func main() {
	logLevel := environment.GetEnv(logLevelEnv, defaultLogLevelEnv)

	customLog := logger.NewLoggerWithOptions(logger.Opts{
		AppName: appName,
		Level:   logLevel,
	})

	// connect to db
	db := dynamodb.New()
	conn, err := db.Connect()
	if err != nil {
		customLog.Fatalf("error initializing db connection: %s", err.Error())
	}

	defer func() {
		if err = db.Disconnect(); err != nil {
			customLog.Error(err.Error())
		}
	}()

	// configure table
	tableName := environment.GetEnv(envTableName, defaultEmpty)
	err = dbInfra.New(conn, customLog).ConfigureTable(tableName)
	if err != nil {
		customLog.Fatalf("error configuring table: %v", err)
	}

	// init dependency injection
	repo := repository.New(conn, tableName, customLog)
	srv := service.New(repo, customLog)

	lambda.Start(handler.New(srv, customLog))
}
