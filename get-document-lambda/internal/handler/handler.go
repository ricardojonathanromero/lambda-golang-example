package handler

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/get-document-lambda/pkg/service"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/utils/encoding"
	"net/http"
)

type Handler interface {
	HandleGetUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
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

func (h *handleImpl) HandleGetUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	h.log.Debug("handleRequest")

	h.log.Info("handle request")

	id, ok := req.PathParameters["id"]
	if !ok {
		h.log.Errorf("id is not valid: %s", id)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusConflict,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"message": "id is required"}`,
		}, nil
	}

	h.log.Debugf("looking for user: %s", id)
	result, err := h.srv.LookingUpUser(ctx, id)
	if err != nil {
		h.log.Errorf("error response from service: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusConflict,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: fmt.Sprintf(`{"message": "%s"}`, err),
		}, nil
	}

	h.log.Info("success response")

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: encoding.ToString(result),
	}, nil
}
