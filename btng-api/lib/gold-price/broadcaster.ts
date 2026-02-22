import { buildMultiCurrencyGoldPrice } from '@/lib/gold-price/service';
import { saveGoldPrice, setLatestGoldPrice } from '@/lib/gold-price/model';

class GoldPriceBroadcaster {
  private intervalId: NodeJS.Timeout | null = null;
  private isRunning = false;
  private updateInterval = 10000; // 10 seconds

  async start() {
    if (this.isRunning) {
      console.log('BTNG Gold Price Broadcaster is already running');
      return;
    }

    console.log('Starting BTNG Gold Price Broadcaster...');
    this.isRunning = true;

    // Initial broadcast
    await this.broadcast();

    // Set up continuous broadcasting
    this.intervalId = setInterval(async () => {
      await this.broadcast();
    }, this.updateInterval);

    console.log(`BTNG Gold Price Broadcaster started - updating every ${this.updateInterval / 1000} seconds`);
  }

  async stop() {
    if (this.intervalId) {
      clearInterval(this.intervalId);
      this.intervalId = null;
    }
    this.isRunning = false;
    console.log('BTNG Gold Price Broadcaster stopped');
  }

  private async broadcast() {
    try {
      const goldPriceData = await buildMultiCurrencyGoldPrice();

      // Save to database
      await saveGoldPrice(goldPriceData);

      // Update in-memory cache
      await setLatestGoldPrice(goldPriceData);

      console.log('✅ BTNG Gold Price Broadcast:', {
        timestamp: new Date(goldPriceData.timestamp).toISOString(),
        base_price_gram: goldPriceData.base_price_gram.toFixed(2),
        currencies: goldPriceData.currencies.length,
        fx_rates_available: Object.keys(goldPriceData.fx_rates).length
      });

    } catch (error: any) {
      console.error('❌ BTNG Gold Price Broadcast Error:', error.message);

      // Continue running even if one update fails
      // The next interval will try again
    }
  }

  getStatus() {
    return {
      isRunning: this.isRunning,
      updateInterval: this.updateInterval,
      lastUpdate: this.isRunning ? new Date().toISOString() : null
    };
  }
}

// Singleton instance
export const goldPriceBroadcaster = new GoldPriceBroadcaster();

// Auto-start in production/development
if (typeof window === 'undefined') { // Server-side only
  // Start broadcaster when module is loaded
  goldPriceBroadcaster.start().catch(console.error);
}