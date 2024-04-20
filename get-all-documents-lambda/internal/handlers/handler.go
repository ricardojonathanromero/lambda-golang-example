package handlers

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/pkg/services"
	"github.com/ricardojonathanromero/lambda-golang-example/utils/encoding"
	"net/http"
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

func (h *handleImpl) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	h.log.Debug("handleRequest")
	users, err := h.srv.LookingUpUsers(ctx)
	if err != nil {
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
