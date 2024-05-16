package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/models"
)

type Repository interface {
	FindDocumentById(ctx context.Context, id string) (*models.UserDB, error)
}

type repoImpl struct {
	conn      *dynamodb.Client
	log       logger.Logger
	tableName string
}

func New(conn *dynamodb.Client, tableName string, log logger.Logger) Repository {
	return &repoImpl{
		conn:      conn,
		log:       log,
		tableName: tableName,
	}
}

func (repo *repoImpl) FindDocumentById(ctx context.Context, id string) (*models.UserDB, error) {
	var result *models.UserDB
	repo.log.Debugf("processing FindDocumentById: %s", id)

	repo.log.Debug("creating request")
	request := &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"Id": &types.AttributeValueMemberS{Value: id}},
		TableName: aws.String(repo.tableName),
	}

	repo.log.Debug("retrieving item")
	out, err := repo.conn.GetItem(ctx, request)
	if err != nil {
		repo.log.Errorf("error GetItem: %s", err)
		return result, err
	}

	repo.log.Debug("processing result from db")
	err = attributevalue.UnmarshalMap(out.Item, &result)
	if err != nil {
		repo.log.Errorf("error unmarshal response into model: %s", err)
		return result, err
	}

	repo.log.Info("result serialized")
	return result, nil
}
