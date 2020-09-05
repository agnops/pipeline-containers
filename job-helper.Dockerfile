FROM golang:1.15-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/app

COPY job-helper .

RUN go mod download

RUN go build -o ./out/job-helper .

RUN ls -alt ./out/job-helper

FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=build_base /tmp/app/out/job-helper /app/job-helper

CMD ["/app/job-helper"]