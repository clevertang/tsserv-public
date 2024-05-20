
---

# Very Trivial Time Series Server

This is a time series server that outputs data from a single deterministic time series.

## Table of Contents
- [Overview](#overview)
- [Installation](#installation)
- [Usage](#usage)
- [Development](#development)
- [License](#license)

## Overview
`tsserv-public` is a simple time series data server designed to provide deterministic time series output.

## Installation
To install the server, clone the repository and build the project:

```sh
git clone -b dev https://github.com/clevertang/tsserv-public.git
cd tsserv-public
go build ./cmd/server
```

## Usage
After building the project, you can run the server using:

```sh
./server
```

## Development
### Running Tests
To run the tests, use the following command:

```sh
go test ./...
```

### Adding Unit Tests
Unit tests can be added in the `_test.go` files. Ensure each test covers a specific functionality of the codebase.


## License
there's no license for this project

---
