package service

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	githubv1 "github.com/roboslone/github-oauth-exchange/proto/github/v1"
	"github.com/roboslone/github-oauth-exchange/proto/github/v1/githubv1connect"
	"github.com/roboslone/github-oauth-exchange/src/github"
)

var (
	ErrAccessTokenNotProvided = connect.NewError(connect.CodeUnauthenticated, errors.New("access_token not provided"))
)

type Service struct {
	githubv1connect.UnimplementedExchangeServiceHandler

	cfg *Config
}

func New(cfg *Config) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) Exchange(ctx context.Context, request *connect.Request[githubv1.ExchangeRequest]) (*connect.Response[githubv1.ExchangeResponse], error) {
	app, ok := s.cfg.GitHub.Index[request.Msg.ClientId]
	if !ok {
		return nil, connect.NewError(
			connect.CodeNotFound,
			errors.New(request.Msg.ClientId),
		)
	}

	response, err := github.ExchangeCode(ctx, app, request.Msg.Code)
	if err != nil {
		return nil, err
	}

	if request.Msg.Resolve {
		response.Account, err = github.Resolve(ctx, response.GetAccessToken().GetValue())
		if err != nil {
			return nil, fmt.Errorf("resolving access token: %w", err)
		}
	}

	return connect.NewResponse(response), nil
}

func (s *Service) Resolve(ctx context.Context, request *connect.Request[githubv1.ResolveRequest]) (*connect.Response[githubv1.ResolveResponse], error) {
	if request.Msg.AccessToken == "" {
		return nil, ErrAccessTokenNotProvided
	}

	account, err := github.Resolve(ctx, request.Msg.AccessToken)
	return connect.NewResponse(&githubv1.ResolveResponse{Account: account}), err
}
