package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	githubv1 "github.com/roboslone/github-oauth-exchange-proto/github/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

type RawExchangeRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type RawExchangeResponse struct {
	Status                int
	Error                 string `json:"error"`
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	TokenType             string `json:"token_type"`
	Scope                 string `json:"scope"`
}

func (r *RawExchangeResponse) ToProto() *githubv1.ExchangeResponse {
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

func ExchangeCode(ctx context.Context, app Application, code string) (*githubv1.ExchangeResponse, error) {
	data, err := json.Marshal(RawExchangeRequest{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		Code:         code,
	})
	if err != nil {
		return nil, fmt.Errorf("encoding exchange request data: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://github.com/login/oauth/access_token",
		strings.NewReader(string(data)),
	)
	if err != nil {
		return nil, fmt.Errorf("constructing exhange request: %w", err)
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

	var result RawExchangeResponse
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

	return result.ToProto(), nil
}
