# libs

Is a collection of reusable configurations or libraries that can be used across different services.
Here we can define a standard way for change/configure something across the whole infra structure, specially things
that might affect the infra, metrics or engineering perspective of things.

## Rules to create a good library:

	- must be indpendent from each others, to receive any format of other library you must use interface.
	- must be configurables with options, and runs with minimal to no configuration as base.
	- must be fully backward compatible.
	- if the library is helping provide a communication layer,
			it must be able to receive & send traces according to [standards docs]()

## Supported Libs

- [x] logger: responsible on providing the logger with base 6 functions.
- [x] config reader: responsible to provide a reader of any struct and pointer of structs using tag `mapstructure`.
- [x] HTTP: Provide HTTP Server with base configurations, that suits our own.
- [x] PPROF server: Provide HTTP server of pprof that can be used for debugging golang.
- [x] otel: responsible to provide TraceProvider and tracer that can be used to uasily with all server types.
- [ ] database: responsible on creating the db connections using the best practices.
- [ ] Circuit breaker
- [ ] metric: metric provider that directly connect to our metric collectors.
- [ ] grpc
- [ ] pubSub
- [ ] Broker
- [ ] Health Checker
- [ ] inmemory cache
