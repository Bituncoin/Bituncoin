import mongoose from 'mongoose';

const MONGODB_URI = process.env.MONGODB_URI || 'mongodb://localhost:27017/btng-sovereign';

const CurrencyPriceSchema = new mongoose.Schema({
  currency: { type: String, required: true }, // e.g., "USD", "GHS", "NGN"
  price_gram: { type: Number, required: true },
  price_ounce: { type: Number, required: true },
  price_kilo: { type: Number, required: true }
}, { _id: false });

const GoldPriceSchema = new mongoose.Schema({
  base_currency: { type: String, default: 'USD' },
  base_price_gram: { type: Number, required: true },
  base_price_ounce: { type: Number, required: true },
  base_price_kilo: { type: Number, required: true },

  currencies: [CurrencyPriceSchema], // all converted prices

  fx_rates: { type: Object, default: {} }, // { GHS: 12.3, NGN: 1450, EUR: 0.92 }

  bid: { type: Number, default: null },
  ask: { type: Number, default: null },
  spread: { type: Number, default: null },
  timestamp: { type: Number, default: Date.now }
}, {
  timestamps: true // adds createdAt and updatedAt
});

export const GoldPrice = mongoose.models.GoldPrice || mongoose.model('GoldPrice', GoldPriceSchema);

// In-memory cache for latest price and fallback storage
let latestGoldPrice: any = null;
let inMemoryStorage: any[] = [];
let useInMemoryFallback = false;

async function ensureConnection() {
  try {
    await mongoose.connect(MONGODB_URI);
    useInMemoryFallback = false;
  } catch (error: any) {
    console.warn('⚠️ MongoDB connection failed, using in-memory storage:', error?.message || error);
    useInMemoryFallback = true;
  }
}

export async function saveGoldPrice(data: any) {
  await ensureConnection();

  if (useInMemoryFallback) {
    // Fallback to in-memory storage
    const entry = { ...data, _id: Date.now().toString(), createdAt: new Date(), updatedAt: new Date() };
    inMemoryStorage.push(entry);
    // Keep only last 1000 entries
    if (inMemoryStorage.length > 1000) {
      inMemoryStorage = inMemoryStorage.slice(-1000);
    }
    console.log('💾 Saved gold price to in-memory storage');
    return entry;
  } else {
    // Use MongoDB
    try {
      const entry = new GoldPrice(data);
      const saved = await entry.save();
      console.log('💾 Saved gold price to MongoDB');
      return saved;
    } catch (error: any) {
      console.warn('⚠️ MongoDB save failed, falling back to in-memory:', error?.message || error);
      // Fallback to in-memory if save fails
      const entry = { ...data, _id: Date.now().toString(), createdAt: new Date(), updatedAt: new Date() };
      inMemoryStorage.push(entry);
      return entry;
    }
  }
}

export async function setLatestGoldPrice(data: any) {
  latestGoldPrice = data;
}

export function getLatestGoldPrice() {
  return latestGoldPrice;
}

export async function getGoldPriceHistory(limit: number = 100) {
  await ensureConnection();

  if (useInMemoryFallback) {
    // Return from in-memory storage
    return inMemoryStorage
      .sort((a, b) => b.timestamp - a.timestamp)
      .slice(0, limit);
  } else {
    // Use MongoDB
    return await GoldPrice.find()
      .sort({ timestamp: -1 })
      .limit(limit)
      .lean();
  }
}

export async function getGoldPriceByTimestampRange(startTime: number, endTime: number) {
  await ensureConnection();

  if (useInMemoryFallback) {
    // Filter from in-memory storage
    return inMemoryStorage
      .filter(item => item.timestamp >= startTime && item.timestamp <= endTime)
      .sort((a, b) => a.timestamp - b.timestamp);
  } else {
    // Use MongoDB
    return await GoldPrice.find({
      timestamp: { $gte: startTime, $lte: endTime }
    })
    .sort({ timestamp: 1 })
    .lean();
  }
}