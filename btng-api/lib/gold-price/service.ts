import axios from 'axios';

const GOLDAPI_KEY = "goldapi-3saduazsmlu7hhzo-io";
const GOLDAPI_URL = "https://www.goldapi.io/api/XAU/USD";
const FX_URL = "https://api.exchangerate.host/latest?base=USD";

export interface GoldPriceData {
  base_currency: string;
  base_price_gram: number;
  base_price_ounce: number;
  base_price_kilo: number;
  currencies: Array<{
    currency: string;
    price_gram: number;
    price_ounce: number;
    price_kilo: number;
  }>;
  fx_rates: Record<string, number>;
  bid: number;
  ask: number;
  spread: number | null;
  timestamp: number;
}

export async function fetchGoldPriceFromAPI(): Promise<Partial<GoldPriceData>> {
  try {
    const response = await axios.get(GOLDAPI_URL, {
      headers: { "x-access-token": GOLDAPI_KEY },
      timeout: 10000 // 10 second timeout
    });

    const data = response.data;

    // GoldAPI returns price per troy ounce and price per gram (24K)
    const ounce = data.price;
    const gram = data.price_gram_24K;
    const kilo = gram * 1000;

    return {
      base_currency: "USD",
      base_price_gram: gram,
      base_price_ounce: ounce,
      base_price_kilo: kilo,
      bid: data.bid || null,
      ask: data.ask || null,
      spread: data.ask && data.bid ? data.ask - data.bid : null,
      timestamp: data.timestamp || Date.now()
    };

  } catch (error: any) {
    console.error('GoldAPI fetch error:', error.message);
    throw new Error(`Failed to fetch gold price: ${error.message}`);
  }
}

export async function fetchFxRates(): Promise<Record<string, number>> {
  try {
    const response = await axios.get(FX_URL, {
      timeout: 5000
    });
    return response.data.rates;
  } catch (error: any) {
    console.error('FX rates fetch error:', error.message);
    // Return fallback rates if FX API fails
    return {
      USD: 1,
      EUR: 0.92,
      GBP: 0.79,
      GHS: 13.77,
      NGN: 1550,
      AED: 3.67,
      SAR: 3.75,
      KWD: 0.31,
      EGP: 47.5
    };
  }
}

export async function buildMultiCurrencyGoldPrice(): Promise<GoldPriceData> {
  const goldData = await fetchGoldPriceFromAPI();
  const fxRates = await fetchFxRates();

  // Target currencies for BTNG sovereign pricing
  const targetCurrencies = ["USD", "EUR", "GBP", "GHS", "NGN", "AED", "SAR", "KWD", "EGP"];

  const currencies = targetCurrencies.map(code => {
    const rate = fxRates[code];
    if (!rate) {
      console.warn(`No FX rate available for ${code}`);
      return null;
    }

    return {
      currency: code,
      price_gram: goldData.base_price_gram! * rate,
      price_ounce: goldData.base_price_ounce! * rate,
      price_kilo: goldData.base_price_kilo! * rate
    };
  }).filter(Boolean) as GoldPriceData['currencies'];

  return {
    base_currency: goldData.base_currency!,
    base_price_gram: goldData.base_price_gram!,
    base_price_ounce: goldData.base_price_ounce!,
    base_price_kilo: goldData.base_price_kilo!,
    currencies,
    fx_rates: fxRates,
    bid: goldData.bid!,
    ask: goldData.ask!,
    spread: goldData.spread!,
    timestamp: goldData.timestamp!
  };
}