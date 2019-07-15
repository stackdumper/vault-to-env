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
FROM alpine as run

WORKDIR /srv/app
COPY --from=build /srv/app/build /usr/local/bin/vte

ENTRYPOINT ["/bin/sh", "-c"]
