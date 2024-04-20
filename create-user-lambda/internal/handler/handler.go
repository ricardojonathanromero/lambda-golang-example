package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/entities"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/service"
	"net/http"
)

type Handle interface {
	HandleCreateUser(ctx context.Context, req entities.UserReq) (events.APIGatewayProxyResponse, error)
}

type handleImpl struct {
	srv service.Service
	log logger.Logger
	v   *validator.Validate
}

func New(srv service.Service, log logger.Logger) Handle {
	return &handleImpl{
		srv: srv,
		log: log,
		v:   validator.New(),
	}
}

func (h *handleImpl) HandleCreateUser(ctx context.Context, req entities.UserReq) (events.APIGatewayProxyResponse, error) {
	var res events.APIGatewayProxyResponse
	h.log.Debug("event received")

	h.log.Debug("validating request")
	if err := h.v.StructCtx(ctx, req); err != nil {
		h.log.Errorf("error occurs validating struct: %v", err)
		return h.getErrorResponse(err), nil
	}

	h.log.Debug("creating user")
	err := h.srv.CreateUser(ctx, req)
	if err != nil {
		h.log.Errorf("error creating user: %v", err)
		return h.getErrorResponse(err), nil
	}

	h.log.Info("event processed")
	res = events.APIGatewayProxyResponse{StatusCode: http.StatusCreated}
	return res, nil
}

func (h *handleImpl) getErrorResponse(err error) events.APIGatewayProxyResponse {
	var statusCode int
	var data string
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		statusCode = http.StatusBadRequest
		data = fmt.Sprintf(`{"code": "bad_request", "message": "%s"}`, ve)
	} else {
		statusCode = http.StatusConflict
		data = fmt.Sprintf(`{"code": "conflict", "message": "%s"}`, err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: data,
	}
}
