FROM golang:alpine AS build

RUN apk add upx

COPY . /app
WORKDIR /app/cmd/server
RUN go get && go build -ldflags "-s -w" -o /server
RUN upx /server

# ---

FROM alpine
COPY --from=build /server /server
ENTRYPOINT ["/server"]
