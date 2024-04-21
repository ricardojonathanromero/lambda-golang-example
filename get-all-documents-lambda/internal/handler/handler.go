package handler

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/service"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/utils/encoding"
	"net/http"
)

type Handler interface {
	HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type handleImpl struct {
	srv service.Service
	log logger.Logger
}

func New(srv service.Service, log logger.Logger) Handler {
	return &handleImpl{
		srv: srv,
		log: log,
	}
}

func (h *handleImpl) HandleRequest(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	h.log.Debug("handleRequest")
	users, err := h.srv.LookingUpUsers(ctx)
	if err != nil {
		h.log.Errorf("error from service: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusConflict,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: fmt.Sprintf(`{"message": "%s"}`, err),
		}, nil
	}

	h.log.Debug("success response!")
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: encoding.ToString(users),
	}, nil
}
