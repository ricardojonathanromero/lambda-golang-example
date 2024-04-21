package empty_test

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/utils/tests"
)

const (
	startPort = 8001
	tableName = "users"
)

var (
	port             int
	dynamodbTestConn tests.DBSuite
	conn             *dynamodb.Client
	log              logger.Logger
)
