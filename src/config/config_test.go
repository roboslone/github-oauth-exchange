package config_test

import (
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/roboslone/github-oauth-exchange/src/config"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	cfg, err := config.New(
		config.WithEnvOptions(&env.Options{
			Environment: map[string]string{
				"GITHUB__APPLICATIONS": "a:ai:as,b:bi:bs",
			},
		}),
	)
	require.NoError(t, err)

	require.Equal(
		t,
		map[string]config.GitHubApplication{
			"a": {Name: "a", ClientID: "ai", ClientSecret: "as"},
			"b": {Name: "b", ClientID: "bi", ClientSecret: "bs"},
		},
		cfg.GitHub.Index,
	)
}
