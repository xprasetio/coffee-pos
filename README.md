# Coffee Shop POS

Coffee Shop Point of Sales System - A comprehensive POS solution for medium-scale coffee shops.

## Features

- Product & Category Management
- Stock Management with Movement History
- Table Management
- Cashier Shift Management
- Transaction & Payment Processing (Midtrans)
- Promo & Discount Management
- Reports & Analytics (Revenue, Best Sellers, Cashier Performance)
- CSV Export

## Tech Stack

- **Backend:** Go (Golang)
- **Database:** MySQL
- **Cache:** Redis
- **Payment Gateway:** Midtrans Snap

## Getting Started

### Prerequisites

- Go 1.21+
- MySQL 8.0+
- Redis 7.0+

### Installation

1. Clone the repository
2. Copy `.env.example` to `.env` and configure
3. Run migrations
4. Start the server

```bash
go run cmd/api/main.go
```

## Project Structure

```
coffee-pos/
├── cmd/api/           # Application entry point
├── internal/          # Private application code
│   ├── entity/        # Domain structs
│   ├── repository/    # Data access layer
│   ├── service/       # Business logic
│   ├── handler/       # HTTP handlers
│   ├── middleware/    # HTTP middleware
│   └── dto/           # Request/Response structs
├── pkg/               # Public utility packages
│   ├── database/      # Database connection
│   ├── redis/         # Redis connection
│   ├── jwt/           # JWT helpers
│   ├── response/      # Response helpers
│   └── validator/     # Input validation
├── migrations/        # SQL migration files
└── config/            # Configuration management
```

## API Documentation

See [API_CONTRACT.md](./API_CONTRACT.md) for complete API documentation.

## License

MIT
