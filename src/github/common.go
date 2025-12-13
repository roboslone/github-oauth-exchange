package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	githubv1 "github.com/roboslone/github-oauth-exchange/proto/github/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

type RawTokenResponse struct {
	Status                int
	Error                 string `json:"error"`
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	TokenType             string `json:"token_type"`
	Scope                 string `json:"scope"`
}

func (r *RawTokenResponse) ToExchangeResponse() *githubv1.ExchangeResponse {
	return &githubv1.ExchangeResponse{
		AccessToken: &githubv1.Token{
			Value:     r.AccessToken,
			ExpiresIn: durationpb.New(time.Second * time.Duration(r.ExpiresIn)),
		},
		RefreshToken: &githubv1.Token{
			Value:     r.RefreshToken,
			ExpiresIn: durationpb.New(time.Second * time.Duration(r.RefreshTokenExpiresIn)),
		},
		TokenType: r.TokenType,
		Scope:     r.Scope,
	}
}

func (r *RawTokenResponse) ToRefreshResponse() *githubv1.RefreshResponse {
	return &githubv1.RefreshResponse{
		AccessToken: &githubv1.Token{
			Value:     r.AccessToken,
			ExpiresIn: durationpb.New(time.Second * time.Duration(r.ExpiresIn)),
		},
		RefreshToken: &githubv1.Token{
			Value:     r.RefreshToken,
			ExpiresIn: durationpb.New(time.Second * time.Duration(r.RefreshTokenExpiresIn)),
		},
		TokenType: r.TokenType,
		Scope:     r.Scope,
	}
}

func requestTokens(ctx context.Context, data []byte) (*RawTokenResponse, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://github.com/login/oauth/access_token",
		strings.NewReader(string(data)),
	)
	if err != nil {
		return nil, fmt.Errorf("constructing refresh request: %w", err)
	}

	request.Header = http.Header{
		"content-type": []string{"application/json"},
		"accept":       []string{"application/json"},
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("exchanging code: %w", err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var result RawTokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}
	result.Status = response.StatusCode

	if result.Error != "" || result.Status != http.StatusOK {
		return nil, fmt.Errorf("bad github response: %d: %s", result.Status, result.Error)
	}
	if result.AccessToken == "" {
		return nil, fmt.Errorf("github response doesn't contain access token")
	}
	return &result, nil
}
