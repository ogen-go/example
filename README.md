# example

Example for server and client with graceful shutdown and open telemetry exporters for
prometheus and jaeger.


* Force `go.mod` to have ogen in [tools.go](./tools.go) 
* Generate code by [gen.go](./gen.go)
* Metrics boilerplate in [internal/app](./internal/app)
* Advanced instrumentation in [internal/httpmiddleware](./internal/httpmiddleware)

```bash
docker compose up
```
