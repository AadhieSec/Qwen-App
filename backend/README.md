# ShopMonitor - Intelligent Product Price & Stock Monitoring System

A high-performance, scalable backend system for monitoring product prices, stock availability, and sending real-time notifications across multiple e-commerce platforms.

## Features

### Core Capabilities
- **Multi-Provider Support**: Amazon, Flipkart, and extensible architecture for additional platforms
- **Intelligent Scraping**: Adaptive strategies using API calls, network inspection, HTML parsing, and browser automation
- **Real-time Monitoring**: Continuous price and stock tracking with configurable intervals
- **Smart Notifications**: Multi-channel alerts via Desktop, Telegram, Discord, Slack, Email, Webhooks, and Push notifications
- **Historical Tracking**: Complete price history, stock changes, variant availability, and coupon tracking
- **User Management**: JWT authentication, OAuth (Google, GitHub), role-based access control

### Technical Highlights
- **High Performance**: Worker pool architecture with job queue for parallel processing
- **Scalable Design**: Redis-backed caching and queue management
- **Resilient**: Automatic retry mechanisms, dead letter queues, circuit breakers
- **Observable**: Prometheus metrics, structured logging, health checks
- **Secure**: Password hashing, rate limiting, input validation, CORS protection

## Project Structure

```
backend/
├── cmd/
│   └── api/              # Main application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── models/           # Database models and schemas
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic layer
│   ├── handler/          # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── provider/         # E-commerce platform providers
│   │   ├── base/         # Base provider interface
│   │   ├── amazon/       # Amazon implementation
│   │   └── flipkart/     # Flipkart implementation
│   ├── queue/            # Job queue management
│   ├── worker/           # Background job processors
│   ├── scheduler/        # Scheduled task management
│   └── notification/     # Notification delivery system
├── pkg/
│   ├── logger/           # Logging utilities
│   └── utils/            # Common utilities
├── migrations/           # Database migration files
├── scripts/              # Utility scripts
└── tests/                # Test files
```

## Quick Start

### Prerequisites
- Go 1.19+
- PostgreSQL 14+
- Redis 6+
- Node.js 18+ (for Playwright browser automation)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd shopmonitor/backend
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Update `.env` with your configuration

4. Install dependencies:
```bash
go mod download
```

5. Run the application:
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/google` - Google OAuth
- `GET /api/v1/auth/github` - GitHub OAuth

### Products
- `GET /api/v1/products` - List user's products
- `POST /api/v1/products` - Add product to monitor
- `GET /api/v1/products/:id` - Get product details
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Remove product
- `GET /api/v1/products/:id/history` - Get price history

### Notifications
- `GET /api/v1/notifications` - List notifications
- `PUT /api/v1/notifications/:id/read` - Mark as read
- `GET /api/v1/notifications/preferences` - Get preferences
- `PUT /api/v1/notifications/preferences` - Update preferences

## Configuration

See `.env.example` for all available configuration options:

- **Server**: Port, timeouts, environment
- **Database**: PostgreSQL connection settings
- **Redis**: Cache and queue configuration
- **Auth**: JWT secrets, OAuth credentials
- **Queue**: Job queue settings
- **Worker**: Worker pool configuration
- **Notification**: Channel-specific settings
- **Playwright**: Browser automation settings

## Architecture

### Provider System
The provider system uses a strategy pattern to support multiple e-commerce platforms:

1. **Base Provider**: Common functionality (HTTP client, HTML parsing)
2. **Platform Providers**: Platform-specific implementations (Amazon, Flipkart)
3. **Adaptive Scraping**: Automatically selects best method (API → Network → HTML → Browser)

### Worker Pool
Background jobs are processed by a configurable worker pool:

- Jobs are queued in Redis
- Workers pick up jobs based on priority
- Failed jobs are retried with exponential backoff
- Dead letter queue for permanently failed jobs

### Notification System
Multi-channel notification delivery:

- User-configurable preferences per event type
- Parallel delivery to multiple channels
- Delivery tracking and retry logic
- Rate limiting per channel

## Development

### Running Tests
```bash
go test ./...
```

### Code Generation
```bash
# Generate Swagger docs
swag init -g cmd/api/main.go

# Run mocks
mockgen -source=internal/provider/base/interface.go -destination=internal/provider/base/mock.go
```

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
