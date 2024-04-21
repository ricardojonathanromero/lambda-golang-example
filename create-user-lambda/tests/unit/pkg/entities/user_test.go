package entities_test

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/entities"
	"os"
)

var _ = Describe("Entities", func() {
	Context("marshal entity into db model", func() {
		When("use default timezone location", func() {
			When("declare user", func() {
				user := &entities.UserReq{
					Name:     "john",
					Lastname: "smith",
					Age:      30,
					Email:    "john.smith@test.com",
				}

				It("can marshal the entity", func() {
					dbModel, err := user.ToDB()
					Expect(err).To(BeNil())
					Expect(dbModel).NotTo(BeNil())
					Expect(dbModel).To(HaveExistingField("ID"))
					Expect(dbModel).To(HaveExistingField("Name"))
					Expect(dbModel).To(HaveExistingField("Lastname"))
					Expect(dbModel).To(HaveExistingField("Age"))
					Expect(dbModel).To(HaveExistingField("Email"))
					Expect(dbModel).To(HaveExistingField("CreatedAt"))
					Expect(dbModel).To(HaveExistingField("UpdatedAt"))
					Expect(dbModel.ID).NotTo(BeEmpty())
					Expect(dbModel.CreatedAt).NotTo(BeNil())
					Expect(dbModel.UpdatedAt).NotTo(BeNil())
				})
			})
		})

		When("set custom timezone location", func() {
			When("declare user and valid timezone", func() {
				var user *entities.UserReq
				BeforeEach(func() {
					user = &entities.UserReq{
						Name:     "john",
						Lastname: "smith",
						Age:      30,
						Email:    "john.smith@test.com",
					}

					err := os.Setenv("TZ_LOCATION", "America/Mexico_City")
					Expect(err).To(BeNil())
				})

				AfterEach(func() {
					val := os.Getenv("TZ_LOCATION")
					if len(val) > 0 {
						err := os.Unsetenv("TZ_LOCATION")
						Expect(err).To(BeNil())
					}
				})

				It("can marshal the entity", func() {
					dbModel, err := user.ToDB()
					Expect(err).To(BeNil())
					Expect(dbModel).NotTo(BeNil())
					Expect(dbModel).To(HaveExistingField("ID"))
					Expect(dbModel).To(HaveExistingField("Name"))
					Expect(dbModel).To(HaveExistingField("Lastname"))
					Expect(dbModel).To(HaveExistingField("Age"))
					Expect(dbModel).To(HaveExistingField("Email"))
					Expect(dbModel).To(HaveExistingField("CreatedAt"))
					Expect(dbModel).To(HaveExistingField("UpdatedAt"))
					Expect(dbModel.ID).NotTo(BeEmpty())
					Expect(dbModel.CreatedAt).NotTo(BeNil())
					Expect(dbModel.UpdatedAt).NotTo(BeNil())
				})
			})

			When("declare user and invalid timezone", func() {
				var user *entities.UserReq
				BeforeEach(func() {
					user = &entities.UserReq{
						Name:     "john",
						Lastname: "smith",
						Age:      30,
						Email:    "john.smith@test.com",
					}

					err := os.Setenv("TZ_LOCATION", "NotValid")
					Expect(err).To(BeNil())
				})

				AfterEach(func() {
					val := os.Getenv("TZ_LOCATION")
					if len(val) > 0 {
						err := os.Unsetenv("TZ_LOCATION")
						Expect(err).To(BeNil())
					}
				})

				It("can marshal the entity", func() {
					dbModel, err := user.ToDB()
					Expect(dbModel).To(BeNil())
					Expect(err).NotTo(BeNil())
					unwErr := errors.Unwrap(err)
					Expect(unwErr.Error()).To(Equal("unknown time zone NotValid"))
				})
			})
		})
	})
})
