package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardojonathanromero/go-utilities/db/dynamodb"
	"github.com/ricardojonathanromero/go-utilities/environment"
	"github.com/ricardojonathanromero/go-utilities/logger"
	dbInfra "github.com/ricardojonathanromero/lambda-golang-example/business/db"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/internal/handlers"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/repository"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/services"
)

const (
	envName      = "LOG_LEVEL"
	envTableName = "DYNAMODB_TABLE_NAME"
	defaultEnv   = "debug"
	defaultEmpty = ""
	appName      = "get-all-documents-lambda"
)

func main() {
	customLog := logger.NewLoggerWithOptions(logger.Opts{
		AppName: appName,
		Level:   environment.GetEnv(envName, defaultEnv),
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
	srv := services.New(repo, customLog)

	lambda.Start(handlers.New(srv, customLog).HandleRequest)
}
