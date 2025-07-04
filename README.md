
# GoRate

GoRate is a lightweight, high-performance currency exchange rate API built in Go using the Gin-Gonic framework. It provides real-time currency conversion rates and cryptocurrency exchange functionality.

## Features

- Real-time currency exchange rates via OpenExchangeRates API
- Cryptocurrency conversion with fixed rates
- RESTful API with JSON responses
- Containerized with Docker for easy deployment
- Configurable via environment variables
- Graceful shutdown handling

## Installation

### Prerequisites

- Go 1.24 or higher
- Docker (optional, for containerized deployment)
- OpenExchangeRates API key (get one at [openexchangerates.org](https://openexchangerates.org/))

### Local Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/IAmRadek/gorate.git
   cd gorate
   ```

2. Create a `.development.env` file with your configuration:
   ```
   ADDR=:8080
   READ_TIMEOUT=10s
   READ_HEADER_TIMEOUT=10s
   WRITE_TIMEOUT=10s
   IDLE_TIMEOUT=10s
   GRACEFUL_TIMEOUT=10s
   MAX_HEADER_BYTES=1024
   OPEN_EXCHANGE_RATES_PROVIDER_APP_ID=your_api_key_here
   ```

3. Run the application:
   ```bash
   make run
   ```

### Docker Setup

1. Build the Docker image:
   ```bash
   make build-dockerimage
   ```

2. Run the container:
   ```bash
   make run-docker
   ```

## API Documentation

### GET /rates

Retrieves exchange rates between multiple currencies.

**Query Parameters:**
- `currencies` (required): Comma-separated list of currency codes (minimum 2)

**Example Request:**
```
GET /rates?currencies=USD,GBP,EUR
```

**Example Response:**
```json
[
  { "from": "USD", "to": "GBP", "rate": 0.732787 },
  { "from": "GBP", "to": "USD", "rate": 1.364653030143820783 },
  { "from": "USD", "to": "EUR", "rate": 0.851239 },
  { "from": "EUR", "to": "USD", "rate": 1.174758205392375114 },
  { "from": "EUR", "to": "GBP", "rate": 0.860847541054862383 },
  { "from": "GBP", "to": "EUR", "rate": 1.161645880726595859 }
]
```

### GET /exchange

Converts between cryptocurrencies using fixed rates.

**Query Parameters:**
- `from` (required): Source cryptocurrency code
- `to` (required): Target cryptocurrency code
- `amount` (required): Amount to convert (must be positive)

**Supported Cryptocurrencies:**
- BEER
- FLOKI
- GATE
- USDT
- WBTC

**Example Request:**
```
GET /exchange?from=WBTC&to=USDT&amount=1.0
```

**Example Response:**
```json
{
  "from": "WBTC",
  "to": "USDT",
  "amount": 57613.353535
}
```

## Configuration

GoRate can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `GIN_MODE` | Gin framework mode (debug/release) | debug |
| `ADDR` | Server address and port | :8080 |
| `READ_TIMEOUT` | HTTP read timeout | 10s |
| `READ_HEADER_TIMEOUT` | HTTP header read timeout | 10s |
| `WRITE_TIMEOUT` | HTTP write timeout | 10s |
| `IDLE_TIMEOUT` | HTTP idle connection timeout | 10s |
| `MAX_HEADER_BYTES` | Maximum HTTP header size | 1024 |
| `GRACEFUL_SHUTDOWN_DURATION` | Graceful shutdown timeout | 5s |
| `OPEN_EXCHANGE_RATES_PROVIDER_APP_ID` | OpenExchangeRates API key | (required) |

## Development

### Project Structure

```
gorate/
├── cmd/
│   └── gorate/           # Application entry point
├── internal/
│   ├── exchanges/        # Exchange functionality
│   └── rates/            # Rate providers and models
├── Dockerfile            # Docker configuration
├── Makefile              # Build and run commands
└── .development.env      # Environment configuration
```

### Available Make Commands

- `make run`: Run the application locally
- `make tests`: Run all tests
- `make build-dockerimage`: Build the Docker image
- `make run-docker`: Run the application in Docker
- `make clean`: Clean build artifacts

## License

Copyright © 2024 IAmRadek. All rights reserved.
