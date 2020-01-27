# mock-server
## General
This tool provides a simple framework for generating performance test-ready mocks from configuration files.

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

- ADDR: `string` The server address mockserver binds to.
