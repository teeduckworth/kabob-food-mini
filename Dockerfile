# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/kabobfood ./cmd/app

FROM gcr.io/distroless/base-debian12 AS app

ENV APP_ENV=prod \
    HTTP_HOST=0.0.0.0 \
    HTTP_PORT=8080

COPY --from=builder /build/kabobfood /bin/kabobfood

EXPOSE 8080

ENTRYPOINT ["/bin/kabobfood"]
