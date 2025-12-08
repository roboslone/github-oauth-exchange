package config

import (
	"fmt"
	"strings"

	env "github.com/caarlos0/env/v11"
)

type InitOptions struct {
	Env *env.Options
}

type InitOption func(*InitOptions)

func WithEnvOptions(o *env.Options) InitOption {
	return func(io *InitOptions) {
		io.Env = o
	}
}

type Config struct {
	Server Server `envPrefix:"SERVER__"`
	GitHub GitHub `envPrefix:"GITHUB__"`
}

type Server struct {
	Address string `env:"ADDRESS" envDefault:":8080"`
}

type GitHubApplication struct {
	Name         string
	ClientID     string
	ClientSecret string
}

type GitHub struct {
	Applications []string `env:"APPLICATIONS"`
	Index        map[string]GitHubApplication
}

func New(opts ...InitOption) (*Config, error) {
	o := &InitOptions{}
	for _, opt := range opts {
		opt(o)
	}

	var cfg Config
	var err error
	if o.Env != nil {
		cfg, err = env.ParseAsWithOptions[Config](*o.Env)
	} else {
		cfg, err = env.ParseAs[Config]()
	}
	if err != nil {
		return nil, err
	}

	if len(cfg.GitHub.Applications) == 0 {
		return nil, fmt.Errorf("at least one GitHub application is required (GITHUB__APPLICATIONS)")
	}

	cfg.GitHub.Index = make(map[string]GitHubApplication, len(cfg.GitHub.Applications))
	for n, s := range cfg.GitHub.Applications {
		app, err := parseApp(s)
		if err != nil {
			return nil, fmt.Errorf("parsing application #%d: %q: %w", n+1, s, err)
		}
		cfg.GitHub.Index[app.Name] = app
	}

	return &cfg, nil
}

func parseApp(s string) (GitHubApplication, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return GitHubApplication{}, fmt.Errorf("invalid format, expected {name}:{client_id}:{client_secret}")
	}

	app := GitHubApplication{
		Name:         strings.TrimSpace(parts[0]),
		ClientID:     strings.TrimSpace(parts[1]),
		ClientSecret: strings.TrimSpace(parts[2]),
	}

	if app.Name == "" {
		return app, fmt.Errorf("name is empty")
	}
	if app.ClientID == "" {
		return app, fmt.Errorf("client_id is empty")
	}
	if app.ClientSecret == "" {
		return app, fmt.Errorf("client_secret is empty")
	}

	return app, nil
}
