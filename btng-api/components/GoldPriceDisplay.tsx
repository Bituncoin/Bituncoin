'use client';

import React from 'react';
import { useGoldPrice } from '@/hooks/useGoldPrice';

export default function GoldPriceDisplay() {
  const { price, isConnected, error } = useGoldPrice();

  if (error) {
    return (
      <div className="gold-price-error">
        <h3>⚠️ Gold Price Error</h3>
        <p>{error}</p>
        <small>Connection: {isConnected ? '✅' : '❌'}</small>
      </div>
    );
  }

  if (!price) {
    return (
      <div className="gold-price-loading">
        <h3>📊 BTNG Gold Price</h3>
        <p>Loading price data...</p>
        <small>Connection: {isConnected ? '✅' : '❌'}</small>
      </div>
    );
  }

  return (
    <div className="gold-price-display">
      <h3>🥇 BTNG Sovereign Gold Price</h3>

      <div className="price-grid">
        <div className="price-item">
          <span className="label">Per Gram (24K)</span>
          <span className="value">${price.base_price_gram?.toFixed(2)}</span>
        </div>

        <div className="price-item">
          <span className="label">Per Ounce</span>
          <span className="value">${price.base_price_ounce?.toFixed(2)}</span>
        </div>

        <div className="price-item">
          <span className="label">Per Kilo</span>
          <span className="value">${price.base_price_kilo?.toFixed(0)}</span>
        </div>
      </div>

      <div className="price-meta">
        <small>
          Last updated: {price.timestamp ? new Date(price.timestamp).toLocaleString() : 'Unknown'}
        </small>
        <br />
        <small>
          Currencies: {price.currencies_available || 0} available
        </small>
        <br />
        <small>
          Status: {isConnected ? '🟢 Live' : '🔴 Offline'}
        </small>
      </div>

      <style jsx>{`
        .gold-price-display {
          border: 2px solid #FFD700;
          border-radius: 12px;
          padding: 20px;
          background: linear-gradient(135deg, #1a1a1a 0%, #2a2a2a 100%);
          color: white;
          font-family: 'Courier New', monospace;
          max-width: 400px;
          margin: 20px auto;
        }

        .gold-price-display h3 {
          margin: 0 0 15px 0;
          color: #FFD700;
          text-align: center;
          font-size: 1.2em;
        }

        .price-grid {
          display: grid;
          gap: 10px;
          margin-bottom: 15px;
        }

        .price-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 8px 12px;
          background: rgba(255, 215, 0, 0.1);
          border-radius: 6px;
        }

        .label {
          font-size: 0.9em;
          opacity: 0.8;
        }

        .value {
          font-weight: bold;
          font-size: 1.1em;
          color: #FFD700;
        }

        .price-meta {
          text-align: center;
          opacity: 0.7;
          font-size: 0.8em;
        }

        .gold-price-loading, .gold-price-error {
          border: 2px solid #666;
          border-radius: 12px;
          padding: 20px;
          background: #1a1a1a;
          color: white;
          text-align: center;
          max-width: 400px;
          margin: 20px auto;
        }

        .gold-price-error {
          border-color: #ff6b6b;
        }

        .gold-price-error h3 {
          color: #ff6b6b;
          margin: 0 0 10px 0;
        }
      `}</style>
    </div>
  );
}