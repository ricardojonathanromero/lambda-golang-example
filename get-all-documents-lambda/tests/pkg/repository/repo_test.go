package repository_test

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	dynamodb2 "github.com/ricardojonathanromero/go-utilities/db/dynamodb"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/repository"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/models"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/utils/tests"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestRepositoryImpl_FindAllDocuments(t *testing.T) {
	// before all
	log := logger.NewLoggerWithOptions(logger.Opts{
		AppName: "get-all-documents-lambda-repository-test",
		Level:   "debug",
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	defaultHttpClient := &http.Client{Timeout: time.Second * 10}
	dynamoDBURL := "http://localhost:8000/"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...any) (aws.Endpoint, error) {
				return aws.Endpoint{URL: dynamoDBURL}, nil
			})),
		config.WithHTTPClient(defaultHttpClient),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummyKey", "dummySecret", "")),
	)

	if err != nil {
		t.Errorf("error configuring db conn: %v", err)
		t.FailNow()
	}

	// global connection
	conn := dynamodb.NewFromConfig(cfg)

	t.Run("when_receives_request", func(t *testing.T) {
		t.Run("return_all_documents", func(t *testing.T) {
			result := `{
    "ConsumedCapacity": {
      "CapacityUnits": 0.5,
      "TableName": "my-table"
    },
    "Count": 4,
    "Items": [
  {
    "Age": {
      "N": "33"
    },
    "CreatedAt": {
      "S": "2024-04-14T13:44:37.609166-06:00"
    },
    "Email": {
      "S": "john.smith@test.com"
    },
    "Id": {
      "S": "1"
    },
    "Lastname": {
      "S": "smith"
    },
    "Name": {
      "S": "john"
    },
    "UpdatedAt": {
      "S": "2024-04-14T13:44:37.609169-06:00"
    }
  },
  {
    "Age": {
      "N": "25"
    },
    "CreatedAt": {
      "S": "2024-04-14T13:44:37.609169-06:00"
    },
    "Email": {
      "S": "john.smith2@test.com"
    },
    "Id": {
      "S": "2"
    },
    "Lastname": {
      "S": "smith"
    },
    "Name": {
      "S": "john"
    },
    "UpdatedAt": {
      "S": "2024-04-14T13:44:37.609169-06:00"
    }
  },
  {
    "Age": {
      "N": "17"
    },
    "CreatedAt": {
      "S": "2024-04-14T13:44:37.609169-06:00"
    },
    "Email": {
      "S": "john.smith3@test.com"
    },
    "Id": {
      "S": "3"
    },
    "Lastname": {
      "S": "smith"
    },
    "Name": {
      "S": "john"
    },
    "UpdatedAt": {
      "S": "2024-04-14T13:44:37.609169-06:00"
    }
  },
  {
    "Age": {
      "N": "49"
    },
    "CreatedAt": {
      "S": "2024-04-14T13:44:37.609169-06:00"
    },
    "Email": {
      "S": "john.smith4@test.com"
    },
    "Id": {
      "S": "4"
    },
    "Lastname": {
      "S": "smith"
    },
    "Name": {
      "S": "john"
    },
    "UpdatedAt": {
      "S": "2024-04-14T13:44:37.609169-06:00"
    }
  }
],
    "ScannedCount": 4
  }`

			resp := httpmock.NewStringResponder(http.StatusOK, result)
			httpmock.RegisterResponder(http.MethodPost, dynamoDBURL, resp)

			repo := repository.New(conn, "my-table", log)
			users, err := repo.FindAllDocuments(context.TODO())
			assert.NoError(t, err)
			assert.Len(t, users, 4)
			log.Debug(toString(users))
		})

		t.Run("return_an_error", func(t *testing.T) {
			result := `{
						"Code": "TransactionConflictException",
						"Message": "A conflict occurs trying to scan documents"
					}`

			resp := httpmock.NewStringResponder(http.StatusBadRequest, result)
			httpmock.RegisterResponder(http.MethodPost, dynamoDBURL, resp)

			repo := repository.New(conn, "my-table", log)
			users, err := repo.FindAllDocuments(context.TODO())
			assert.Error(t, err)
			assert.Nil(t, users)
		})

		t.Run("return_all_documents_but_not_match_with_spec", func(t *testing.T) {
			result := `{
    "ConsumedCapacity": {
      "CapacityUnits": 0.5,
      "TableName": "my-table"
    },
    "Count": 2,
    "Items": [
  {
    "Age": {
      "S": "33"
    },
    "CreatedAt": {
      "S": "2024-04-14T13:44:37.609166-06:00"
    },
    "Email": {
      "S": "john.smith@test.com"
    },
    "Id": {
      "S": "1"
    },
    "Lastname": {
      "S": "smith"
    },
    "Name": {
      "S": "john"
    },
    "UpdatedAt": {
      "S": "2024-04-14T13:44:37.609169-06:00"
    }
  },
  {
    "Age": {
      "S": "25"
    },
    "CreatedAt": {
      "S": "2024-04-14T13:44:37.609169-06:00"
    },
    "Email": {
      "S": "john.smith2@test.com"
    },
    "Id": {
      "S": "2"
    },
    "Lastname": {
      "S": "smith"
    },
    "Name": {
      "S": "john"
    },
    "UpdatedAt": {
      "S": "2024-04-14T13:44:37.609169-06:00"
    }
  }],
    "ScannedCount": 4
  }`

			resp := httpmock.NewStringResponder(http.StatusOK, result)
			httpmock.RegisterResponder(http.MethodPost, dynamoDBURL, resp)

			repo := repository.New(conn, "my-table", log)
			users, err := repo.FindAllDocuments(context.TODO())
			assert.Error(t, err)
			assert.Nil(t, users)
		})
	})
}

func TestRepositoryImpl_DockerTest(t *testing.T) {
	dynamodbTableName := "users"
	log := logger.NewLoggerWithOptions(logger.Opts{
		AppName: "get-all-documents-lambda-repository-test",
		Level:   "debug",
	})

	dynamodbSuite := tests.New("9000")
	err := dynamodbSuite.StartDynamoDB()
	assert.NoError(t, err)

	defer dynamodbSuite.Shutdown()

	// create dynamodb table
	err = dynamodbSuite.CreateTable(getTable(dynamodbTableName))
	assert.NoError(t, err)

	// set items
	items := []*models.UserDB{
		{
			ID:        uuid.NewString(),
			Name:      "john",
			Lastname:  "smith",
			Age:       34,
			Email:     "john.smith@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.NewString(),
			Name:      "josep",
			Lastname:  "smith",
			Age:       18,
			Email:     "josep.smith@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.NewString(),
			Name:      "milli",
			Lastname:  "smith",
			Age:       24,
			Email:     "mlli.smith@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.NewString(),
			Name:      "jose",
			Lastname:  "hernandez",
			Age:       20,
			Email:     "jose.hernandez@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.NewString(),
			Name:      "martin",
			Lastname:  "caballero",
			Age:       17,
			Email:     "martin.caballero@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, item := range items {
		err = dynamodbSuite.PutItem(dynamodbTableName, item)
		assert.NoError(t, err)
	}

	// init local conn
	err = os.Setenv("ENV", "local")
	assert.NoError(t, err)

	db := dynamodb2.New()
	conn, err := db.Connect()
	assert.NoError(t, err)

	// configure repository
	repo := repository.New(conn, dynamodbTableName, log)

	users, err := repo.FindAllDocuments(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 5)
}

func getTable(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("CreatedAt"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("CreatedAt"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:   aws.String(tableName),
		BillingMode: types.BillingModePayPerRequest,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		Tags: []types.Tag{
			{
				Key:   aws.String("OWNER"),
				Value: aws.String("Ricardo Romero"),
			},
		},
	}
}

func toString(input any) string {
	data, err := json.Marshal(input)
	if err != nil {
		return ""
	}

	return string(data)
}
