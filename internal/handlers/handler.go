package handlers

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/pkg/services"
)

type Handler interface {
	HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type handleImpl struct {
	srv services.Service
	log logger.Logger
}

func New(srv services.Service, log logger.Logger) Handler {
	return &handleImpl{
		srv: srv,
		log: log,
	}
}

// implement methods

func (h handleImpl) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//TODO implement me
	panic("implement me")
}
