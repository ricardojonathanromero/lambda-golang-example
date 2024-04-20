package entities

import (
	"github.com/google/uuid"
	"github.com/ricardojonathanromero/lambda-golang-example/business/models"
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

func (u *UserReq) ToDB() (*models.UserDB, error) {
	u.ID = uuid.NewString()

	lc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, err
	}

	now := time.Now().In(lc)
	u.CreatedAt = now
	u.UpdatedAt = now

	return (*models.UserDB)(u), nil
}
