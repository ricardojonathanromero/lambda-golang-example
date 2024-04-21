package handler_test

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/lambdacontext"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/internal/handler"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/entities"
	"github.com/stretchr/testify/mock"
	"net/http"
	"time"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateUser(ctx context.Context, req entities.UserReq) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

var _ = Describe("Handler", func() {
	var mockService *MockService
	var lambdaCtx *lambdacontext.LambdaContext
	var ctx context.Context
	var log logger.Logger

	appName := "create-user-lambda-handler-test"
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
				var req entities.UserReq
				var err error

				BeforeEach(func() {
					req = entities.UserReq{
						Name:     "John",
						Lastname: "Smith",
						Age:      30,
						Email:    "john.smith@test.com",
					}

					mockService.On("CreateUser", ctx, req).
						Times(1).
						Return(err)
				})

				It("can get 201 http code from response", func() {
					defer cancel()

					res, errRes := handler.New(mockService, log).HandleCreateUser(ctx, req)
					Expect(errRes).To(BeNil())
					Expect(res).NotTo(BeNil())
					Expect(res.StatusCode).To(Equal(http.StatusCreated))
				})
			})

			When("request not pass validations", func() {
				var req entities.UserReq

				BeforeEach(func() {
					req = entities.UserReq{
						Lastname: "Smith",
						Age:      30,
						Email:    "john.smith@test.com",
					}
				})

				It("can get 201 http code from response", func() {
					defer cancel()

					res, errRes := handler.New(mockService, log).HandleCreateUser(ctx, req)
					Expect(errRes).To(BeNil())
					Expect(res).NotTo(BeNil())
					Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				})
			})

			When("result from service is not valid", func() {
				var req entities.UserReq

				BeforeEach(func() {
					req = entities.UserReq{
						Name:     "john",
						Lastname: "Smith",
						Age:      30,
						Email:    "john.smith@test.com",
					}

					mockService.On("CreateUser", ctx, req).
						Times(1).
						Return(errors.New("generic error"))
				})

				It("can get 409 http code from response", func() {
					defer cancel()

					res, errRes := handler.New(mockService, log).HandleCreateUser(ctx, req)
					Expect(errRes).To(BeNil())
					Expect(res).NotTo(BeNil())
					Expect(res.StatusCode).To(Equal(http.StatusConflict))
				})
			})
		})
	})
})
