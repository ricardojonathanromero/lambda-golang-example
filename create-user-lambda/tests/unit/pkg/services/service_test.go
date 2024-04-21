package services_test

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/entities"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/service"
	"github.com/stretchr/testify/mock"
	"os"
	"time"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) InsertUser(ctx context.Context, user any) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

var _ = Describe("Service", func() {
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
		var req entities.UserReq

		BeforeEach(func() {
			req = entities.UserReq{
				Name:     "john",
				Lastname: "smith",
				Age:      20,
				Email:    "john.smith@test.com",
			}
		})

		Context("deadline is up to 10 secs", func() {
			var cancel context.CancelFunc
			BeforeEach(func() {
				ctx, cancel = context.WithTimeout(ctx, time.Second*10)
			})

			When("request is valid and mock valid response from db", func() {
				BeforeEach(func() {
					mockRepo.On("InsertUser", ctx, mock.Anything).
						Times(1).
						Return(nil)
				})

				It("can save the record", func() {
					defer cancel()

					err := service.New(mockRepo, log).CreateUser(ctx, req)
					Expect(err).To(BeNil())
				})
			})

			When("request cannot be marshalled as db model", func() {
				BeforeEach(func() {
					err := os.Setenv("TZ_LOCATION", "NotValid")
					Expect(err).To(BeNil())
				})

				AfterEach(func() {
					err := os.Unsetenv("TZ_LOCATION")
					Expect(err).To(BeNil())
				})

				It("cannot save record in db due to timezone not exist", func() {
					defer cancel()

					err := service.New(mockRepo, log).CreateUser(ctx, req)
					Expect(err).NotTo(BeNil())
					unwErr := errors.Unwrap(err)
					Expect(unwErr.Error()).To(Equal("unknown time zone NotValid"))
				})
			})

			When("db returns an error", func() {
				BeforeEach(func() {
					mockRepo.On("InsertUser", ctx, mock.Anything).
						Times(1).
						Return(context.DeadlineExceeded)
				})

				It("cannot save the record due to deadline exceeded", func() {
					defer cancel()

					err := service.New(mockRepo, log).CreateUser(ctx, req)
					Expect(err).NotTo(BeNil())
					Expect(err).To(Equal(context.DeadlineExceeded))
				})
			})
		})
	})
})
