package services

import (
	"context"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/repository"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/models"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/utils/encoding"
)

type Service interface {
	LookingUpUsers(ctx context.Context) ([]*models.UserDB, error)
}

type serviceImpl struct {
	repo repository.Repository
	log  logger.Logger
}

func New(repo repository.Repository, log logger.Logger) Service {
	return &serviceImpl{
		repo: repo,
		log:  log,
	}
}

func (srv *serviceImpl) LookingUpUsers(ctx context.Context) ([]*models.UserDB, error) {
	srv.log.Debug("looking for all users")
	users, err := srv.repo.FindAllDocuments(ctx)
	if err != nil {
		// do something
		srv.log.Errorf("error from repository: %s", err)
		return nil, err
	}

	srv.log.Debug(encoding.ToString(users))
	return users, err
}
