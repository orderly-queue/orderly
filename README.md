# Go Template

A template repo for golang APIs, comes with:

- YAML configs
- Postgres db connection
- Atlas/golang-migrate migrations
- Metrics server
- Probes server
- HTTP server
- Test containers setup
- SQLC
- OpenTelemetry Tracing
- Sentry error tracking
- User creation/JWT authentication

## Get Started

First, create a new repo from the template, then run:

```
./hack/rename.sh
```

and follow the promts to rename the go module etc.

## Running Tests

On every PR, the Dockerfile will be built and unit tests will be run, you can run these manually with:

```
task build
```

```
task test:unit
```
