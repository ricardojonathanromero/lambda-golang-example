package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ricardojonathanromero/go-utilities/db/dynamodb"
	"github.com/ricardojonathanromero/go-utilities/environment"
	"github.com/ricardojonathanromero/go-utilities/logger"
	dbInfra "github.com/ricardojonathanromero/lambda-golang-example/business/db"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/internal/handler"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/repository"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/service"
)

const (
	appName      = "create-user-lambda"
	envTableName = "DYNAMODB_TABLE_NAME"
	envName      = "ENV"
	defaultEnv   = "local"
	defaultEmpty = ""
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
	repo := repository.New(tableName, conn, customLog)
	srv := service.New(repo, customLog)
	lambda.Start(handler.New(srv, customLog).HandleCreateUser)
}
