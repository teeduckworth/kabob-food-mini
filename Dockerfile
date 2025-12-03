# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/kabobfood ./cmd/app

FROM gcr.io/distroless/base-debian12

ENV APP_ENV=prod \
    HTTP_HOST=0.0.0.0 \
    HTTP_PORT=8080

COPY --from=builder /bin/kabobfood /bin/kabobfood

EXPOSE 8080

ENTRYPOINT ["/bin/kabobfood"]
