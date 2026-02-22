package global.genesis

import global.genesis.system.definition.*

// BTNG Sovereign Gold Standard - Genesis System Definition
// This file defines the Genesis platform configuration for BTNG integration

systemDefinition {
    // Global system settings
    global {
        item(name = "ENVIRONMENT", value = "development")
        item(name = "LOG_LEVEL", value = "INFO")
        item(name = "TIMEZONE", value = "UTC")
    }

    // BTNG Smart Contract Configuration
    item(name = "BTNG_RPC_URL", value = "https://eth-sepolia.g.alchemy.com/v2/YOUR_ALCHEMY_KEY")
    item(name = "BTNG_CHAIN_ID", value = "11155111")
    item(name = "BTNG_TOKEN_ADDRESS", value = "0x0000000000000000000000000000000000000000")
    item(name = "BTNG_ORACLE_ADDRESS", value = "0x0000000000000000000000000000000000000000")

    // Database configuration
    database {
        // Use FoundationDB or PostgreSQL
        item(name = "DB_TYPE", value = "POSTGRESQL")
        item(name = "DB_HOST", value = "localhost")
        item(name = "DB_PORT", value = "5432")
        item(name = "DB_NAME", value = "btng_genesis")
        item(name = "DB_USERNAME", value = "btng")
        item(name = "DB_PASSWORD", value = "btng_password")
    }

    // Web server configuration
    web {
        item(name = "PORT", value = "8080")
        item(name = "HOST", value = "0.0.0.0")
        item(name = "ENABLE_CORS", value = "true")
        item(name = "CORS_ORIGINS", value = "http://localhost:3003,http://localhost:3000")
    }

    // Authentication configuration
    auth {
        item(name = "JWT_SECRET", value = "btng-genesis-jwt-secret-key-256-bits")
        item(name = "SESSION_TIMEOUT", value = "3600000") // 1 hour
    }

    // BTNG Integration settings
    btng {
        item(name = "API_BASE_URL", value = "http://localhost:3003")
        item(name = "PRICE_POLL_INTERVAL", value = "10000") // 10 seconds
        item(name = "ENABLE_REAL_TIME_STREAMING", value = "true")
    }

    // Process configuration
    processes {
        // Gold price poller process
        process(name = "GOLD_PRICE_POLLER") {
            item(name = "ENABLED", value = "true")
            item(name = "SCHEDULE", value = "0/10 * * * * ?") // Every 10 seconds
        }

        // Health check process
        process(name = "HEALTH_CHECK") {
            item(name = "ENABLED", value = "true")
            item(name = "SCHEDULE", value = "0 */5 * * * ?") // Every 5 minutes
        }
    }

    // Data server configuration for real-time streaming
    dataServer {
        item(name = "ENABLED", value = "true")
        item(name = "PORT", value = "9000")
        item(name = "WEBSOCKET_ENABLED", value = "true")
    }

    // Monitoring and metrics
    monitoring {
        item(name = "ENABLED", value = "true")
        item(name = "METRICS_PORT", value = "9090")
        item(name = "HEALTH_CHECK_ENDPOINT", value = "/api/health")
    }
}