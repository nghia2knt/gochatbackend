FROM golang:1.19-alpine AS build

RUN apk add alpine-sdk bash

WORKDIR /app

COPY . .

COPY vendor/ ./vendor/

RUN  GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -tags dynamic --mod=vendor -o main ./cmd/conversation

FROM alpine:3.14

COPY --from=build /app/main /usr/bin/conversation-service

LABEL name=conversation-service

ENTRYPOINT ["conversation-service"]
