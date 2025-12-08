# GitHub code exchange

This is a simple server capable of exchanging GitHub authorization codes for access tokens.  
`GITHUB__APPLICATIONS` environment variable is used to configure supported applications:

```sh
export GITHUB__APPLICATIONS='<name>:<client_id>:<client_secret>
```

Server address is configured via `SERVER__ADDRESS` and defaults to `:8080`.

Run using Docker image:

```sh
docker run -p 8080:8080 -e "GITHUB__APPLICATIONS=..." ghcr.io/roboslone/github-oauth-exchange:main
```

Codes then can be exchanged like so:

```sh
curl http://localhost:8080/github.v1.ExchangeService/Exchange \
    -H 'content-type: application/json' \
    -d '{"app_name": "into-the-void", "code": "..."}'
# {"accessToken":{"value":"***","expiresIn":"28800s"},"refreshToken":{"value":"***","expiresIn":"15724800s"},"tokenType":"bearer"}
```
