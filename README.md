# example

Example for generated server and client, instrumented with OpenTelemetry.

* Force `go.mod` to have ogen in [tools.go](./tools.go) 
* Generate code by [gen.go](./gen.go)
* Metrics boilerplate in [internal/app](./internal/app)
* Advanced instrumentation in [internal/httpmiddleware](./internal/httpmiddleware)

```bash
docker compose up
```

You can open Grafana dashboard on http://localhost:3000 to observe telemetry.
For example, you can see client traces in [TraceQL explore][traces].

[traces]: http://localhost:3000/explore?orgId=1&left=%7B%22datasource%22:%22tempo-oteldb%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22tempo%22,%22uid%22:%22tempo-oteldb%22%7D,%22queryType%22:%22nativeSearch%22,%22limit%22:20,%22serviceName%22:%22client%22%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D