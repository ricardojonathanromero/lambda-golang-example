package one_record_test

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/internal/handler"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/repository"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/service"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/models"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/utils/tests"
	"net/http"
	"time"
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

var _ = Describe("one record", func() {
	var hdl handler.Handler
	var lambdaCtx *lambdacontext.LambdaContext
	var ctx context.Context

	BeforeEach(func() {
		repo := repository.New(conn, tableName, log)
		srv := service.New(repo, log)
		hdl = handler.New(srv, log)

		lambdaCtx = &lambdacontext.LambdaContext{
			AwsRequestID:       "awsRequestId1234",
			InvokedFunctionArn: "arn:aws:lambda:xxx",
			Identity:           lambdacontext.CognitoIdentity{},
			ClientContext:      lambdacontext.ClientContext{},
		}
	})

	Describe("iterate records", func() {
		When("deadline is up to 10 secs", func() {
			var cancel context.CancelFunc

			BeforeEach(func() {
				// configure
				ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
				ctx = lambdacontext.NewContext(ctx, lambdaCtx)
			})

			Context("set one item in db", func() {
				item := &models.UserDB{
					ID:        uuid.NewString(),
					Name:      "john",
					Lastname:  "smith",
					Age:       23,
					Email:     "john.smith@test.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				BeforeEach(func() {
					err := dynamodbTestConn.PutItem(tableName, item)
					Expect(err).To(BeNil())
				})

				It("return empty list", func() {
					defer cancel()

					req := events.APIGatewayProxyRequest{
						Resource:   "/",
						Path:       "/",
						HTTPMethod: http.MethodGet,
					}

					log.Debug("start send handler request")
					res, errRes := hdl.HandleRequest(ctx, req)

					Expect(errRes).To(BeNil())
					Expect(res).NotTo(BeNil())
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					Expect(res.Body).NotTo(BeEmpty())

					var result []*models.UserDB
					err := json.Unmarshal([]byte(res.Body), &result)
					Expect(err).To(BeNil())
					Expect(result).To(HaveLen(1))
					Expect(result[0].ID).To(Equal(item.ID))
				})
			})
		})
	})
})
