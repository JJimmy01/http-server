# syntax=docker/dockerfile:1

##
## build
##
FROM golang:1.17-buster AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o /http-server

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
ENV VERSION v1.0.0

WORKDIR /
COPY --from=build /http-server /http-server

EXPOSE 8080

ENTRYPOINT ["/http-server"]