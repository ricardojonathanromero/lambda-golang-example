package entities

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/ricardojonathanromero/go-utilities/environment"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/models"
	"time"
)

type UserReq struct {
	ID        string    `json:"-"`
	Name      string    `json:"name" validate:"required,min=3,max=50"`
	Lastname  string    `json:"lastname" validate:"required,min=3,max=50"`
	Age       int32     `json:"age" validate:"required,gt=0,lt=99"`
	Email     string    `json:"email" validate:"required,email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

const (
	envLocation     = "TZ_LOCATION"
	defaultLocation = "America/Los_Angeles"
)

func (u *UserReq) ToDB() (*models.UserDB, error) {
	u.ID = uuid.NewString()

	tz := environment.GetEnv(envLocation, defaultLocation)
	lc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("location %s not found: %w", tz, err)
	}

	now := time.Now().In(lc)
	u.CreatedAt = now
	u.UpdatedAt = now

	return (*models.UserDB)(u), nil
}
