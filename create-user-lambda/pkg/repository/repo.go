package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardojonathanromero/go-utilities/logger"
)

type Repository interface {
	InsertUser(ctx context.Context, user any) error
}

type repoImpl struct {
	tableName string
	client    *dynamodb.Client
	log       logger.Logger
}

func New(tableName string, client *dynamodb.Client, log logger.Logger) Repository {
	return &repoImpl{
		tableName: tableName,
		client:    client,
		log:       log,
	}
}

func (repo *repoImpl) InsertUser(ctx context.Context, user any) error {
	repo.log.Debug("marshalling input")
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		repo.log.Errorf("error marshalling input: %v", err)
		return err
	}

	repo.log.Debug("sending input")
	req := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(repo.tableName),
		ConditionExpression: aws.String("attribute_not_exists(Id)"),
	}

	_, err = repo.client.PutItem(ctx, req)
	if err != nil {
		repo.log.Errorf("error put item: %v", err)
		return err
	}

	repo.log.Debug("item inserted")
	return nil
}
