# === INSTALL CERTS STAGE === #
FROM alpine:latest as certs
RUN apk --update add ca-certificates

# === BUILD STAGE === #
FROM golang:1.12-alpine as build

RUN apk add --no-cache git

WORKDIR /srv/app
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# RUN go test -v ./...
RUN go build -ldflags="-w -s" -o build

# === RUN STAGE === #
FROM scratch as run

WORKDIR /srv/app
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /srv/app/build /srv/app/build

ENTRYPOINT ["/srv/app/build"]
