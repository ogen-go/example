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