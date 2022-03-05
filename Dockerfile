# --- Base ----
FROM golang:1.17.8-stretch AS base
WORKDIR $GOPATH/src/github.com/esequielvirtuoso/oauth_go_lib

# ---- Dependencies ----
FROM base AS dependencies
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN	go mod tidy
RUN	go mod vendor

# ---- Test ----
FROM dependencies AS test
COPY . .
RUN go get -u github.com/axw/gocov/gocov
RUN GO111MODULE=off go get -u github.com/matm/gocov-html
ARG POSTGRES_URL
RUN go test -v -cpu 1 -failfast -coverprofile=coverage.out -covermode=set ./...
RUN gocov convert coverage.out | gocov-html > /index.html
RUN grep -v "_mock" coverage.out >> filtered_coverage.out
RUN go tool cover -func filtered_coverage.out
