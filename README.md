# GoUtils

A comprehensive Go utilities library providing essential components for building scalable Go applications. This library includes utilities for logging, configuration management, caching, database connections, event queues, and more.

## Features

- 🔧 **HTTP Router**: Pre-configured Gin router with essential middleware
- 📝 **Structured Logging**: Advanced logging with zap integration and middleware support
- ⚙️ **Configuration Management**: Support for multiple config sources (file, Consul, Zookeeper)
- 🗃️ **Caching**: Redis and in-memory cache implementations
- 🗄️ **Database Connections**: MongoDB, MySQL, PostgreSQL connection utilities
- 📨 **Event Queue**: Kafka integration for event-driven architecture
- 🛠️ **Common Functions**: Email/mobile validation and utility functions
- 🔄 **Application Lifecycle**: Application interface with graceful shutdown
- 🔍 **Utilities**: Generic utility functions and helpers

## Installation

```bash
go get github.com/gofreego/goutils
```

## Quick Start

### HTTP Router

```go
package main

import (
    "github.com/gofreego/goutils"
    "github.com/gin-gonic/gin"
)

func main() {
    router := goutils.GetHTTPRouter(gin.ReleaseMode)
    
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    
    router.Run(":8080")
}
```

### Logging

```go
package main

import (
    "context"
    "github.com/gofreego/goutils/logger"
)

func main() {
    // Initialize logger
    config := logger.Config{
        AppName: "my-app",
        Build:   "v1.0.0",
    }
    config.InitiateLogger()
    
    ctx := context.Background()
    
    // Simple logging
    logger.Info(ctx, "Application started")
    
    // Structured logging
    fields := &logger.Fields{}
    fields.AddField("user_id", "12345")
    fields.AddField("action", "login")
    logger.Infof(ctx, "User action performed", fields)
}
```

### Configuration Management

```go
package main

import (
    "context"
    "github.com/gofreego/goutils/configutils"
)

type AppConfig struct {
    Name   string `yaml:"name"`
    Port   int    `yaml:"port"`
    Reader configutils.Config `yaml:"reader"`
}

func (c *AppConfig) GetReaderConfig() *configutils.Config {
    return &c.Reader
}

func main() {
    var config AppConfig
    err := configutils.ReadConfig(context.Background(), "config.yaml", &config)
    if err != nil {
        panic(err)
    }
    
    configutils.LogConfig(context.Background(), &config)
}
```

### Caching

```go
package main

import (
    "context"
    "time"
    "github.com/gofreego/goutils/cache"
)

func main() {
    ctx := context.Background()
    
    // Redis cache
    redisConfig := &cache.Config{
        Name: cache.REDIS,
        Redis: redis.Config{
            Host:     "localhost:6379",
            Password: "",
            DB:       0,
        },
    }
    
    redisCache := cache.NewCache(ctx, redisConfig)
    
    // Set value
    redisCache.Set(ctx, "key", "value")
    
    // Set with timeout
    redisCache.SetWithTimeout(ctx, "temp_key", "temp_value", 5*time.Minute)
    
    // Get value
    var result string
    redisCache.GetV(ctx, "key", &result)
}
```

### Database Connections

#### MongoDB

```go
package main

import (
    "context"
    "time"
    "github.com/gofreego/goutils/databases/connections/mongo"
)

func main() {
    config := &mongo.Config{
        Hosts:                  "localhost:27017",
        Username:               "admin",
        Password:               "password",
        Database:               "mydb",
        MaxPoolSize:            100,
        MinPoolSize:            10,
        ConnectTimeout:         10 * time.Second,
        ServerSelectionTimeout: 5 * time.Second,
    }
    
    client, err := mongo.GetClient(context.Background(), config)
    if err != nil {
        panic(err)
    }
    
    database := client.Database(config.Database)
    // Use database...
}
```

### Event Queue (Kafka)

```go
package main

import (
    "context"
    "github.com/gofreego/goutils/eventqueue"
    "github.com/gofreego/goutils/eventqueue/kafka"
    "github.com/gofreego/goutils/eventqueue/models"
)

func main() {
    config := &kafka.Config{
        Brokers: []string{"localhost:9092"},
        Topic:   "my-topic",
        GroupID: "my-group",
    }
    
    queue := kafka.NewKafkaEventQueue(context.Background(), config)
    
    // Publish message
    message := &models.Message{
        Topic: "my-topic",
        Key:   "event-key",
        Value: []byte("event data"),
    }
    
    err := queue.Publish(context.Background(), message)
    if err != nil {
        panic(err)
    }
    
    // Consume message
    consumedMessage, err := queue.Consume(context.Background())
    if err != nil {
        panic(err)
    }
    
    // Process message and commit
    // ... process consumedMessage ...
    queue.Commit(context.Background(), consumedMessage)
}
```

## Modules Overview

### API
- **Router**: HTTP router with middleware for CORS, request timing, and request ID tracking

### Logger
- **Structured Logging**: Built on zap logger with context support
- **Middleware Support**: Extensible logging middleware system
- **Multiple Levels**: Support for Info, Error, Warn, Debug levels

### ConfigUtils
- **Multiple Sources**: File, Consul, Zookeeper support
- **Format Support**: JSON and YAML configuration formats
- **Interface-based**: Clean abstraction for different config readers

### Cache
- **Redis**: Full Redis integration with connection pooling
- **Memory**: In-memory cache for development and testing
- **Timeout Support**: TTL-based cache entries

### Databases
- **MongoDB**: Connection management with advanced pool configuration
- **SQL Databases**: MySQL and PostgreSQL connection utilities
- **Migration Support**: Database migration utilities

### EventQueue
- **Kafka**: Producer and consumer implementations
- **Message Interface**: Generic message handling
- **Commit Support**: Manual and automatic message commitment

### Utils
- **Common Functions**: Email and mobile validation
- **Generic Utilities**: Type-safe utility functions
- **Application Interface**: Standard application lifecycle management

## Configuration Examples

### File-based Configuration (config.yaml)
```yaml
name: "my-application"
port: 8080
reader:
  Name: "file"
  Format: "yaml"
  File:
    Path: "./config"
    Filename: "app.yaml"
```

### Consul Configuration
```yaml
reader:
  Name: "consul"
  Format: "json"
  Consul:
    Address: "localhost:8500"
    Scheme: "http"
    Datacenter: "dc1"
```

## Middleware

The library includes several built-in middleware components:

- **Request ID Middleware**: Adds unique request IDs to all requests
- **Request Time Middleware**: Logs request duration
- **CORS Middleware**: Handles cross-origin requests
- **Logging Middleware**: Structured request/response logging

## Contributing

We welcome contributions! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Requirements

- Go 1.23.0 or higher

## Dependencies

Key dependencies include:
- Gin Web Framework
- Zap Logger
- Redis Client
- MongoDB Driver
- Kafka Client
- And more... (see go.mod for full list)

---

For more examples and detailed documentation, please check the individual package documentation and example files in the repository.
