package github

import (
	"fmt"
	"strings"
)

type Application struct {
	ClientID     string
	ClientSecret string
}

func ApplicationFromString(s string) (Application, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return Application{}, fmt.Errorf("invalid format, expected {client_id}:{client_secret}")
	}

	app := Application{
		ClientID:     strings.TrimSpace(parts[0]),
		ClientSecret: strings.TrimSpace(parts[1]),
	}

	if app.ClientID == "" {
		return app, fmt.Errorf("client_id is empty")
	}
	if app.ClientSecret == "" {
		return app, fmt.Errorf("client_secret is empty")
	}

	return app, nil
}
