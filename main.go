package main

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"net/http"
	"slices"

	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"connectrpc.com/validate"
	"github.com/roboslone/github-oauth-exchange-proto/github/v1/githubv1connect"
	"github.com/roboslone/github-oauth-exchange/src/config"
	"github.com/roboslone/github-oauth-exchange/src/service"
	"github.com/rs/cors"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("available applications (%d):\n", len(cfg.GitHub.Index))
	for _, id := range slices.Sorted(maps.Keys(cfg.GitHub.Index)) {
		fmt.Printf("\t%q\n", id)
	}

	mux := http.NewServeMux()

	path, handler := githubv1connect.NewExchangeServiceHandler(
		service.New(cfg),
		connect.WithInterceptors(
			validate.NewInterceptor(),
		),
	)
	mux.Handle(path, addCORS(cfg, handler))

	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)

	s := http.Server{
		Addr:      cfg.Server.Address,
		Handler:   mux,
		Protocols: p,
	}

	fmt.Printf("listening on %q\n", cfg.Server.Address)
	err = s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func addCORS(cfg *config.Config, handler http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: cfg.Server.AllowedOrigins,
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	}).Handler(handler)
}
