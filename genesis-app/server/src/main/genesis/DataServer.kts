// Data Server Configuration for BTNG Gold Price Streaming
// This enables real-time updates to frontend components

dataServer {
    // Stream all gold price updates to subscribed clients
    query("ALL_GOLD_PRICES") {
        source {
            GOLD_PRICE_CACHE
        }

        // Send updates when any gold price record changes
        onCommit = true

        // Include all fields in the stream
        fields {
            INSTRUMENT
            PRICE
            BASE_PRICE_GRAM
            BASE_PRICE_OUNCE
            BASE_PRICE_KILO
            LAST_UPDATE
            SOURCE
            CURRENCIES
            FX_RATES
            BID
            ASK
            SPREAD
        }
    }

    // Stream only the latest gold price (most recent update)
    query("LATEST_GOLD_PRICE") {
        source {
            GOLD_PRICE_CACHE
        }

        // Filter to get only the most recent record
        where {
            LAST_UPDATE maxOf LAST_UPDATE
        }

        onCommit = true

        fields {
            INSTRUMENT
            PRICE
            BASE_PRICE_GRAM
            BASE_PRICE_OUNCE
            BASE_PRICE_KILO
            LAST_UPDATE
            SOURCE
        }
    }

    // Stream gold price history (last 24 hours)
    query("GOLD_PRICE_HISTORY_24H") {
        source {
            GOLD_PRICE_CACHE
        }

        where {
            LAST_UPDATE >= DateTime.now().minusHours(24)
        }

        orderBy {
            LAST_UPDATE descending
        }

        onCommit = true

        fields {
            INSTRUMENT
            PRICE
            LAST_UPDATE
            SOURCE
        }
    }
}