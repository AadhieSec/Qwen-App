package main

import (
"context"
"fmt"
"log"
"net/http"
"os"
"os/signal"
"syscall"
"time"

"github.com/gofiber/fiber/v2"
"github.com/gofiber/fiber/v2/middleware/cors"
"github.com/gofiber/fiber/v2/middleware/logger"
"github.com/gofiber/fiber/v2/middleware/recover"
"github.com/gofiber/fiber/v2/middleware/requestid"

"shopmonitor/internal/config"
"shopmonitor/internal/repository"
)

func main() {
// Load configuration
cfg, err := config.Load()
if err != nil {
log.Fatalf("Failed to load configuration: %v", err)
}

// Initialize database
db, err := repository.NewDatabase(&cfg.Database)
if err != nil {
log.Fatalf("Failed to connect to database: %v", err)
}
defer db.Close()

// Run migrations
if err := db.AutoMigrate(); err != nil {
log.Fatalf("Failed to run migrations: %v", err)
}

// Create Fiber app
app := fiber.New(fiber.Config{
AppName:      "ShopMonitor API v1.0",
ReadTimeout:  cfg.Server.ReadTimeout,
WriteTimeout: cfg.Server.WriteTimeout,
IdleTimeout:  cfg.Server.IdleTimeout,
ErrorHandler: customErrorHandler,
})

// Middleware
app.Use(requestid.New())
app.Use(logger.New())
app.Use(recover.New())
app.Use(cors.New(cors.Config{
AllowOrigins: "*",
AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
AllowHeaders: "Origin,Content-Type,Accept,Authorization",
}))

// Health check endpoint
app.Get("/health", func(c *fiber.Ctx) error {
return c.JSON(fiber.Map{
"status":    "healthy",
"timestamp": time.Now(),
})
})

// API routes would be registered here
// registerRoutes(app, db, cfg)

// Start server in goroutine
go func() {
addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
log.Printf("Server starting on %s", addr)
if err := app.Listen(addr); err != nil {
log.Fatalf("Failed to start server: %v", err)
}
}()

// Graceful shutdown
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

log.Println("Shutting down server...")

ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownDelay)
defer cancel()

if err := app.ShutdownWithContext(ctx); err != nil {
log.Fatalf("Server forced to shutdown: %v", err)
}

log.Println("Server exited gracefully")
}

func customErrorHandler(c *fiber.Ctx, err error) error {
code := fiber.StatusInternalServerError
if e, ok := err.(*fiber.Error); ok {
code = e.Code
}

c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
return c.Status(code).JSON(fiber.Map{
"error": map[string]interface{}{
"code":    code,
"message": err.Error(),
},
})
}
