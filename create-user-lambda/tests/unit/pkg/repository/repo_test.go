package repository_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/repository"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/models"
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

type UnsupportedType struct {
	Channel chan int
}

// MarshalDynamoDBAttributeValue implements the Marshaller interface incorrectly
func (u UnsupportedType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return nil, fmt.Errorf("unsupported type: channel")
}

var _ = Describe("Repository", func() {
	var ctx context.Context
	var log logger.Logger
	var conn *dynamodb.Client

	appName := "create-user-lambda-repository-test"
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

	Describe("receives a user to insert", func() {
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
					result := `{}`
					resp := httpmock.NewStringResponder(http.StatusOK, result)
					httpmock.RegisterResponder(http.MethodPost, dynamodbLocalURL, resp)
					repo = repository.New(tableName, conn, log)
				})

				It("can be finish the process without any error", func() {
					defer cancel()

					usr := models.UserDB{
						ID:        uuid.NewString(),
						Name:      "john",
						Lastname:  "smith",
						Age:       30,
						Email:     "john.smith@test.com",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}

					err := repo.InsertUser(ctx, usr)
					Expect(err).To(BeNil())
				})
			})

			Context("the db return not valid response", func() {
				When("set an invalid structure", func() {
					//var result string
					var repo repository.Repository

					BeforeEach(func() {
						//result = `{"code":"ConditionalCheckFailedException","message":"The id set already exists"}`
						//resp := httpmock.NewStringResponder(http.StatusBadRequest, result)
						//httpmock.RegisterResponder(http.MethodPost, dynamodbLocalURL, resp)
						repo = repository.New(tableName, conn, log)
					})

					It("cannot be marshalled due to unsupported channel type", func() {
						defer cancel()

						unsupported := UnsupportedType{
							Channel: make(chan int),
						}

						err := repo.InsertUser(ctx, unsupported)
						Expect(err).NotTo(BeNil())
						Expect(err.Error()).To(Equal("unsupported type: channel"))
					})
				})

				When("id set already exists", func() {
					//var result string
					var repo repository.Repository

					BeforeEach(func() {
						result := `{"code":"ConditionalCheckFailedException","message":"The id set already exists"}`
						resp := httpmock.NewStringResponder(http.StatusBadRequest, result)
						httpmock.RegisterResponder(http.MethodPost, dynamodbLocalURL, resp)
						repo = repository.New(tableName, conn, log)
					})

					It("cannot be marshalled due to unsupported channel type", func() {
						defer cancel()

						usr := models.UserDB{
							ID:        uuid.NewString(),
							Name:      "john",
							Lastname:  "smith",
							Age:       30,
							Email:     "john.smith@test.com",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}

						err := repo.InsertUser(ctx, usr)
						Expect(err).NotTo(BeNil())

						var ae smithy.APIError
						ok := errors.As(err, &ae)
						Expect(ok).To(BeTrue())
						Expect(ae).To(HaveExistingField("Message"))
						Expect(ae.ErrorCode()).To(Equal("ConditionalCheckFailedException"))
						Expect(ae.ErrorMessage()).To(Equal("The id set already exists"))
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
					repo = repository.New(tableName, conn, log)
				})

				It("cannot send request by timeout", func() {
					defer cancel()

					time.Sleep(2 * time.Second) // sleep 2 secs

					usr := models.UserDB{
						ID:        uuid.NewString(),
						Name:      "john",
						Lastname:  "smith",
						Age:       30,
						Email:     "john.smith@test.com",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}

					err := repo.InsertUser(ctx, usr)
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
