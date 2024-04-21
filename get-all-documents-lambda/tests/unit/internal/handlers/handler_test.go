package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/internal/handler"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/models"
	"github.com/stretchr/testify/mock"
	"net/http"
	"time"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) LookingUpUsers(ctx context.Context) ([]*models.UserDB, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.UserDB), args.Error(1)
}

var _ = Describe("Handler", func() {
	var mockService *MockService
	var lambdaCtx *lambdacontext.LambdaContext
	var ctx context.Context
	var log logger.Logger

	appName := "get-all-documents-lambda-handler-test"
	logLevel := "debug"

	BeforeEach(func() {
		log = logger.NewLoggerWithOptions(logger.Opts{AppName: appName, Level: logLevel})
		mockService = new(MockService)
		lambdaCtx = &lambdacontext.LambdaContext{
			AwsRequestID:       "awsRequestId1234",
			InvokedFunctionArn: "arn:aws:lambda:xxx",
			Identity:           lambdacontext.CognitoIdentity{},
			ClientContext:      lambdacontext.ClientContext{},
		}
	})

	Describe("process request event", func() {
		Context("with context timeout of 10s", func() {
			var c context.Context
			var cancel context.CancelFunc

			BeforeEach(func() {
				c, cancel = context.WithTimeout(context.Background(), 10*time.Second)
				ctx = lambdacontext.NewContext(c, lambdaCtx)
			})

			Context("mocking success result from service layer", func() {
				var req events.APIGatewayProxyRequest
				var err error

				BeforeEach(func() {
					req = events.APIGatewayProxyRequest{
						Resource:   "/",
						Path:       "/",
						HTTPMethod: http.MethodGet,
						Headers: map[string]string{
							"Accept": "application/json",
						},
					}

					mockService.On("LookingUpUsers", ctx).
						Times(1).
						Return([]*models.UserDB{
							{
								ID:        "1",
								Name:      "john",
								Lastname:  "smith",
								Age:       28,
								Email:     "john.smith@test.com",
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							},
						}, err)
				})

				It("can get 200 http code from response", func() {
					defer cancel()

					res, errRes := handler.New(mockService, log).HandleRequest(ctx, req)
					Expect(errRes).To(BeNil())
					Expect(res).NotTo(BeNil())
					Expect(res.Headers).To(HaveKeyWithValue("Content-Type", "application/json"))
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					Expect(res.Body).NotTo(BeEmpty())

					var expectRes []*models.UserDB
					err = json.Unmarshal([]byte(res.Body), &expectRes)
					Expect(err).To(BeNil())
					Expect(expectRes).To(HaveLen(1))
					Expect(expectRes[0].ID).To(Equal("1"))
				})
			})

			When("service fails", func() {
				var req events.APIGatewayProxyRequest

				BeforeEach(func() {
					req = events.APIGatewayProxyRequest{
						Resource:   "/",
						Path:       "/",
						HTTPMethod: http.MethodGet,
						Headers: map[string]string{
							"Accept": "application/json",
						},
					}

					var result []*models.UserDB
					mockService.On("LookingUpUsers", ctx).
						Times(1).
						Return(result, errors.New("internal error"))
				})

				It("can get conflict response", func() {
					defer cancel()

					res, errRes := handler.New(mockService, log).HandleRequest(ctx, req)
					Expect(errRes).To(BeNil())
					Expect(res).NotTo(BeNil())
					Expect(res.Headers).To(HaveKeyWithValue("Content-Type", "application/json"))
					Expect(res.StatusCode).To(Equal(http.StatusConflict))
					Expect(res.Body).NotTo(BeEmpty())

					var expectRes map[string]string
					err := json.Unmarshal([]byte(res.Body), &expectRes)
					Expect(err).To(BeNil())
					Expect(expectRes).To(HaveKeyWithValue("message", "internal error"))
				})
			})
		})
	})
})
