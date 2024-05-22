### Very Trivial Time Series Server

This repository provides a simple HTTP server that generates and serves deterministic time series data. The server is
implemented in Go and contains endpoints for retrieving data points in a specified time range.

#### Features

- **/hello Endpoint:** Responds with a simple greeting message.
- **/data Endpoint:** Returns time series data points within a specified range.

#### Project Structure

- **cmd/tsserv:** Contains the main entry point for the server.
- **pkg/tsserv:** Implements the server logic, including HTTP handlers.
- **pkg/datasource:** Generates deterministic time series data.
- **pkg/logger:** Provides logging utilities.

#### Getting Started

1. Clone the repository:
   ```sh
   git clone https://github.com/clevertang/tsserv-public.git
   cd tsserv-public
   ```

2. Build and run the server:
   ```sh
   go build -o tsserv ./cmd/server
   ./tsserv -p 8080
   ```

3. Access the endpoints:
    - Greeting: `http://localhost:8080/hello`
    - Time Series Data: `http://localhost:8080/data?begin=<RFC3339_start>&end=<RFC3339_end>`

#### Code Overview

- **main.go:**
    - Sets up and starts the server.
    - Implements graceful shutdown.

- **server.go:**
    - Defines the `Server` struct.
    - Initializes the HTTP server with endpoints.

- **handlers.go:**
    - Implements the `/hello` and `/data` endpoints.
    - Uses `RequestParams` struct and `parseRequestParams` function for parsing query parameters.

- **datasource/core.go:**
    - Implements the logic for generating deterministic time series data.
    - `Query` function creates a `Cursor` to iterate over the data points.
    - `Cursor` struct handles the generation of data points using cosine functions for a pseudo-random yet deterministic
      output.

#### Testing

Run tests:

   ```sh
   go test ./...
   ```

#### Contribution

Contributions are welcome. Please fork the repository and create a pull request for any improvements or bug fixes.

#### License

This project is licensed under the MIT License.

For more details, check out the [repository](https://github.com/clevertang/tsserv-public).

#### performance testing
benchmarking the server using ab tool
```sh
ab -n 1000 -c 10 "http://localhost:8080/data?begin=2010-01-01T00:00:00Z&end=2020-01-02T00:00:00Z"
```
results:
```
Concurrency Level:      10
Time taken for tests:   920.180 seconds
Complete requests:      1000
Failed requests:        216
   (Connect: 0, Receive: 0, Length: 216, Exceptions: 0)
Total transferred:      132859540464 bytes
HTML transferred:       132859443464 bytes
Requests per second:    1.09 [#/sec] (mean)
Time per request:       9201.802 [ms] (mean)
Time per request:       920.180 [ms] (mean, across all concurrent requests)
Transfer rate:          141000.26 [Kbytes/sec] received
```



updated the server by using cache for updateBaseValues and write to http response in chunks
```sh
ab -n 1000 -c 10 "http://localhost:8080/data?begin=2010-01-01T00:00:00Z&end=2020-01-02T00:00:00Z"
```
results:
```
Concurrency Level:      10
Time taken for tests:   1053.168 seconds
Complete requests:      1000
Failed requests:        999
   (Connect: 0, Receive: 0, Length: 999, Exceptions: 0)
Total transferred:      52198527000 bytes
HTML transferred:       52198445000 bytes
Requests per second:    0.95 [#/sec] (mean)
Time per request:       10531.676 [ms] (mean)
Time per request:       1053.168 [ms] (mean, across all concurrent requests)
Transfer rate:          48401.72 [Kbytes/sec] received
```