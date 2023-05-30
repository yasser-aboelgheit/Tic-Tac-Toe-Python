# bc-go-tracer

This package should provide the basic tracer according to standard lib of `opentracing`; Building your service should 
provide tracing for every action/function in the a way we can detrmine the time it took from `jaeger`;

Thanks to `go-kit` for inspiration.

How we should do tracing:
- must exist tracing attributes: service, version, env, resource (function name)
- may exist when possible: 
  ```go
	type customer struct {
		id string
		device struct{
						id string
						version string
					} 
		session struct{
						id string
					}
	}
	```
- if you are going to call another service you must add all above in addition to `trace_id`, last `span_id`; 
	later we can detrmine the communications between services.
- the library must provide all necessary functions to this out of the box:
	* Inject to logger
	* grpc Middleware:
		= to extract keys from context to `RPC.Meta`.
		= from RPC.Meta inject to context.
	* http Middleware:
		= extract keys from headers to context.
		= inject keys to context from headers.

