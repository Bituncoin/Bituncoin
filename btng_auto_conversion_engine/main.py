"""
BTNG Auto-Conversion Engine
Sovereign Currency-to-Gold conversion.
"""

import time
import os

import requests
from fastapi import FastAPI
from pydantic import BaseModel, Field

app = FastAPI(title="BTNG Auto-Conversion Engine")

# Gold Oracle endpoint
ORACLE_URL = os.getenv("BTNG_ORACLE_URL", "http://127.0.0.1:38991/oracle/price/live")
FALLBACK_GOLD_PRICE_USD_GRAM = 165.23

# Mock USD conversion rates
FX_RATES = {
    "ghs": 13.50,
    "btc": 65000.0,
    "eth": 3500.0,
    "usdt": 1.0,
    "btng": 1.0,
}


def get_live_gold_price():
    """
    Fetches the live gold price from the BTNG Oracle.
    Using port 38991 (Gold Oracle).
    """
    try:
        response = requests.get(ORACLE_URL, timeout=2)
        if response.status_code == 200:
            return response.json()["p_gold_usd_gram"]
        return FALLBACK_GOLD_PRICE_USD_GRAM
    except (requests.RequestException, ValueError, KeyError, TypeError):
        return FALLBACK_GOLD_PRICE_USD_GRAM


def convert_to_usd(currency, amount):
    """
    Converts incoming currency amount to USD.
    """
    c = currency.strip().lower()
    if c == "usd":
        return amount
    if c in FX_RATES:
        if c in ["btc", "eth", "usdt", "btng"]:
            return amount * FX_RATES[c]
        if c == "ghs":
            return amount / FX_RATES[c]
    return amount


class ConversionRequest(BaseModel):
    currency: str = Field(min_length=1, max_length=10)
    amount: float = Field(gt=0)


class ConversionResponse(BaseModel):
    grams: float
    currency: str
    amount: float
    status: str
    timestamp: int


@app.post("/convert/to-btng", response_model=ConversionResponse)
def convert_to_btng_api(req: ConversionRequest):
    """
    Unified entry point for Auto-Conversion to BTNG gold-grams.
    """
    gold_price_usd_gram = get_live_gold_price()
    usd_value = convert_to_usd(req.currency, req.amount)
    grams = usd_value / gold_price_usd_gram

    return ConversionResponse(
        grams=round(grams, 6),
        currency=req.currency,
        amount=req.amount,
        status="ACTIVE",
        timestamp=int(time.time()),
    )


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        app,
        host=os.getenv("BTNG_HOST", "127.0.0.1"),
        port=int(os.getenv("BTNG_PORT", "38994")),
    )
