FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o /app/api-client ./cmd/api-client

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/api-client /api-client

ENTRYPOINT ["/api-client"]
