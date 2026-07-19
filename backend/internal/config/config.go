package config

import (
"os"
"strconv"
"time"

"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
Server       ServerConfig
Database     DatabaseConfig
Redis        RedisConfig
Auth         AuthConfig
Queue        QueueConfig
Worker       WorkerConfig
Notification NotificationConfig
Playwright   PlaywrightConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
Port          string
Host          string
ReadTimeout   time.Duration
WriteTimeout  time.Duration
IdleTimeout   time.Duration
ShutdownDelay time.Duration
Environment   string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
Host            string
Port            string
User            string
Password        string
DBName          string
SSLMode         string
MaxOpenConns    int
MaxIdleConns    int
ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
Host            string
Port            string
Password        string
DB              int
PoolSize        int
MinIdleConns    int
ConnMaxLifetime time.Duration
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
JWTSecret           string
JWTExpirationHours  int
RefreshTokenDays    int
GoogleClientID      string
GoogleClientSecret  string
GitHubClientID      string
GitHubClientSecret  string
}

// QueueConfig holds queue configuration
type QueueConfig struct {
Type              string // redis or rabbitmq
MaxRetries        int
RetryDelaySeconds int
DeadLetterQueue   string
}

// WorkerConfig holds worker pool configuration
type WorkerConfig struct {
NumWorkers      int
MaxJobRetry     int
JobTimeoutSecs  int
HealthCheckSecs int
}

// NotificationConfig holds notification configuration
type NotificationConfig struct {
TelegramBotToken  string
TelegramChatID    string
DiscordWebhookURL string
SlackWebhookURL   string
SMTPHost          string
SMTPPort          int
SMTPUser          string
SMTPPassword      string
SMTPFrom          string
}

// PlaywrightConfig holds browser automation configuration
type PlaywrightConfig struct {
Headless    bool
BrowserType string // chromium, firefox, webkit
Timeout     time.Duration
UseProxy    bool
ProxyURL    string
UserAgent   string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
// Load .env file if it exists
godotenv.Load()

cfg := &Config{
Server: ServerConfig{
Port:          getEnv("SERVER_PORT", "8080"),
Host:          getEnv("SERVER_HOST", "0.0.0.0"),
ReadTimeout:   getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
WriteTimeout:  getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
IdleTimeout:   getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
ShutdownDelay: getDurationEnv("SERVER_SHUTDOWN_DELAY", 10*time.Second),
Environment:   getEnv("ENVIRONMENT", "development"),
},
Database: DatabaseConfig{
Host:            getEnv("DB_HOST", "localhost"),
Port:            getEnv("DB_PORT", "5432"),
User:            getEnv("DB_USER", "postgres"),
Password:        getEnv("DB_PASSWORD", "postgres"),
DBName:          getEnv("DB_NAME", "shopmonitor"),
SSLMode:         getEnv("DB_SSLMODE", "disable"),
MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
},
Redis: RedisConfig{
Host:            getEnv("REDIS_HOST", "localhost"),
Port:            getEnv("REDIS_PORT", "6379"),
Password:        getEnv("REDIS_PASSWORD", ""),
DB:              getIntEnv("REDIS_DB", 0),
PoolSize:        getIntEnv("REDIS_POOL_SIZE", 10),
MinIdleConns:    getIntEnv("REDIS_MIN_IDLE_CONNS", 5),
ConnMaxLifetime: getDurationEnv("REDIS_CONN_MAX_LIFETIME", 5*time.Minute),
},
Auth: AuthConfig{
JWTSecret:          getEnv("JWT_SECRET", "change-this-secret-key-in-production"),
JWTExpirationHours: getIntEnv("JWT_EXPIRATION_HOURS", 24),
RefreshTokenDays:   getIntEnv("REFRESH_TOKEN_DAYS", 7),
GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
},
Queue: QueueConfig{
Type:              getEnv("QUEUE_TYPE", "redis"),
MaxRetries:        getIntEnv("QUEUE_MAX_RETRIES", 3),
RetryDelaySeconds: getIntEnv("QUEUE_RETRY_DELAY_SECONDS", 60),
DeadLetterQueue:   getEnv("QUEUE_DEAD_LETTER_QUEUE", "dead_letter"),
},
Worker: WorkerConfig{
NumWorkers:      getIntEnv("WORKER_NUM_WORKERS", 10),
MaxJobRetry:     getIntEnv("WORKER_MAX_JOB_RETRY", 3),
JobTimeoutSecs:  getIntEnv("WORKER_JOB_TIMEOUT_SECS", 300),
HealthCheckSecs: getIntEnv("WORKER_HEALTH_CHECK_SECS", 30),
},
Notification: NotificationConfig{
TelegramBotToken:  getEnv("TELEGRAM_BOT_TOKEN", ""),
TelegramChatID:    getEnv("TELEGRAM_CHAT_ID", ""),
DiscordWebhookURL: getEnv("DISCORD_WEBHOOK_URL", ""),
SlackWebhookURL:   getEnv("SLACK_WEBHOOK_URL", ""),
SMTPHost:          getEnv("SMTP_HOST", ""),
SMTPPort:          getIntEnv("SMTP_PORT", 587),
SMTPUser:          getEnv("SMTP_USER", ""),
SMTPPassword:      getEnv("SMTP_PASSWORD", ""),
SMTPFrom:          getEnv("SMTP_FROM", ""),
},
Playwright: PlaywrightConfig{
Headless:    getBoolEnv("PLAYWRIGHT_HEADLESS", true),
BrowserType: getEnv("PLAYWRIGHT_BROWSER_TYPE", "chromium"),
Timeout:     getDurationEnv("PLAYWRIGHT_TIMEOUT", 30*time.Second),
UseProxy:    getBoolEnv("PLAYWRIGHT_USE_PROXY", false),
ProxyURL:    getEnv("PLAYWRIGHT_PROXY_URL", ""),
UserAgent:   getEnv("PLAYWRIGHT_USER_AGENT", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
},
}

return cfg, nil
}

func getEnv(key, defaultValue string) string {
if value := os.Getenv(key); value != "" {
return value
}
return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
if value := os.Getenv(key); value != "" {
if intValue, err := strconv.Atoi(value); err == nil {
return intValue
}
}
return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
if value := os.Getenv(key); value != "" {
if boolValue, err := strconv.ParseBool(value); err == nil {
return boolValue
}
}
return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
if value := os.Getenv(key); value != "" {
if durationValue, err := time.ParseDuration(value); err == nil {
return durationValue
}
}
return defaultValue
}
