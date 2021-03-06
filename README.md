# mockserver
![Build tests](https://github.com/ncatelli/mockserver/workflows/Build%20tests/badge.svg?branch=master)

## General
This tool provides a simple framework for generating performance test-ready mocks from configuration files.

## Table of Contents
<!-- TOC -->

- [mockserver](#mockserver)
    - [General](#general)
    - [Table of Contents](#table-of-contents)
    - [Dependencies](#dependencies)
    - [Building](#building)
        - [Docker](#docker)
        - [Locally](#locally)
    - [Testing](#testing)
        - [Locally](#locally-1)
    - [Configuration](#configuration)
        - [Services](#services)
        - [Response Bodies](#response-bodies)
            - [Template Parameters](#template-parameters)
                - [Path Variables](#path-variables)
            - [Template Functions](#template-functions)
                - [GTF Functions](#gtf-functions)
                - [Custom Generators](#custom-generators)
        - [Drivers](#drivers)
            - [yaml](#yaml)
                - [Parameters](#parameters)
                    - [path](#path)
                    - [method](#method)
                    - [middleware](#middleware)
                - [request_headers](#request_headers)
                - [query_params](#query_params)
                    - [Handlers](#handlers)
                - [Example](#example)
        - [Middlewares](#middlewares)
            - [logging](#logging)
                - [settings](#settings)
            - [latency](#latency)
                - [settings](#settings-1)

<!-- /TOC -->

## Dependencies
- make

## Building
### Docker
The tool can be built and run entirely via docker using the following command.

```sh
$> docker build -t ncatelli/mockserver .
```

### Locally
The tool can also be built and installed locally by running a pip install from the root of the project.

```
$> make build
```

## Testing
Tests can be run using the built in go testing library and a convenient wrapper to test all subpackages has been provided below.

### Locally
Local tests default to running tests on all subpackages along with coverage tests.
Tests can be run with the following make command.

```
$> make test
```

## Configuration
### Services
The mockserver service can be configured via the following environment variables:

- ADDR:        `string`  The server address mockserver binds to.
- CONFIG_PATH: `string`  A filesystem path to the simple driver config file.
- CONFIG_URL:  `url.URL` A URL path to fetch the configuration file from. This
    is useful for when a service wants to publish its own configuration file.

It's worth noting that _EITHER_ `CONFIG_PATH` or `CONFIG_URL` should be sent. If both are set, `CONFIG_PATH` takes priority.

### Response Bodies
All response bodies in for handlers are valid [go templates](https://golang.org/pkg/html/template/). In addition some helper data is included in each template variable to be referenced for rendering. This includes the following:

- Template Parameters
- GTF Functions

#### Template Parameters
Template parameters are passed directly into the template via the [data argument at time of execution](https://golang.org/pkg/html/template/#Template.Execute) and include both the `http.Request` object for the request that generated the template as well as any path variables that are parsed from the request URL.

##### Path Variables
The mockserver allow for the parsing of variables directly out of a url path through the [gorilla/mux router](http://www.gorillatoolkit.org/pkg/mux#Vars) and more information on what kind of pattern matching can be accomplished by the router can be found at the preceeding link.

#### Template Functions
Mocking functionality is implatemented via golang's [stdlib template functions](https://golang.org/pkg/html/template/#FuncMap). A few additional libraries and features have been included to aid in extending this functionality.

##### GTF Functions
[GTF](https://github.com/leekchan/gtf) is a template function library with the stated goal of implementing the functions included in jinja2. Further documentation on the functions included can be found at their github page.

##### Custom Generators
Custom generators provides a simple way to add new functions that will be compiled into the mockserver at build time.

New generators can be added by:

- Creating a new package under github.com/ncatelli/mockserver/pkg/router/generator/plugins
  - This package __MUST__ include a `Generator` type.
  - This new `Generator` __MUST__ satisfy the interface `github.com/ncatelli/mockserver/pkg/router/generator.Generator`
- Rerun `make` to generate the correct package imports and build mockserver with the new plugin.

### Drivers
#### yaml
The yaml driver implements a simple configuration format that maps directly to the implementation of the Route struct.
##### Parameters
###### path
**Required**

This field represents a url path to be passed to the router and supports all [gorilla path matching and variables](https://github.com/gorilla/mux#matching-routes). All variables specified in the path are passed back to the handlers via [path variables](#path-variables).

###### method
**Required**
The HTTP that this route will match against. This field currently only matches 1 method.

###### middleware
This field takes a map of logging drivers and a map of strings to be passed in for configuring the middlewares. Further information on the available middleware and their configuration parameters and their settings can be found in the [middlewares section](#middlewares).

##### request_headers
This field represents a key-value mapping of headers that must be defined to be routeable to the defined route.

##### query_params
This field represents a key-value mapping of query parameters that must be set to be routable to the defined route.

###### Handlers
The handlers field takes a weighted list of objects that map directly to the Handler structure. Subfields of handlers represent

- weight: A positive weighted value to determine the frequency a handler is hit. Higher represents more frequent hits. Zero represents unrouteable (good for a temporarily disabled handler).
- response_headers: A key-value store of additional headers to be attached to the response body.
- static_response: A response body template to respond with. This supercedes the response_path setting and is suitable for short responses.
- response_path: A file path to a file that will be used to generate the response body. This is more suitable for multi-line responses that will be difficult to fit into a static_response.
- response_status: A status code to assign to the response.

##### Example
```yaml
---
- path: "/test/pathvar/{embed}"
  method: GET
  middleware:
    logging:
      target: stdout
  handlers:
  - weight: 1
    response_headers:
      content-type: text/plain
    static_response: '{{ .PathVars.embed }}'
    response_status: 200
- path: "/test/weighted"
  method: GET
  handlers:
  - weight: 2
    response_headers:
      content-type: application/json
    response_path: /examples/example_response_body.txt
    response_status: 200
  - weight: 1
    response_headers:
      content-type: text/plain
    static_response: ''
    response_status: 500
- path: "/test/with/required/headers"
  method: GET
  request_headers:
    status: ok
  handlers:
  - weight: 1
    response_headers:
      content-type: text/plain
    static_response: 'ok'
    response_status: 200
- path: "/test/with/required/query/params"
  method: GET
  query_params:
    status: ok
  handlers:
  - weight: 1
    response_headers:
      content-type: text/plain
    static_response: 'ok'
    response_status: 200
```

### Middlewares
#### logging
The logging handler implements the [gorilla logging handler](https://godoc.org/github.com/gorilla/handlers#LoggingHandler) and outputs logs to a target in Apache CLF format.

##### settings
target (default: `stdout`): a target to write files to. Currently this only supports stdout.

#### latency
The latency middleware allows injection of artificial latency into a route to mimic either transit or processing time. This latency can be specified either as a static value or as a range of time.

##### settings
latency (default: `0`): A static latency in milliseconds to inject into a response.
min     (default: `0`): A minimum value for a range of latency in a response.
max     (default: `0`): A maximum value for a range of latency in a response.
