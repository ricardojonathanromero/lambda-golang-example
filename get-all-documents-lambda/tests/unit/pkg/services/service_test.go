package services_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/service"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/models"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindAllDocuments(ctx context.Context) ([]*models.UserDB, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.UserDB), args.Error(1)
}

var _ = Describe("Service", func() {
	Expect(nil)
	var mockRepo *MockRepo
	var log logger.Logger
	var ctx context.Context

	BeforeEach(func() {
		mockRepo = new(MockRepo)
		log = logger.NewLoggerWithOptions(logger.Opts{
			AppName: "create-user-lambda-service-test",
			Level:   "debug",
		})
		ctx = context.Background()
	})

	Describe("service return response", func() {
		Context("deadline is up to 10 secs", func() {
			var cancel context.CancelFunc
			BeforeEach(func() {
				ctx, cancel = context.WithTimeout(ctx, time.Second*10)
			})

			When("request is valid and mock valid response from db", func() {
				BeforeEach(func() {
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
				})

				It("can save the record", func() {
					defer cancel()

					users, err := service.New(mockRepo, log).LookingUpUsers(ctx)
					Expect(err).To(BeNil())
					Expect(users).NotTo(BeNil())
					Expect(users).To(HaveLen(1))
					Expect(users[0].ID).To(Equal("1"))
				})
			})

			When("db returns an error", func() {
				BeforeEach(func() {
					var resp []*models.UserDB
					mockRepo.On("FindAllDocuments", ctx).
						Times(1).
						Return(resp, context.DeadlineExceeded)
				})

				It("can save the record", func() {
					defer cancel()

					users, err := service.New(mockRepo, log).LookingUpUsers(ctx)
					Expect(users).To(BeNil())
					Expect(err).NotTo(BeNil())
					Expect(err).To(Equal(context.DeadlineExceeded))
				})
			})
		})
	})
})
