package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardojonathanromero/go-utilities/db/dynamodb"
	"github.com/ricardojonathanromero/go-utilities/environment"
	"github.com/ricardojonathanromero/go-utilities/logger"
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
	log := logger.NewLoggerWithOptions(logger.Opts{
		AppName: appName,
		Level:   environment.GetEnv(envName, defaultEnv),
	})

	// connect to db
	db := dynamodb.New()
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("error initializing db connection: %s", err.Error())
	}

	defer func() {
		if err = db.Disconnect(); err != nil {
			log.Error(err.Error())
		}
	}()

	tableName := environment.GetEnv(envTableName, defaultEmpty)

	// init dependency injection
	repo := repository.New(conn, tableName, log)
	srv := services.New(repo, log)

	lambda.Start(handlers.New(srv, log).HandleRequest)
}
