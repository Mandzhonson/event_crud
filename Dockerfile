FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/main.go

FROM alpine

COPY --from=builder ./app/app ./
COPY  .env .

EXPOSE 8080
CMD ["./app"]

