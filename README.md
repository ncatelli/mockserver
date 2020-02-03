# mockserver
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
                - [URL Variables](#url-variables)
        - [Drivers](#drivers)
            - [Simple](#simple)
                - [Example](#example)

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

- ADDR:        `string` The server address mockserver binds to.
- CONFIG_PATH: `string` A filesystem path to the simple driver config file.

### Response Bodies
All response bodies in for handlers are valid [go templates](https://golang.org/pkg/html/template/). In addition some helper data is included in each template variable to be referenced for rendering. This includes the following:

- Template Parameters
- GTF Functions

#### Template Parameters
Template parameters are passed directly into the template via the [data argument at time of execution](https://golang.org/pkg/html/template/#Template.Execute) and include both the `http.Request` object for the request that generated the template as well as any path variables that are parsed from the request URL.

##### URL Variables
The mockserver allow for the parsing of variables directly out of a url path through the [gorilla/mux router](http://www.gorillatoolkit.org/pkg/mux#Vars) and more information on what kind of pattern matching can be accomplished by the router can be found at the preceeding link.

### Drivers
#### Simple
##### Example
```yaml
---
- path: "/test/pathvar/{embed}"
  method: GET
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
    static_response: '{"resp": "Ok"}'
    response_status: 200
  - weight: 1
    response_headers:
      content-type: text/plain
    static_response: ''
    response_status: 500
```