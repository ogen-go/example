# example

Example for server and client with graceful shutdown and open telemetry exporters for
prometheus and jaeger.


* Force `go.mod` to have ogen in [tools.go](./tools.go) 
* Generate code by [gen.go](./gen.go)
* Metrics boilerplate in [internal/app](./internal/app)

## Start server
```go
go run ./cmd/api-server
```

## Check metrics
```console
$ curl localhost:9090
Service is up and running.

Resource:
  service.name                     api
  service.namespace                example
  telemetry.sdk.language           go
  telemetry.sdk.name               opentelemetry
  telemetry.sdk.version            1.8.0

Available debug endpoints:
/metrics             - prometheus metrics
/debug/pprof         - exported pprof
```

## Use client
```console
$ go run ./cmd/api-client --id 1337
pet: {
  "id": 1337,
  "name": "Pet 1337",
  "status": "available"
}
```

## Environmental variables

| Name                              | Value                       | Description                                  |
|-----------------------------------|-----------------------------|----------------------------------------------|
| `OTEL_SERVICE_NAME`               | `api`                       | OpenTelemetry service name                   |
| `OTEL_RESOURCE_ATTRIBUTES`        | `service.namespace=example` | Additional OpenTelemetry resource attributes |
| `OTEL_EXPORTER_JAEGER_AGENT_HOST` | `localhost`                 | Jaeger host to use                           |
| `OTEL_EXPORTER_JAEGER_AGENT_PORT` | `6831`                      | Jaeger port to use (UDP)                     |

## TODO

- [ ] Add server context propagation example
- [ ] Add client context propagation example