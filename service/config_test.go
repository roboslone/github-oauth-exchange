package service_test

import (
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/roboslone/github-oauth-exchange/github"
	"github.com/roboslone/github-oauth-exchange/service"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	cfg, err := service.NewConfig(
		service.WithEnvOptions(&env.Options{
			Environment: map[string]string{
				"GITHUB__APPLICATIONS": "ai:as,bi:bs",
			},
		}),
	)
	require.NoError(t, err)

	require.Equal(
		t,
		map[string]github.Application{
			"ai": {ClientID: "ai", ClientSecret: "as"},
			"bi": {ClientID: "bi", ClientSecret: "bs"},
		},
		cfg.GitHub.Index,
	)
}
