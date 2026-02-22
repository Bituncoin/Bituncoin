import global.genesis.httpclient.GenesisHttpClient
import kotlinx.serialization.json.*

// BTNG Price Bridge Service
// Polls the BTNG Node.js API every minute and stores gold prices in Genesis database

job("BTNG_PRICE_POLLER") {
    val client = GenesisHttpClient()
    val loginUrl = "http://localhost:3003/api/auth/login"
    val priceUrl = "http://localhost:3003/api/btng/gold/price/status"

    // Schedule: Runs every 60 seconds
    schedule("0 * * * * ?") {
        try {
            // 1. Authenticate with BTNG API to get JWT token
            val loginResponse: String = client.post(
                loginUrl,
                mapOf(
                    "username" to "admin",
                    "password" to "password123"
                )
            )
            val loginJson = Json.parseToJsonElement(loginResponse).jsonObject
            val token = loginJson["token"]?.jsonPrimitive?.content

            if (token.isNullOrEmpty()) {
                LOG.error("BTNG Bridge: Failed to authenticate - no token received")
                return@schedule
            }

            // 2. Fetch gold price data with authentication
            val priceResponse: String = client.get(
                priceUrl,
                headers = mapOf("Authorization" to "Bearer $token")
            )
            val priceJson = Json.parseToJsonElement(priceResponse).jsonObject

            // 3. Extract price data from the response
            val latestPrice = priceJson["latest_price"]?.jsonObject
            val basePriceGram = latestPrice?.get("base_price_gram")?.jsonPrimitive?.double ?: 0.0
            val basePriceOunce = latestPrice?.get("base_price_ounce")?.jsonPrimitive?.double ?: 0.0
            val basePriceKilo = latestPrice?.get("base_price_kilo")?.jsonPrimitive?.double ?: 0.0

            // 4. Upsert into Genesis Database
            db.upsert(
                GoldPriceCache {
                    instrument = "GOLD/USD"
                    price = basePriceOunce  // Using ounce price as primary
                    basePriceGram = basePriceGram
                    basePriceOunce = basePriceOunce
                    basePriceKilo = basePriceKilo
                    lastUpdate = DateTime.now()
                    source = "BTNG_API"
                }
            )

            LOG.info("BTNG Bridge: Successfully synced gold price - Gram: $basePriceGram, Ounce: $basePriceOunce, Kilo: $basePriceKilo")

        } catch (e: Exception) {
            LOG.error("BTNG Bridge Error: ${e.message}")
            LOG.error("BTNG Bridge StackTrace: ${e.stackTraceToString()}")
        }
    }
}

// Health check job - runs every 5 minutes to ensure bridge is working
job("BTNG_BRIDGE_HEALTH_CHECK") {
    schedule("0 */5 * * * ?") {
        try {
            val client = GenesisHttpClient()
            val loginUrl = "http://localhost:3003/api/auth/login"

            // Test authentication
            val loginResponse: String = client.post(
                loginUrl,
                mapOf(
                    "username" to "admin",
                    "password" to "password123"
                )
            )

            val loginJson = Json.parseToJsonElement(loginResponse).jsonObject
            val token = loginJson["token"]?.jsonPrimitive?.content

            if (!token.isNullOrEmpty()) {
                LOG.info("BTNG Bridge Health: ✅ Authentication successful")
            } else {
                LOG.warn("BTNG Bridge Health: ❌ Authentication failed - no token")
            }

        } catch (e: Exception) {
            LOG.error("BTNG Bridge Health: ❌ Connection failed - ${e.message}")
        }
    }
}