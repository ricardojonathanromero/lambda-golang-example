package service

import (
	"context"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-document-lambda/pkg/repository"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/models"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/utils/encoding"
)

type Service interface {
	LookingUpUser(ctx context.Context, id string) (*models.UserDB, error)
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

func (srv *serviceImpl) LookingUpUser(ctx context.Context, id string) (*models.UserDB, error) {
	srv.log.Debug("processing service layer")

	srv.log.Debug("looking document")
	user, err := srv.repo.FindDocumentById(ctx, id)
	if err != nil {
		srv.log.Errorf("error from repository: %s", err)
		return nil, err
	}

	srv.log.Debug(encoding.ToString(user))
	return user, nil
}
