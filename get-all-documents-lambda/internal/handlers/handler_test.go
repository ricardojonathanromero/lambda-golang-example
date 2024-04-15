package handlers_test

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/business/models"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"time"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) LookingUpUsers(ctx context.Context) ([]*models.UserDB, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.UserDB), args.Error(1)
}

func TestHandleImpl_HandleRequest(t *testing.T) {
	log := logger.NewLoggerWithOptions(logger.Opts{
		AppName: "get-all-documents-lambda-handler-test",
		Level:   "debug",
	})

	mockService := new(MockService)

	var ct = &lambdacontext.LambdaContext{
		AwsRequestID:       "awsRequestId1234",
		InvokedFunctionArn: "arn:aws:lambda:xxx",
		Identity:           lambdacontext.CognitoIdentity{},
		ClientContext:      lambdacontext.ClientContext{},
	}

	var ctx = lambdacontext.NewContext(context.TODO(), ct)

	t.Run("when_receives_request", func(t *testing.T) {
		t.Run("then_return_200", func(t *testing.T) {
			t.Run("and_2_documents", func(t *testing.T) {
				req := events.APIGatewayProxyRequest{
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
					}, nil)

				res, err := handlers.New(mockService, log).HandleRequest(ctx, req)
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				assert.Equal(t, res.Headers["Content-Type"], "application/json")
				assert.Equal(t, res.StatusCode, http.StatusOK)
			})
		})

		t.Run("then_return_error", func(t *testing.T) {
			req := events.APIGatewayProxyRequest{
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

			res, err := handlers.New(mockService, log).HandleRequest(ctx, req)
			assert.NoError(t, err)
			assert.NotEmpty(t, res)
			assert.Equal(t, res.Headers["Content-Type"], "application/json")
			assert.Equal(t, res.StatusCode, http.StatusConflict)
		})
	})
}
