package github

import (
	"context"
	"encoding/json"
	"fmt"

	githubv1 "github.com/roboslone/github-oauth-exchange/proto/github/v1"
)

type RawRefreshRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"`
}

func RefreshToken(ctx context.Context, app Application, refreshToken string) (*githubv1.RefreshResponse, error) {
	data, err := json.Marshal(RawRefreshRequest{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("encoding request data: %w", err)
	}

	result, err := requestTokens(ctx, data)
	if err != nil {
		return nil, err
	}
	return result.ToRefreshResponse(), nil
}
