import { useEffect, useState } from 'react';

interface GoldPriceData {
  price?: number;
  timestamp?: string;
  base_price_gram?: number;
  base_price_ounce?: number;
  base_price_kilo?: number;
  currencies_available?: number;
}

export function useGoldPrice() {
  const [price, setPrice] = useState<GoldPriceData | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let reconnectTimeout: NodeJS.Timeout;

    const connect = () => {
      try {
        // For now, we'll poll the REST API since WebSocket isn't implemented yet
        // TODO: Replace with actual WebSocket when /ws/btng/gold-price is available
        const pollPrice = async () => {
          try {
            const response = await fetch('/api/btng/gold/price/status');
            if (response.ok) {
              const data = await response.json();
              if (data.latest_price) {
                setPrice({
                  price: data.latest_price.base_price_ounce,
                  timestamp: data.latest_price.timestamp,
                  base_price_gram: data.latest_price.base_price_gram,
                  base_price_ounce: data.latest_price.base_price_ounce,
                  base_price_kilo: data.latest_price.base_price_kilo,
                  currencies_available: data.latest_price.currencies_available
                });
                setIsConnected(true);
                setError(null);
              }
            } else {
              throw new Error(`HTTP ${response.status}`);
            }
          } catch (err: any) {
            setError(`Failed to fetch price: ${err.message}`);
            setIsConnected(false);
          }
        };

        // Initial fetch
        pollPrice();

        // Poll every 10 seconds
        const interval = setInterval(pollPrice, 10000);

        return () => clearInterval(interval);

      } catch (err: any) {
        setError(`Connection failed: ${err.message}`);
        setIsConnected(false);

        // Retry connection after 5 seconds
        reconnectTimeout = setTimeout(connect, 5000);
      }
    };

    connect();

    return () => {
      if (reconnectTimeout) clearTimeout(reconnectTimeout);
    };
  }, []);

  return { price, isConnected, error };
}

// Alternative polling-only version (simpler)
export function useGoldPricePolling(intervalMs: number = 10000) {
  const [price, setPrice] = useState<GoldPriceData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPrice = async () => {
      try {
        setLoading(true);
        const response = await fetch('/api/btng/gold/price/status');
        if (!response.ok) throw new Error(`HTTP ${response.status}`);

        const data = await response.json();
        if (data.latest_price) {
          setPrice({
            price: data.latest_price.base_price_ounce,
            timestamp: data.latest_price.timestamp,
            base_price_gram: data.latest_price.base_price_gram,
            base_price_ounce: data.latest_price.base_price_ounce,
            base_price_kilo: data.latest_price.base_price_kilo,
            currencies_available: data.latest_price.currencies_available
          });
          setError(null);
        } else {
          setPrice(null);
          setError('No price data available');
        }
      } catch (err: any) {
        setError(err.message);
        setPrice(null);
      } finally {
        setLoading(false);
      }
    };

    // Initial fetch
    fetchPrice();

    // Set up polling
    const interval = setInterval(fetchPrice, intervalMs);

    return () => clearInterval(interval);
  }, [intervalMs]);

  return { price, loading, error };
}