package repository_test

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jarcoal/httpmock"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/repository"
	"github.com/stretchr/testify/assert"
	"net/http"
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

func toString(input any) string {
	data, err := json.Marshal(input)
	if err != nil {
		return ""
	}

	return string(data)
}
