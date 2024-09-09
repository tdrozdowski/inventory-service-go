# inventory-service-go
![Current Build](https://github.com/tdrozdowski/inventory-service-go/actions/workflows/go.yml/badge.svg)
[![codecov](https://codecov.io/gh/tdrozdowski/inventory-service-go/graph/badge.svg?token=JNAPPBO1OV)](https://codecov.io/gh/tdrozdowski/inventory-service-go)

Example of inventory service in Go, using SQLx and Echo.

NOTE - this does assume you ran the migrations or have a db from the companion Rust project.  See [here](https://github.com/tdrozdowski/inventory-service-rs?tab=readme-ov-file#getting-started) for more details.
The supporting docker-compose file is included with this project to run the database.

## Getting Started
This project builds using standard Go tookit tools - nothing extra is needed.

On MacOS - use Homebrew:
```bash
brew install go
```

To start up the server just do:
```bash
go build
```

And then you can use the http client scripts in the `client-http` folder or use Postman.

To access the OpenAPI docs go here: http://localhost:8080/docs

## Development
If you care to extend/modify this code, you will need to install a couple tools.
First off, Google Wire is used to manage Dependency Injection.  Uber's Mock lib is also used to manage the mocks.
Both of these have command line tools to generate code.

Be sure to install the following:

Google Wire
```bash
go install github.com/google/wire/cmd/wire@latest
```

Uber Mock
```bash
go install go.uber.org/mock/mockgen@latest
```

Now you can regen the wires or mocks as needed.