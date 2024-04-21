package services_test

import (
	"context"
	"errors"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/services"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindAllDocuments(ctx context.Context) ([]*models.UserDB, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.UserDB), args.Error(1)
}

func TestServiceImpl_LookingUpUsers(t *testing.T) {
	log := logger.NewLoggerWithOptions(logger.Opts{
		AppName: "get-all-documents-lambda-service-test",
		Level:   "debug",
	})

	mockRepo := new(MockRepo)

	t.Run("when_receives_request", func(t *testing.T) {
		t.Run("then_return_all_documents", func(t *testing.T) {
			ctx := context.TODO()

			mockRepo.On("FindAllDocuments", ctx).
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

			users, err := services.New(mockRepo, log).LookingUpUsers(ctx)
			assert.NoError(t, err)
			assert.NotNil(t, users)
			assert.Len(t, users, 1)
		})

		t.Run("and_db_connection_is_not_initialized", func(t *testing.T) {
			t.Run("then_return_an_error", func(t *testing.T) {
				ctx := context.TODO()

				var resp []*models.UserDB
				result := errors.New("db not connected")

				mockRepo.On("FindAllDocuments", ctx).
					Times(1).
					Return(resp, result)

				users, err := services.New(mockRepo, log).LookingUpUsers(ctx)
				assert.Nil(t, users)
				assert.NotNil(t, err)
				assert.Equal(t, err, result)
			})
		})
	})
}
