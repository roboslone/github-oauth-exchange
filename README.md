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
docker run -p 8080:8080 -e "GITHUB__APPLICATIONS=..." ghcr.io/roboslone/github-oauth-exchange:main
```

Codes then can be exchanged like so:

```sh
curl http://localhost:8080/github.v1.ExchangeService/Exchange \
    -H 'content-type: application/json' \
    -d '{"client_id": "...", "code": "..."}'
# {"accessToken":{"value":"***","expiresIn":"28800s"},"refreshToken":{"value":"***","expiresIn":"15724800s"},"tokenType":"bearer"}
```
