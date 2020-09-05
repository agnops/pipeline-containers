FROM golang:1.13-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/app

COPY helm-deploy/ .

RUN apk add build-base

RUN go mod download

RUN go build -o ./out/helm-deploy .

RUN ls -alt ./out/helm-deploy


FROM alpine:3.9

COPY --from=build_base /tmp/app/out/helm-deploy /usr/bin/helm-deploy