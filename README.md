# GF-Scorer

GF-Scorer is a high-performance Go application designed to process large volumes of log files, calculate various scores for each line, and store the results in a PostgreSQL database. It's built with concurrency in mind and includes Prometheus metrics for monitoring.

## Features

- Concurrent processing of multiple log files
- Calculation of repeat, increasing, decreasing, and ML scores for each line
- Efficient batch insertion into PostgreSQL database
- Configurable processing parameters
- Prometheus metrics for monitoring performance
- Graceful shutdown handling

## Prerequisites

- Go 1.22 or later
- PostgreSQL database
- Prometheus (for metrics collection)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/iyuangang/gf-scorer.git
   cd gf-scorer
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Build the application:
   ```
   go build -o gf-scorer cmd/scorer/main.go
   ```

## Configuration

Create a `config.json` file in the project root with the following structure:

```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "your_username",
    "password": "your_password",
    "dbname": "your_database",
    "maxOpenConns": 25,
    "maxIdleConns": 10,
    "connMaxLifetime": 300
  },
  "processing": {
    "batchSize": 1000,
    "maxConcurrentFiles": 10
  },
  "metrics": {
    "port": 8080
  }
}
```

Adjust the values according to your environment and requirements.

## Usage

Run the application with the following command:

```
./gf-scorer -config config.json -input /path/to/your/logfiles
```

- `-config`: Path to the configuration file (default: "config.json")
- `-input`: Path to the input file or directory (required)

## Metrics

Prometheus metrics are exposed on the configured port (default: 8080) at the `/metrics` endpoint. You can configure Prometheus to scrape these metrics for monitoring the application's performance.

## Project Structure

- `cmd/scorer`: Contains the main application entry point
- `internal/config`: Configuration loading and structures
- `internal/database`: Database connection management
- `internal/metrics`: Prometheus metrics definitions
- `internal/scorer`: Core scoring logic and file processing

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [MIT License](LICENSE).
