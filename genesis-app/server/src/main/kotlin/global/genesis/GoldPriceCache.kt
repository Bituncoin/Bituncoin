package global.genesis

import global.genesis.db.TableDefinition
import global.genesis.db.FieldType

// Gold Price Cache Table Definition
// Stores BTNG gold price data polled from the Node.js API

val GOLD_PRICE_CACHE = TableDefinition {
    tableName = "GOLD_PRICE_CACHE"

    // Primary key
    field("INSTRUMENT", FieldType.STRING, 50) {
        primaryKey = true
    }

    // Price data
    field("PRICE", FieldType.DOUBLE)
    field("BASE_PRICE_GRAM", FieldType.DOUBLE)
    field("BASE_PRICE_OUNCE", FieldType.DOUBLE)
    field("BASE_PRICE_KILO", FieldType.DOUBLE)

    // Metadata
    field("LAST_UPDATE", FieldType.DATETIME)
    field("SOURCE", FieldType.STRING, 50)

    // Optional fields for future expansion
    field("CURRENCIES", FieldType.STRING, 1000) // JSON string of currency data
    field("FX_RATES", FieldType.STRING, 1000)   // JSON string of FX rates
    field("BID", FieldType.DOUBLE)
    field("ASK", FieldType.DOUBLE)
    field("SPREAD", FieldType.DOUBLE)

    // Indexing for performance
    index("LAST_UPDATE_INDEX") {
        field("LAST_UPDATE")
    }

    index("SOURCE_INDEX") {
        field("SOURCE")
    }
}