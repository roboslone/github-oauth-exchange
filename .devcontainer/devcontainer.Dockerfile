FROM node:24-trixie-slim

ENV DEVCONTAINER="true" \
    CGO_ENABLED="0" \
    PATH="/root/.local/bin:${PATH}"

RUN apt update -y && apt install -y --no-install-recommends ssh
