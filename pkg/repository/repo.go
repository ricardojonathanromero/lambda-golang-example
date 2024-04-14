package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	db "github.com/ricardojonathanromero/go-utilities/db/dynamodb"
	"github.com/ricardojonathanromero/go-utilities/logger"
)

type Repository interface {
}

type repositoryImpl struct {
	client    db.DB
	log       logger.Logger
	tableName string
}

func New(client db.DB, tableName string, log logger.Logger) Repository {
	return &repositoryImpl{
		client:    client,
		log:       log,
		tableName: tableName,
	}
}

func (repo *repositoryImpl) FindAllDocuments(ctx context.Context) {
	conn, err := repo.conn()
	if err != nil {
		// connection not established
		// do something
	}

	// query
	conn.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 aws.String(repo.tableName),
		AttributesToGet:           nil,
		ConditionalOperator:       "",
		ConsistentRead:            nil,
		ExclusiveStartKey:         nil,
		ExpressionAttributeNames:  nil,
		ExpressionAttributeValues: nil,
		FilterExpression:          nil,
		IndexName:                 nil,
		Limit:                     nil,
		ProjectionExpression:      nil,
		ReturnConsumedCapacity:    "",
		ScanFilter:                nil,
		Segment:                   nil,
		Select:                    "",
		TotalSegments:             nil,
	})
}

func (repo *repositoryImpl) conn() (*dynamodb.Client, error) {
	conn, err := repo.client.Connect()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
