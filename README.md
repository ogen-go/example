# example

Example for generated server and client, instrumented with OpenTelemetry.

* Force `go.mod` to have ogen in [tools.go](./tools.go) 
* Generate code by [gen.go](./gen.go)
* Metrics boilerplate in [internal/app](./internal/app)
* Advanced instrumentation in [internal/httpmiddleware](./internal/httpmiddleware)

```bash
docker compose up
```
