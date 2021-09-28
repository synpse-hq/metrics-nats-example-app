# CPU and Memory gathering example application with NATS

Example application which gathers `CPU` and `Memory` stats (see `pkg/metricsbridge/config/config.yaml` for all metrics) 
and sending them to NATS (topic `metrics`) running on the same device. Application from first sight might look bit complicated (but it is not). 
This is so it represent more real life scenario. Where we have API server so external entities to the device could interact with the application,
metrics gathering/backend process, messaging for async communication. This should be representing real world application usecase.

All packages are explained bellow.

* `cmd` - entrypoint for execution
* `agent` - package level entrypoint, where we initiate all the services and run them as go routines
* `api` - all shared types and structs
* `metrics` - metrics collection package. We gather metrics and set them as prometheus exporters (optional in real world)
* `metricsbridge` - application internally accessing prometheus metrics set by `metrics` package and on timely basis sending them to NATS queue.
* `service` - (Optional) API layer of the application. Currently exposing single `/metrics` endpoint with prometheus metrics. 


# Build

To build and push all images with your own name:

```
export APP_REPO=quay.io/example/repo
make push-all
```