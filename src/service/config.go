package service

import (
	"fmt"
	"slices"

	env "github.com/caarlos0/env/v11"
	"github.com/roboslone/github-oauth-exchange/src/github"
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
	Address        string   `env:"ADDRESS" envDefault:":8080"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
}

type GitHub struct {
	RawApplications []string `env:"APPLICATIONS"`
	Index           map[string]github.Application
}

func NewConfig(opts ...InitOption) (*Config, error) {
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

	if len(cfg.GitHub.RawApplications) == 0 {
		return nil, fmt.Errorf("at least one GitHub application is required (GITHUB__APPLICATIONS)")
	}

	cfg.GitHub.Index = make(map[string]github.Application, len(cfg.GitHub.RawApplications))
	for n, s := range cfg.GitHub.RawApplications {
		app, err := github.ApplicationFromString(s)
		if err != nil {
			return nil, fmt.Errorf("parsing application #%d: %q: %w", n+1, s, err)
		}
		cfg.GitHub.Index[app.ClientID] = app
	}

	slices.Sort(cfg.Server.AllowedOrigins)

	return &cfg, nil
}
