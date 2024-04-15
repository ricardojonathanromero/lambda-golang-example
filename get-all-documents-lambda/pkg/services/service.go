package services

import (
	"context"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/business/models"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/internal/utils"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/repository"
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

	srv.log.Debug(utils.ToString(users))
	return users, err
}
