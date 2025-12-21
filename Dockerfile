FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_DATE=unknown

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X 'main.Version=${VERSION}' -X 'main.Commit=${COMMIT}' -X 'main.BuildDate=${BUILD_DATE}'" -o ./proxygate ./cmd/server/

FROM alpine:3.22.0

RUN apk add --no-cache tzdata

RUN mkdir /proxygate

COPY --from=builder ./app/proxygate /proxygate/proxygate

COPY config.example.yaml /proxygate/config.example.yaml

WORKDIR /proxygate

EXPOSE 8317

ENV TZ=UTC

CMD ["./proxygate"]