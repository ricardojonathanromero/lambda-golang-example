package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/business/models"
)

type Repository interface {
	FindAllDocuments(ctx context.Context) ([]*models.UserDB, error)
}

type repositoryImpl struct {
	conn      *dynamodb.Client
	log       logger.Logger
	tableName string
}

func New(conn *dynamodb.Client, tableName string, log logger.Logger) Repository {
	return &repositoryImpl{
		conn:      conn,
		log:       log,
		tableName: tableName,
	}
}

func (repo *repositoryImpl) FindAllDocuments(ctx context.Context) ([]*models.UserDB, error) {
	// scan input
	input := &dynamodb.ScanInput{TableName: aws.String(repo.tableName)}

	repo.log.Debug("executing scan")
	output, err := repo.conn.Scan(ctx, input)
	if err != nil {
		// eval error
		repo.log.Errorf("error executing scan: %s", err)
		return nil, err
	}

	repo.log.Debug("scan response received, serializing response ...")
	var users []*models.UserDB
	err = attributevalue.UnmarshalListOfMaps(output.Items, &users)
	if err != nil {
		repo.log.Errorf("error serializing reponse into model: %s", err)
		return nil, err
	}

	repo.log.Debugf("response serialized - total items: %d", len(users))
	return users, nil
}
