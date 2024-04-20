package service

import (
	"context"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/entities"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/repository"
)

type Service interface {
	CreateUser(ctx context.Context, req entities.UserReq) error
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

func (s *serviceImpl) CreateUser(ctx context.Context, req entities.UserReq) error {
	s.log.Debug("converting req model into db model")
	dbReq, err := req.ToDB()
	if err != nil {
		s.log.Errorf("error loading location: %v", err)
		return err
	}

	s.log.Info("saving request")
	err = s.repo.InsertUser(ctx, dbReq)
	if err != nil {
		s.log.Errorf("error inserting user: %v", err)
		return err
	}

	s.log.Info("record saved")
	return nil
}
