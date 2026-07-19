# ShopMonitor - Personal Shopping Monitoring Platform

## Executive Summary

ShopMonitor is a production-grade, cross-platform personal shopping monitoring platform that continuously tracks products across multiple e-commerce websites and instantly notifies users of price drops, stock availability, size/color variants, coupons, and delivery options.

### Vision
Build a polished, extensible application comparable to CamelCamelCamel, Keepa, and Honey—but focused on personal monitoring with superior real-time capabilities and multi-provider support.

---

## System Architecture

### High-Level Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           CLIENT LAYER                                   │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                  │
│  │   Web App    │  │  Desktop App │  │  Mobile App  │  (Future)        │
│  │  React+TS    │  │   Tauri/Electron │  Android/iOS │                  │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘                  │
│         │                 │                 │                           │
│         └─────────────────┴─────────────────┘                           │
│                           │                                             │
│                    HTTPS / WebSocket                                    │
└───────────────────────────┼─────────────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────────────┐
│                         API GATEWAY LAYER                                │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │                    Go Fiber REST API                             │  │
│  │  • Authentication Middleware                                     │  │
│  │  • Rate Limiting                                                 │  │
│  │  • Request Validation                                            │  │
│  │  • CORS Handler                                                  │  │
│  │  • API Versioning (/api/v1)                                      │  │
│  └──────────────────────────────────────────────────────────────────┘  │
└───────────────────────────┬─────────────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────────────┐
│                        APPLICATION LAYER                                 │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐         │
│  │ Auth Service    │  │ Product Service │  │ Monitor Service │         │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐         │
│  │ Alert Service   │  │ Worker Service  │  │ Notify Service  │         │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘         │
└───────────────────────────┬─────────────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────────────┐
│                         DOMAIN LAYER                                     │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │                    Provider Interface                            │  │
│  │  • Discover()       • FetchPrice()      • FetchVariants()       │  │
│  │  • Monitor()        • FetchCoupons()    • FetchDelivery()       │  │
│  │  • FetchStock()     • FetchMetadata()   • HealthCheck()         │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                                                        │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐          │
│  │  Amazon    │ │  Flipkart  │ │   Myntra   │ │   Snitch   │ ...      │
│  │  Provider  │ │  Provider  │ │  Provider  │ │  Provider  │          │
│  └────────────┘ └────────────┘ └────────────┘ └────────────┘          │
└───────────────────────────┬─────────────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────────────┐
│                        INFRASTRUCTURE LAYER                              │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                  │
│  │  Repository  │  │   Queue      │  │    Cache     │                  │
│  │   (GORM)     │  │  (Redis)     │  │   (Redis)    │                  │
│  └──────────────┘  └──────────────┘  └──────────────┘                  │
│  ┌──────────────┐  ┌──────────────┐                                    │
│  │   Browser    │  │   HTTP       │                                    │
│  │  Automation  │  │   Client     │                                    │
│  │  (Playwright)│  │   Pool       │                                    │
│  └──────────────┘  └──────────────┘                                    │
└───────────────────────────┬─────────────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────────────┐
│                          DATA LAYER                                      │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐                    ┌─────────────────┐            │
│  │   PostgreSQL    │                    │     Redis       │            │
│  │   (Primary DB)  │                    │  (Cache/Queue)  │            │
│  └─────────────────┘                    └─────────────────┘            │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Technology Stack Decisions

### Frontend Stack

| Technology | Purpose | Rationale |
|------------|---------|-----------|
| **React 18** | UI Framework | Component-based, large ecosystem, excellent TypeScript support |
| **TypeScript 5.x** | Type Safety | Compile-time error detection, better IDE support, maintainability |
| **Vite** | Build Tool | Fast HMR, optimized production builds, modern ES modules |
| **Tailwind CSS** | Styling | Utility-first, consistent design system, small bundle size |
| **shadcn/ui** | Component Library | Accessible, customizable, built on Radix UI primitives |
| **TanStack Query** | Data Fetching | Caching, background updates, optimistic updates |
| **React Router v6** | Routing | Declarative routing, nested routes, loaders |
| **Recharts** | Charts | Composable, responsive, built on D3 |
| **Zustand** | State Management | Minimal boilerplate, no context provider needed |

### Backend Stack

| Technology | Purpose | Rationale |
|------------|---------|-----------|
| **Go 1.25+** | Backend Language | Performance, concurrency, small memory footprint |
| **Fiber** | Web Framework | Express-like API, fastest Go web framework, middleware support |
| **GORM** | ORM | Full-featured, migrations, associations, hooks |
| **PostgreSQL** | Primary Database | ACID compliance, JSONB support, full-text search |
| **Redis** | Cache & Queue | In-memory speed, pub/sub, streams for job queue |
| **Playwright** | Browser Automation | Reliable, fast, multi-browser support, auto-wait |

### Infrastructure

| Technology | Purpose | Rationale |
|------------|---------|-----------|
| **Docker** | Containerization | Consistent environments, easy deployment |
| **Docker Compose** | Local Development | Multi-container orchestration |
| **GitHub Actions** | CI/CD | Native GitHub integration, free for public repos |
| **Prometheus** | Metrics | Time-series database, powerful querying |
| **Grafana** | Visualization | Beautiful dashboards, alerting |

---

## Folder Structure

```
shopmonitor/
├── .github/
│   └── workflows/
│       ├── ci.yml              # Continuous Integration
│       ├── cd.yml              # Continuous Deployment
│       └── security-scan.yml   # Security scanning
│
├── docs/
│   ├── ARCHITECTURE.md         # This file
│   ├── API.md                  # OpenAPI/Swagger documentation
│   ├── DATABASE.md             # ER diagrams, schema details
│   ├── DEPLOYMENT.md           # Deployment guides
│   ├── DEVELOPMENT.md          # Developer setup guide
│   ├── CONTRIBUTING.md         # Contribution guidelines
│   └── PROVIDERS.md            # How to add new providers
│
├── backend/
│   ├── cmd/
│   │   ├── api/                # API server entry point
│   │   │   └── main.go
│   │   ├── worker/             # Worker process entry point
│   │   │   └── main.go
│   │   └── migrate/            # Database migration tool
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── config/             # Configuration management
│   │   │   ├── config.go
│   │   │   └── env.go
│   │   │
│   │   ├── domain/             # Domain models & interfaces
│   │   │   ├── entities/
│   │   │   │   ├── user.go
│   │   │   │   ├── product.go
│   │   │   │   ├── price_history.go
│   │   │   │   └── ...
│   │   │   ├── repository/
│   │   │   │   ├── user_repository.go
│   │   │   │   ├── product_repository.go
│   │   │   │   └── ...
│   │   │   └── provider/
│   │   │       ├── interface.go       # Provider contract
│   │   │       ├── base.go            # Base provider implementation
│   │   │       └── registry.go        # Provider registry
│   │   │
│   │   ├── infrastructure/     # External dependencies
│   │   │   ├── database/
│   │   │   │   ├── postgres.go
│   │   │   │   └── migrations/
│   │   │   ├── cache/
│   │   │   │   └── redis.go
│   │   │   ├── queue/
│   │   │   │   └── redis_queue.go
│   │   │   ├── browser/
│   │   │   │   └── playwright.go
│   │   │   └── http/
│   │   │       └── client.go
│   │   │
│   │   ├── repository/         # Repository implementations
│   │   │   ├── user_repository.go
│   │   │   ├── product_repository.go
│   │   │   ├── price_repository.go
│   │   │   └── ...
│   │   │
│   │   ├── service/            # Business logic layer
│   │   │   ├── auth_service.go
│   │   │   ├── product_service.go
│   │   │   ├── monitor_service.go
│   │   │   ├── alert_service.go
│   │   │   ├── notification_service.go
│   │   │   └── scheduler_service.go
│   │   │
│   │   ├── handler/            # HTTP handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── product_handler.go
│   │   │   ├── monitor_handler.go
│   │   │   ├── alert_handler.go
│   │   │   └── admin_handler.go
│   │   │
│   │   ├── middleware/         # HTTP middleware
│   │   │   ├── auth.go
│   │   │   ├── ratelimit.go
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   │
│   │   ├── provider/           # Shopping website providers
│   │   │   ├── amazon/
│   │   │   │   ├── amazon.go
│   │   │   │   ├── parser.go
│   │   │   │   └── api.go
│   │   │   ├── flipkart/
│   │   │   ├── myntra/
│   │   │   ├── ajio/
│   │   │   ├── snitch/
│   │   │   ├── nike/
│   │   │   ├── adidas/
│   │   │   ├── hnm/
│   │   │   ├── zara/
│   │   │   ├── levis/
│   │   │   ├── decathlon/
│   │   │   ├── croma/
│   │   │   ├── apple/
│   │   │   ├── samsung/
│   │   │   ├── boat/
│   │   │   ├── nothing/
│   │   │   ├── swiggy/
│   │   │   ├── zomato/
│   │   │   ├── blinkit/
│   │   │   ├── zepto/
│   │   │   └── bigbasket/
│   │   │
│   │   ├── worker/             # Background job processing
│   │   │   ├── worker.go
│   │   │   ├── pool.go
│   │   │   ├── job.go
│   │   │   ├── scheduler.go
│   │   │   └── retry.go
│   │   │
│   │   ├── notification/       # Notification channels
│   │   │   ├── email.go
│   │   │   ├── telegram.go
│   │   │   ├── discord.go
│   │   │   ├── slack.go
│   │   │   ├── desktop.go
│   │   │   └── webhook.go
│   │   │
│   │   └── observability/      # Monitoring & logging
│   │       ├── logger.go
│   │       ├── metrics.go
│   │       └── tracing.go
│   │
│   ├── pkg/                    # Public packages
│   │   ├── crypto/             # Cryptography utilities
│   │   ├── validator/          # Input validation
│   │   └── utils/              # Common utilities
│   │
│   ├── api/                    # API specifications
│   │   └── openapi.yaml
│   │
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── Makefile
│
├── frontend/
│   ├── public/
│   │   ├── favicon.ico
│   │   └── manifest.json
│   │
│   ├── src/
│   │   ├── components/         # Reusable components
│   │   │   ├── ui/             # shadcn/ui components
│   │   │   ├── layout/
│   │   │   │   ├── Header.tsx
│   │   │   │   ├── Sidebar.tsx
│   │   │   │   └── Footer.tsx
│   │   │   ├── product/
│   │   │   │   ├── ProductCard.tsx
│   │   │   │   ├── ProductList.tsx
│   │   │   │   ├── ProductDetail.tsx
│   │   │   │   └── PriceChart.tsx
│   │   │   ├── monitor/
│   │   │   │   ├── MonitorForm.tsx
│   │   │   │   ├── MonitorList.tsx
│   │   │   │   └── MonitorSettings.tsx
│   │   │   ├── alert/
│   │   │   │   ├── AlertList.tsx
│   │   │   │   └── AlertSettings.tsx
│   │   │   └── dashboard/
│   │   │       ├── Dashboard.tsx
│   │   │       ├── StatsCard.tsx
│   │   │       └── RecentAlerts.tsx
│   │   │
│   │   ├── pages/              # Page components
│   │   │   ├── Home.tsx
│   │   │   ├── Dashboard.tsx
│   │   │   ├── Products.tsx
│   │   │   ├── Monitors.tsx
│   │   │   ├── Alerts.tsx
│   │   │   ├── Settings.tsx
│   │   │   ├── Login.tsx
│   │   │   └── Admin.tsx
│   │   │
│   │   ├── hooks/              # Custom React hooks
│   │   │   ├── useAuth.ts
│   │   │   ├── useProducts.ts
│   │   │   ├── useMonitors.ts
│   │   │   └── useNotifications.ts
│   │   │
│   │   ├── services/           # API clients
│   │   │   ├── api.ts
│   │   │   ├── auth.ts
│   │   │   ├── products.ts
│   │   │   └── monitors.ts
│   │   │
│   │   ├── store/              # Zustand stores
│   │   │   ├── authStore.ts
│   │   │   ├── productStore.ts
│   │   │   └── uiStore.ts
│   │   │
│   │   ├── types/              # TypeScript types
│   │   │   ├── api.ts
│   │   │   ├── product.ts
│   │   │   └── user.ts
│   │   │
│   │   ├── utils/              # Utility functions
│   │   │   ├── format.ts
│   │   │   ├── validators.ts
│   │   │   └── constants.ts
│   │   │
│   │   ├── styles/             # Global styles
│   │   │   └── globals.css
│   │   │
│   │   ├── App.tsx
│   │   ├── main.tsx
│   │   └── vite-env.d.ts
│   │
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   ├── tailwind.config.js
│   ├── postcss.config.js
│   ├── vite.config.ts
│   └── Dockerfile
│
├── deploy/
│   ├── docker-compose.yml      # Production deployment
│   ├── docker-compose.dev.yml  # Development environment
│   ├── kubernetes/             # K8s manifests (future)
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── ingress.yaml
│   └── prometheus/
│       ├── prometheus.yml
│       └── grafana/
│           └── dashboards/
│
├── scripts/
│   ├── setup.sh                # Initial setup script
│   ├── dev.sh                  # Start development environment
│   ├── test.sh                 # Run all tests
│   ├── build.sh                # Build all components
│   └── migrate.sh              # Run database migrations
│
├── .gitignore
├── .env.example
├── LICENSE
└── README.md
```

---

## Database Design

### Entity Relationship Diagram

```
┌─────────────────────┐       ┌─────────────────────┐
│       users         │       │      sessions       │
├─────────────────────┤       ├─────────────────────┤
│ id (PK)             │◄──────│ id (PK)             │
│ email               │       │ user_id (FK)        │
│ password_hash       │       │ refresh_token       │
│ name                │       │ expires_at          │
│ avatar_url          │       │ created_at          │
│ timezone            │       │ updated_at          │
│ currency            │       └─────────────────────┘
│ language            │
│ created_at          │       ┌─────────────────────┐
│ updated_at          │       │     providers       │
└─────────┬───────────┘       ├─────────────────────┤
          │                   │ id (PK)             │
          │                   │ name                │
          │                   │ base_url            │
          │                   │ status              │
          │                   │ last_check          │
          │                   └─────────────────────┘
          │
          │                   ┌─────────────────────┐
          │                   │      products       │
          │                   ├─────────────────────┤
          │                   │ id (PK)             │
          │                   │ provider_id (FK)    │
          │                   │ external_id         │
          │                   │ url                 │
          │                   │ title               │
          │                   │ brand               │
          │                   │ category            │
          │                   │ description         │
          │                   │ image_url           │
          │                   │ rating              │
          │                   │ review_count        │
          │                   │ seller              │
          │                   │ shipping_cost       │
          │                   │ return_policy       │
          │                   │ metadata (JSONB)    │
          │                   │ created_at          │
          │                   │ updated_at          │
          │                   └─────────┬───────────┘
          │                             │
          │                   ┌─────────┴───────────┐
          │                   │  product_variants   │
          │                   ├─────────────────────┤
          │                   │ id (PK)             │
          │                   │ product_id (FK)     │
          │                   │ sku                 │
          │                   │ variant_type        │
          │                   │ variant_value       │
          │                   │ available           │
          │                   │ price               │
          │                   │ mrp                 │
          │                   │ discount_percent    │
          │                   │ created_at          │
          │                   │ updated_at          │
          │                   └─────────┬───────────┘
          │                             │
          │                   ┌─────────┴───────────┐
          │                   │   product_images    │
          │                   ├─────────────────────┤
          │                   │ id (PK)             │
          │                   │ product_id (FK)     │
          │                   │ url                 │
          │                   │ alt_text            │
          │                   │ position            │
          │                   └─────────────────────┘
          │
┌─────────┴───────────┐
│ user_products       │  (Wishlist/Favorites)
├─────────────────────┤
│ id (PK)             │
│ user_id (FK)        │
│ product_id (FK)     │
│ nickname            │
│ tags (TEXT[])       │
│ is_favorite         │
│ created_at          │
└─────────┬───────────┘
          │
┌─────────┴───────────┐
│     monitors        │
├─────────────────────┤
│ id (PK)             │
│ user_id (FK)        │
│ product_id (FK)     │
│ status              │
│ check_interval      │
│ target_price        │
│ target_discount     │
│ desired_sizes ([])  │
│ desired_colors ([]) │
│ max_price           │
│ min_discount        │
│ delivery_pincode    │
│ seller_preference   │
│ notification_channels│
│ last_check          │
│ next_check          │
│ paused_at           │
│ created_at          │
│ updated_at          │
└─────────┬───────────┘
          │
┌─────────┴───────────────────┐
│     price_history           │
├─────────────────────────────┤
│ id (PK)                     │
│ monitor_id (FK)             │
│ product_variant_id (FK)     │
│ price                       │
│ mrp                         │
│ discount_percent            │
│ currency                    │
│ seller                      │
│ in_stock                    │
│ observed_at                 │
└─────────────────────────────┘

┌─────────────────────────────┐
│      stock_history          │
├─────────────────────────────┤
│ id (PK)                     │
│ monitor_id (FK)             │
│ product_variant_id (FK)     │
│ available                   │
│ quantity_available          │
│ observed_at                 │
└─────────────────────────────┘

┌─────────────────────────────┐
│       coupons               │
├─────────────────────────────┤
│ id (PK)                     │
│ product_id (FK)             │
│ code                        │
│ description                 │
│ discount_type               │
│ discount_value              │
│ min_order_value             │
│ max_discount                │
│ valid_from                  │
│ valid_until                 │
│ terms (TEXT[])              │
│ active                      │
│ discovered_at               │
└─────────────────────────────┘

┌─────────────────────────────┐
│    delivery_history         │
├─────────────────────────────┤
│ id (PK)                     │
│ monitor_id (FK)             │
│ pincode                     │
│ available                   │
│ delivery_type               │
│ estimated_days              │
│ shipping_cost               │
│ observed_at                 │
└─────────────────────────────┘

┌─────────────────────────────┐
│      notifications          │
├─────────────────────────────┤
│ id (PK)                     │
│ user_id (FK)                │
│ monitor_id (FK)             │
│ type                        │
│ title                       │
│ message                     │
│ data (JSONB)                │
│ read                        │
│ read_at                     │
│ created_at                  │
└─────────────────────────────┘

┌─────────────────────────────┐
│   notification_logs         │
├─────────────────────────────┤
│ id (PK)                     │
│ notification_id (FK)        │
│ channel                     │
│ status                      │
│ error_message               │
│ sent_at                     │
│ delivered_at                │
└─────────────────────────────┘

┌─────────────────────────────┐
│          jobs               │
├─────────────────────────────┤
│ id (PK)                     │
│ type                        │
│ payload (JSONB)             │
│ status                      │
│ priority                    │
│ attempts                    │
│ max_attempts                │
│ scheduled_at                │
│ started_at                  │
│ completed_at                │
│ failed_at                   │
│ error_message               │
│ created_at                  │
└─────────────────────────────┘

┌─────────────────────────────┐
│        workers              │
├─────────────────────────────┤
│ id (PK)                     │
│ hostname                    │
│ pid                         │
│ status                      │
│ current_job_id (FK)         │
│ jobs_completed              │
│ jobs_failed                 │
│ last_heartbeat              │
│ started_at                  │
└─────────────────────────────┘

┌─────────────────────────────┐
│       user_settings         │
├─────────────────────────────┤
│ id (PK)                     │
│ user_id (FK)                │
│ theme                       │
│ timezone                    │
│ currency                    │
│ language                    │
│ notification_prefs (JSONB)  │
│ created_at                  │
│ updated_at                  │
└─────────────────────────────┘
```

### Schema Details

See `docs/DATABASE.md` for complete SQL schema definitions, indexes, and constraints.

---

## Provider Interface Design

### Core Interface

Every shopping website provider must implement the following interface:

```go
// Provider defines the contract for all shopping website providers
type Provider interface {
    // ID returns the unique identifier for this provider
    ID() string
    
    // Name returns the human-readable name
    Name() string
    
    // BaseURL returns the base URL of the shopping website
    BaseURL() string
    
    // Capabilities returns what this provider supports
    Capabilities() ProviderCapabilities
    
    // Discovery
    
    // Discover extracts product information from a URL
    Discover(ctx context.Context, url string) (*ProductDiscovery, error)
    
    // ValidateURL checks if the URL belongs to this provider
    ValidateURL(url string) bool
    
    // Monitoring
    
    // Monitor performs a complete product check
    Monitor(ctx context.Context, product *Product) (*MonitorResult, error)
    
    // FetchPrice gets the current price for a product
    FetchPrice(ctx context.Context, product *Product) (*PriceInfo, error)
    
    // FetchVariants gets all available variants (sizes, colors, etc.)
    FetchVariants(ctx context.Context, product *Product) ([]*Variant, error)
    
    // FetchCoupons gets available coupons for the product
    FetchCoupons(ctx context.Context, product *Product) ([]*Coupon, error)
    
    // FetchDelivery checks delivery availability for a pincode
    FetchDelivery(ctx context.Context, product *Product, pincode string) (*DeliveryInfo, error)
    
    // FetchStock gets current stock status
    FetchStock(ctx context.Context, product *Product) (*StockInfo, error)
    
    // FetchMetadata gets additional product metadata
    FetchMetadata(ctx context.Context, product *Product) (*ProductMetadata, error)
    
    // Health
    
    // HealthCheck verifies the provider is working
    HealthCheck(ctx context.Context) (*HealthStatus, error)
}

// ProviderCapabilities describes what features a provider supports
type ProviderCapabilities struct {
    PriceTracking     bool
    StockTracking     bool
    VariantTracking   bool
    CouponDetection   bool
    DeliveryTracking  bool
    SupportsSizes     bool
    SupportsColors    bool
    SupportsElectronics bool
    SupportsFood      bool
    HasAPI            bool
    RequiresBrowser   bool
}
```

### Base Provider Implementation

To reduce code duplication, we provide a base provider with common functionality:

```go
// BaseProvider provides common functionality for all providers
type BaseProvider struct {
    id           string
    name         string
    baseURL      string
    httpClient   *http.Client
    browser      *playwright.Browser
    cache        cache.Cache
    rateLimiter  *ratelimit.Limiter
}

// Common methods implemented by BaseProvider:
// - Default HTTP client with retries
// - Browser lifecycle management
// - Response caching
// - Rate limiting
// - User agent rotation
// - Cookie persistence
```

### Provider Registry

Providers are registered in a central registry for dynamic discovery:

```go
// Registry manages all available providers
type Registry struct {
    providers map[string]Provider
    mu        sync.RWMutex
}

func (r *Registry) Register(provider Provider)
func (r *Registry) Get(id string) (Provider, bool)
func (r *Registry) GetAll() []Provider
func (r *Registry) MatchURL(url string) (Provider, bool)
```

### Adding a New Provider

To add a new shopping website:

1. Create a new directory under `backend/internal/provider/{website}/`
2. Implement the `Provider` interface
3. Register the provider in `backend/internal/provider/registry.go`
4. Add tests

Example structure for a new provider:

```
providers/amazon/
├── amazon.go          # Provider implementation
├── parser.go          # HTML parsing logic
├── api.go             # API endpoint calls
├── types.go           # Amazon-specific types
└── amazon_test.go     # Tests
```

See `docs/PROVIDERS.md` for detailed implementation guide.

---

## API Design

### REST API Endpoints

#### Authentication

```
POST   /api/v1/auth/register          # Register new user
POST   /api/v1/auth/login             # Login
POST   /api/v1/auth/logout            # Logout
POST   /api/v1/auth/refresh           # Refresh token
POST   /api/v1/auth/forgot-password   # Request password reset
POST   /api/v1/auth/reset-password    # Reset password
GET    /api/v1/auth/me                # Get current user
PUT    /api/v1/auth/me                # Update current user

# OAuth
GET    /api/v1/auth/oauth/google      # Initiate Google OAuth
GET    /api/v1/auth/oauth/github      # Initiate GitHub OAuth
GET    /api/v1/auth/oauth/callback    # OAuth callback
```

#### Products

```
GET    /api/v1/products               # List products (with filters)
POST   /api/v1/products               # Add product to monitor
GET    /api/v1/products/:id           # Get product details
PUT    /api/v1/products/:id           # Update product
DELETE /api/v1/products/:id           # Remove product
GET    /api/v1/products/:id/history   # Get price history
GET    /api/v1/products/:id/variants  # Get product variants
POST   /api/v1/products/discover      # Discover product from URL
```

#### Monitors

```
GET    /api/v1/monitors               # List user's monitors
POST   /api/v1/monitors               # Create new monitor
GET    /api/v1/monitors/:id           # Get monitor details
PUT    /api/v1/monitors/:id           # Update monitor
DELETE /api/v1/monitors/:id           # Delete monitor
POST   /api/v1/monitors/:id/pause     # Pause monitoring
POST   /api/v1/monitors/:id/resume    # Resume monitoring
POST   /api/v1/monitors/:id/check     # Force immediate check
GET    /api/v1/monitors/:id/history   # Get monitoring history
```

#### Alerts

```
GET    /api/v1/alerts                 # List alerts
GET    /api/v1/alerts/unread          # List unread alerts
PUT    /api/v1/alerts/:id/read        # Mark alert as read
PUT    /api/v1/alerts/read-all        # Mark all as read
DELETE /api/v1/alerts/:id             # Delete alert
GET    /api/v1/alerts/settings        # Get alert settings
PUT    /api/v1/alerts/settings        # Update alert settings
```

#### Notifications

```
GET    /api/v1/notifications          # List notifications
PUT    /api/v1/notifications/:id/read # Mark as read
PUT    /api/v1/notifications/read-all # Mark all as read
GET    /api/v1/notifications/channels # List configured channels
POST   /api/v1/notifications/channels # Add notification channel
DELETE /api/v1/notifications/channels/:id
```

#### Dashboard

```
GET    /api/v1/dashboard/stats        # Get dashboard statistics
GET    /api/v1/dashboard/savings      # Get savings summary
GET    /api/v1/dashboard/recent-alerts
GET    /api/v1/dashboard/top-discounts
```

#### Settings

```
GET    /api/v1/settings               # Get user settings
PUT    /api/v1/settings               # Update user settings
GET    /api/v1/settings/preferences   # Get preferences
PUT    /api/v1/settings/preferences   # Update preferences
```

#### Admin

```
GET    /api/v1/admin/workers          # List workers
GET    /api/v1/admin/jobs             # List jobs
GET    /api/v1/admin/metrics          # Get system metrics
GET    /api/v1/admin/providers        # List providers status
POST   /api/v1/admin/providers/:id/check # Force health check
```

### WebSocket Events

Real-time updates via WebSocket:

```
# Server → Client events
alert:new           # New alert triggered
notification:new    # New notification
monitor:update      # Monitor status changed
product:price       # Price update
product:stock       # Stock update
dashboard:update    # Dashboard stats changed

# Client → Server events
subscribe:alerts    # Subscribe to alerts
subscribe:monitors  # Subscribe to specific monitors
unsubscribe         # Unsubscribe from channel
```

### API Response Format

```json
{
  "success": true,
  "data": { },
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  },
  "error": null
}
```

### Error Response Format

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  }
}
```

See `docs/API.md` for complete OpenAPI specification.

---

## Development Roadmap

### Phase 1: Foundation ✅ (Current)
- [x] Architecture design
- [x] Technology decisions
- [x] Database schema
- [x] Provider interface
- [x] API design
- [ ] Project scaffolding

### Phase 2: Backend Core
- [ ] Database migrations
- [ ] Repository implementations
- [ ] Authentication service
- [ ] User management
- [ ] Basic CRUD APIs
- [ ] Provider registry
- [ ] First provider (Amazon)

### Phase 3: Monitoring Engine
- [ ] Scheduler implementation
- [ ] Worker pool
- [ ] Job queue
- [ ] Retry mechanism
- [ ] Circuit breaker
- [ ] Additional providers (Flipkart, Myntra)

### Phase 4: Frontend Foundation
- [ ] Project setup
- [ ] Authentication UI
- [ ] Dashboard layout
- [ ] Product listing
- [ ] Product detail page

### Phase 5: Monitoring UI
- [ ] Add monitor form
- [ ] Monitor list
- [ ] Monitor settings
- [ ] Price charts
- [ ] Stock history

### Phase 6: Notifications
- [ ] Email notifications
- [ ] Telegram integration
- [ ] Discord integration
- [ ] Desktop notifications
- [ ] Notification preferences UI

### Phase 7: Advanced Features
- [ ] Coupon detection
- [ ] Delivery tracking
- [ ] Size/color monitoring
- [ ] Food delivery providers
- [ ] Import/export

### Phase 8: Polish & Production
- [ ] Admin panel
- [ ] Observability (Prometheus, Grafana)
- [ ] Performance optimization
- [ ] Security hardening
- [ ] Documentation
- [ ] Testing suite
- [ ] CI/CD pipeline

### Phase 9: Future Enhancements
- [ ] Browser extension
- [ ] Mobile apps
- [ ] AI features
- [ ] Multi-user support
- [ ] Price comparison

---

## Key Architectural Decisions

### 1. Clean Architecture
We use Clean Architecture to ensure separation of concerns, testability, and maintainability. Dependencies point inward, with business logic at the core.

### 2. Provider-Based Design
Each shopping website is an independent module implementing a common interface. This allows:
- Easy addition of new providers
- Independent testing
- Different strategies per provider (API vs scraping)
- Graceful degradation if one provider fails

### 3. Worker Pool Pattern
Background jobs are processed by a configurable worker pool:
- Dynamic scaling based on queue size
- Isolation of failures
- Resource control
- Priority queues for urgent checks

### 4. Event-Driven Notifications
Notifications are decoupled from monitoring via events:
- Multiple channels per event
- Async processing
- Retry on failure
- User preferences respected

### 5. Caching Strategy
Multi-layer caching:
- L1: In-memory cache for hot data
- L2: Redis for shared cache
- TTL-based invalidation
- Write-through for critical data

### 6. Database Design
- Normalized schema for data integrity
- JSONB columns for flexible metadata
- Time-series tables for history
- Proper indexing for query patterns

### 7. Security First
- JWT with short expiry + refresh tokens
- Password hashing with Argon2
- Input validation on all endpoints
- Rate limiting per user/IP
- CSP headers
- Encrypted secrets storage

---

## Next Steps

This completes Phase 1. The architecture is designed for:

1. **Extensibility**: New providers require minimal code
2. **Performance**: Efficient caching, connection pooling, worker pools
3. **Reliability**: Retries, circuit breakers, health checks
4. **Maintainability**: Clean architecture, comprehensive testing
5. **Security**: Defense in depth, secure defaults

**Awaiting approval before proceeding to Phase 2: Backend Core Implementation**

Please review and confirm if you'd like me to proceed with implementing the backend foundation.
