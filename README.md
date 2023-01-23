# Go metrics

library to capture go service metrics, abstracting protocol (http, gRPC) and details of timeseries db.


## [Setup](Setup)

Expose prometheus `/metrics` endpoint in `9101` port

```go
obs, err := metrics.Setup(
	    metrics.WithAddress(":9101"),
        metrics.WithServiceName("service-name")
    )
m := mux.NewRouter()
m.Use(mux.MiddlewareFunc(obs.Middleware))

# add required handlers ... 
err = http.ListenAndServe(addr, m)
```

Once integrated prometheus can be configured to scrape `localhost:9101` and you could setup a dashboard to monitor your service.

![service dashboard](./assets/service_dashboard.png)


## Architecture 

In order to decide which means to use, you've to be aware of Pull vs Poll architecture and what fits your need based on ecosystem.

![pull vs poll](./assets/pull_vs_poll_metrics.png)


In k8s ecosystem prometheus is a standard and in vm ecosystem influx is used widely.


## TODO

We'll support prometheus, influxdb and add open-telemetry for db transactions

- [x] Prometheus
- [x] Pprof
- [ ] open-telemetry
- [ ] gRPC interceptor
- [ ] sql spans tracing
- [ ] WIP - Influx middleware

