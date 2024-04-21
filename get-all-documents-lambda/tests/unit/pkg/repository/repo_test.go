package repository_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/repository"
	"net/http"
	"time"
)

func getDBClientWithHttpHandler(url string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...any) (aws.Endpoint, error) {
				return aws.Endpoint{URL: url}, nil
			})),
		config.WithHTTPClient(http.DefaultClient),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummyKey", "dummySecret", "")),
	)

	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

var _ = Describe("Repository", func() {
	var ctx context.Context
	var log logger.Logger
	var conn *dynamodb.Client

	appName := "get-all-documents-lambda-repository-test"
	dynamodbLocalURL := "http://localhost:8000/"
	tableName := "my-table"
	logLevel := "debug"

	BeforeEach(func() {
		var err error
		// configure dynamodb local session
		ctx = context.Background()
		log = logger.NewLoggerWithOptions(logger.Opts{AppName: appName, Level: logLevel})
		conn, err = getDBClientWithHttpHandler(dynamodbLocalURL)
		Expect(err).To(BeNil())
	})

	Describe("retrieve all users in db", func() {
		When("connection db has been initialized and context deadline is set to 10 secs", func() {
			var cancel context.CancelFunc
			BeforeEach(func() {
				// remove any mocks
				httpmock.Reset()
				ctx, cancel = context.WithTimeout(ctx, time.Second*10)
			})

			Context("the db returns success response", func() {
				var repo repository.Repository

				BeforeEach(func() {
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
					httpmock.RegisterResponder(http.MethodPost, dynamodbLocalURL, resp)
					repo = repository.New(conn, tableName, log)
				})

				It("can get 4 elements", func() {
					defer cancel()

					users, err := repo.FindAllDocuments(ctx)
					Expect(err).To(BeNil())
					Expect(users).NotTo(BeNil())
					Expect(users).To(HaveLen(4))
				})
			})

			Context("the db return not valid response", func() {
				When("occurs a transaction conflict exception", func() {
					//var result string
					var repo repository.Repository

					BeforeEach(func() {
						result := `{"code":"TransactionConflictException","message":"A conflict occurs trying to scan documents"}`
						resp := httpmock.NewStringResponder(http.StatusBadRequest, result)
						httpmock.RegisterResponder(http.MethodPost, dynamodbLocalURL, resp)
						repo = repository.New(conn, tableName, log)
					})

					It("cannot be marshalled due to unsupported channel type", func() {
						defer cancel()

						users, err := repo.FindAllDocuments(ctx)
						Expect(users).To(BeNil())
						Expect(err).NotTo(BeNil())

						var ae smithy.APIError
						ok := errors.As(err, &ae)
						Expect(ok).To(BeTrue())
						Expect(ae).To(HaveExistingField("Message"))
						Expect(ae.ErrorCode()).To(Equal("TransactionConflictException"))
						Expect(ae.ErrorMessage()).To(Equal("A conflict occurs trying to scan documents"))
					})
				})
				When("items received not match with model", func() {
					var repo repository.Repository

					BeforeEach(func() {
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
						httpmock.RegisterResponder(http.MethodPost, dynamodbLocalURL, resp)
						repo = repository.New(conn, tableName, log)
					})

					It("receives an error unmarshalling result", func() {
						defer cancel()

						users, err := repo.FindAllDocuments(ctx)
						Expect(users).To(BeNil())
						Expect(err).NotTo(BeNil())

						log.Debugf("previous to compare with object: %v", err)
						var ute *attributevalue.UnmarshalTypeError
						ok := errors.As(err, &ute)

						Expect(ok).To(BeTrue())
						log.Debugf("value: %s, type: %s, error: %v, unwrap: %v", ute.Value, ute.Value, ute.Error(), ute.Unwrap())
						expectedErr := fmt.Errorf("unmarshal failed, %s", "cannot unmarshal string into Go value type int32")
						Expect(ute.Error()).To(Equal(expectedErr.Error()))
					})
				})
			})
		})

		When("connection db has been initialized and context deadline is set to 1 secs", func() {
			var cancel context.CancelFunc
			BeforeEach(func() {
				// remove any mocks
				httpmock.Reset()
				ctx, cancel = context.WithTimeout(ctx, time.Second*1)
			})

			When("configure success response", func() {
				var repo repository.Repository

				BeforeEach(func() {
					result := `{}`
					resp := httpmock.NewStringResponder(http.StatusOK, result)
					httpmock.RegisterResponder(http.MethodPost, dynamodbLocalURL, resp)
					repo = repository.New(conn, tableName, log)
				})

				It("cannot send request by timeout", func() {
					defer cancel()

					time.Sleep(2 * time.Second) // sleep 2 secs

					users, err := repo.FindAllDocuments(ctx)
					Expect(users).To(BeNil())
					Expect(err).NotTo(BeNil())

					var oe *smithy.OperationError
					ok := errors.As(err, &oe)
					Expect(ok).To(BeTrue())

					unerr := oe.Unwrap()
					Expect(errors.Is(unerr, context.DeadlineExceeded)).To(BeTrue())
				})
			})
		})
	})
})
