# Builder stage
FROM golang:alpine AS builder
WORKDIR /build

## Install dependencies
RUN apk update && apk add --no-cache build-base git
COPY go.mod go.sum ./
RUN go mod download

## Build application
COPY main.go ./
COPY cmd ./cmd
COPY config ./config
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go build -tags static_all,musl -o main .

# Runner stage
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

COPY --from=builder /build/main .
COPY --from=builder /build/config ./config

CMD ["./main"]  
