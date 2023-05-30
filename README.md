# StartupBuilder
is a general project of trying to apply microservices practices (Doesn't have to be the best)

## On Boarding

The project will try to create basic needs of any startups, in addition on creating an operation framework of startup.
the general project is created in a monorepo pattern, which means the whole startup is here, but also means they are all
isolated by and act independently (keep in mind that you might need to breakdown into separate repos, always be ready);
below is expalaination of each folder:
```bash
apps/				# <- Should include each service (defined responsibility) each app has it's own language.
						# the important thing is that LSP can work within that directory as if it's separate repo.
						# it might include specific tools needed to run the app. each app must be containerized and providing deployment
						# guideline
pkg/ # collection of independent libraries; this just to enable easy importing at the begging of any project, JS/TS monorepo can take the advantage of this by configuration.
		# goLang based projects can use the go package system the same way by init each folder.
		# Rust can also benefit from Cargo system.
		# so in general having all libraries in one folder is just a metaphor for now as if those
		# are distributed across different repos.
deployments/ # this app should have all deployments related files that are not tools, example: Dockerfile,
						# extra compose files, k8s files etc..
tools/      # helper tools like extra Makefiles, setups ... to add more functions to environment runner, or test runner.
						# the purpose here is tools might be written in other language, like python app create environment files
						# most of the time you don't need much of tools, except maybe add extra compose files for testing
						# instrumentation that won't be needed in deployments, or maybe needed for local runners only.
test/       # any test helper functions to be used within the code (should be the same language as app).
db/
 |_migrations/*.sql  # migration files should always be in SQL format only
 |_initdb.d/         # might be needed to control infra related changes to the database on local
										# this changes should be applied MANUALLY only on production environments
										# for security
scripts/            # Should be any script that can as helper script in running the application related only. example (run tests in specific way, run codecoverage) but if you are creating codecoverage tool you should add it to /tools.
templates/ # apps barebone structure templates to speed up app creation. (this can be ignored in general).
inspire/ # folder with submodules only to get inspiration from projects with LSP working.
```
## Setup

### Local configurations

* add the following domains to `/etc/hosts`
    ```hosts
    127.0.0.1 permify.localhost
    127.0.0.1 traefik.localhost
    127.0.0.1 jaeger.localhost
    127.0.0.1 prometheus.localhost
    127.0.0.1 admin.localhost
    ```

* you can change this URLs as well in `.env` file; check `.env.sample` for more information.
    ```env
    ## URLS
    BASE_URL=localhost
    TRAEFIK_URL=traefik.${BASE_URL}
    PROMETHEUS_URL=prometheus.${BASE_URL}
    JAEGER_URL=jaeger.${BASE_URL}
    PERMIFY_URL=permify.${BASE_URL}
    ADMIN_URL=admin.${BASE_URL}
    ```
## Checklist

- Company Level
    - Organise Docs structure.
- Backend
	- OTEL
		- add HTTP Client instrumentations. [net/HTTP](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/net/http)
		- add HTTP Server instrumentations. [otelmux](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/github.com/gorilla/mux/otelmux)
		- add gRPC unray Interceptor. [otelGrpc](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/google.golang.org/grpc/otelgrpc)
		- Server CPU and memory [runtime](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/runtime)
	- gRPC server runner.
	- database driver creator.
	- circuit-breaker
	- metric
    - http
        - rate limiter
        - CORS
	- authentication service
		- Should manage JWTs
		- Should manage user sessions
		- should manage save user devices
		- Should provide pkg `auth` as:
			- gRPC inerceptor
			- HTTP middleware
	- identity service
		- Should manage users
		- Should manage permission
	- notification service
		- Should manage user preferred notification channels.
		- Should manage user devices registry.
		- Should Send user notification according to preferred channels.
	- Subscription service
		- Should manage user Subscription topics
		- Should handle marketing events to send user according to subscribed topics.
		- Should do marketing according to user notification channel choice.
	- Product Service
		- should manage product {description, image, tags}
	- Insights Service
		- Should manage user tracable actions
	- Make tracer work with HTTP headers

- infrastructure
	- skaffold
	- grafana
	- jeager
	- metric
	- pubsub (I.E: Kafka)
	- devHub (I.E: spotify backoffice)

## Inspiring Projects

### Authentication

- [GoTrue](https://github.com/netlify/gotrue) An SWT based API for managing users and issuing SWT tokens.
- [SupaBase Auth](https://github.com/supabase/auth) A JWT based API for managing users and issuing JWT tokens.
- [Golang SCIM](https://github.com/elimity-com/scim) Golang Implementation of the SCIM v2 Specification.
- [Django SCIM2](https://github.com/15five/django-scim2) Django implementation of SCIM v2.
- [go-auth](https://github.com/Sirneij/go-auth) A fullstack session-based authentication system using golang and sveltekit.
- [keycloak](https://github.com/keycloak/keycloak) Java Based
- [Authelia](https://github.com/authelia/authelia) Golang to provide only authentication.
- [authentik](https://github.com/goauthentik/authentik) (Python, Golang Based)
- [Hanko](https://github.com/teamhanko/hanko) (Go Based)

