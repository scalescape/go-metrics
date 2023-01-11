# Go metrics

library to capture go service metrics, abstracting protocol (http, gRPC) and details of timeseries db.

```
obs, err := metrics.Setup(
	    metrics.WithAddress(":9101"),
        metrics.WithServiceName("service-name")
    )
m := mux.NewRouter()
m.Use(mux.MiddlewareFunc(obs.Middleware))

# add required handlers ... 
err = http.ListenAndServe(addr, m)
```

We'll support prometheus, influxdb and add open-telemetry for db transactions

## TODO
- [ ] open-telemetry
- [ ] gRPC interceptor
- [ ] sql spans tracing

