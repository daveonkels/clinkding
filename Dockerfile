FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o clinkding .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /build/clinkding /usr/local/bin/clinkding

ENTRYPOINT ["clinkding"]
