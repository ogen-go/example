FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o /app/api-server ./cmd/api-server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/api-server /api-server

ENTRYPOINT ["/api-server"]
