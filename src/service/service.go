package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	githubv1 "github.com/roboslone/github-oauth-exchange-proto/github/v1"
	"github.com/roboslone/github-oauth-exchange-proto/github/v1/githubv1connect"
	"github.com/roboslone/github-oauth-exchange/src/config"
)

type Service struct {
	githubv1connect.UnimplementedExchangeServiceHandler

	cfg *config.Config
}

func New(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) Exchange(ctx context.Context, request *connect.Request[githubv1.ExchangeRequest]) (*connect.Response[githubv1.ExchangeResponse], error) {
	app, ok := s.cfg.GitHub.Index[request.Msg.AppName]
	if !ok {
		return nil, connect.NewError(
			connect.CodeNotFound,
			errors.New(request.Msg.AppName),
		)
	}

	response, err := ExchangeCode(ctx, app, request.Msg.Code)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(response), nil
}
