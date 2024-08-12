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

## Get Started

First, create a new repo from the template, then run:

```
task rename
```

and follow the promts to rename the go module etc.

## Running Tests

On every PR, the Dockerfile will be built and unit tests will be run, you can run these manually with:

```
# Builds the docker image
task build
```

```
task test:unit
```
