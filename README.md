# GitHub code exchange

This is a simple server capable of exchanging GitHub authorization codes for access tokens.  
`GITHUB__APPLICATIONS` environment variable is used to configure supported applications:

```sh
export GITHUB__APPLICATIONS='<client_id_1>:<client_secret_1>,<client_id_2>:<client_secret_2>'
```

Server address is configured via `SERVER__ADDRESS` and defaults to `:8080`.

CORS is configured with `SERVER__ALLOWED_ORIGINS`:

```sh
export SERVER__ALLOWED_ORIGINS='example.com,whatever.org'
```

Run using Docker image:

```sh
docker run \
    -d \                                            # run in background
    -p 11275:8080 \                                 # expose port 11275 to host
    -e 'GITHUB__APPLICATIONS=...' \                 # github apps
    -e 'SERVER__ALLOWED_ORIGINS=...' \              # allowed origins for CORS 
    ghcr.io/roboslone/github-oauth-exchange:main
```

Codes then can be exchanged like so:

```sh
curl http://localhost:11275/github.v1.ExchangeService/Exchange \
    -H 'content-type: application/json' \
    -d '{"client_id": "...", "code": "..."}'
# {"accessToken":{"value":"***","expiresIn":"28800s"},"refreshToken":{"value":"***","expiresIn":"15724800s"},"tokenType":"bearer"}
```

Server then can be exposed to the internet via TLS-terminating proxy, e.g. nginx.  
Don't expose unencrypted HTTP server, this will leak GitHub tokens.
