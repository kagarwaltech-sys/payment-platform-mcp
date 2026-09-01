FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/payment-platform-mcp ./cmd/server

FROM alpine:3.22

RUN adduser -D -H -u 10001 app
COPY --from=build /out/payment-platform-mcp /usr/local/bin/payment-platform-mcp
USER app
ENTRYPOINT ["/usr/local/bin/payment-platform-mcp"]
