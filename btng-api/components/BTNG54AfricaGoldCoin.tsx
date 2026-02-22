'use client';

import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import styles from './BTNG54AfricaGoldCoin.module.css';

interface Country {
  name: string;
  code: string;
  flag: string;
  goldReserves: number;
  population: number;
  sovereignValue: number;
}

interface GoldCoinData {
  countries: Country[];
  totalGoldValue: number;
  totalCountries: number;
  totalPopulation: number;
  totalGoldReserves: number;
  lastUpdated: string;
}

interface BTNG54AfricaGoldCoinProps {
  className?: string;
}

const BTNG54AfricaGoldCoin: React.FC<BTNG54AfricaGoldCoinProps> = ({ className }) => {
  const [selectedCountry, setSelectedCountry] = useState<Country | null>(null);
  const [goldCoinData, setGoldCoinData] = useState<GoldCoinData | null>(null);
  const [isAnimating, setIsAnimating] = useState(false);
  const [currentView, setCurrentView] = useState<'unified' | 'countries' | 'sovereign'>('unified');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Fetch gold coin data from API
  useEffect(() => {
    const fetchGoldCoinData = async () => {
      try {
        const response = await fetch('/api/gold-coin');
        const result = await response.json();

        if (result.success) {
          setGoldCoinData(result.data);
        } else {
          setError('Failed to load gold coin data');
        }
      } catch (err) {
        setError('Network error while loading data');
      } finally {
        setLoading(false);
      }
    };

    fetchGoldCoinData();
  }, []);

  if (loading) {
    return (
      <div className={`${styles.btng54AfricaGoldCoin} ${className || ''}`}>
        <div className={styles.loading}>
          <motion.div
            className={styles.loadingSpinner}
            animate={{ rotate: 360 }}
            transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
          />
          <p>Loading BTNG 54 Africa Gold Coin...</p>
        </div>
      </div>
    );
  }

  if (error || !goldCoinData) {
    return (
      <div className={`${styles.btng54AfricaGoldCoin} ${className || ''}`}>
        <div className={styles.error}>
          <h2>⚠️ Error Loading Gold Coin Data</h2>
          <p>{error || 'Unknown error occurred'}</p>
          <button onClick={() => window.location.reload()}>Retry</button>
        </div>
      </div>
    );
  }

  const { countries, totalGoldValue, totalCountries, totalPopulation, totalGoldReserves } = goldCoinData;

  const handleCountryClick = (country: Country) => {
    setSelectedCountry(country);
    setIsAnimating(true);
    setTimeout(() => setIsAnimating(false), 1000);
  };

  const renderUnifiedView = () => (
    <motion.div
      className={styles.unifiedView}
      initial={{ opacity: 0, scale: 0.8 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.8 }}
    >
      <div className={styles.coinContainer}>
        <motion.div
          className={styles.goldCoin}
          animate={{ rotateY: 360 }}
          transition={{ duration: 2, repeat: Infinity, ease: "linear" }}
        >
          <div className={styles.coinFace}>
            <div className={styles.coinCenter}>
              <div className={styles.africaMap}>
                {countries.map((country, index) => (
                  <motion.div
                    key={country.code}
                    className={styles.countryDot}
                    style={{
                      left: `${Math.random() * 80 + 10}%`,
                      top: `${Math.random() * 80 + 10}%`,
                    }}
                    whileHover={{ scale: 1.5 }}
                    onClick={() => handleCountryClick(country)}
                  >
                    {country.flag}
                  </motion.div>
                ))}
              </div>
              <div className={styles.coinText}>
                <h2>BTNG</h2>
                <p>54 Africa</p>
                <p>Gold Coin</p>
              </div>
            </div>
          </div>
        </motion.div>
      </div>

      <div className={styles.statsContainer}>
        <motion.div
          className={styles.stat}
          initial={{ x: -100, opacity: 0 }}
          animate={{ x: 0, opacity: 1 }}
          transition={{ delay: 0.5 }}
        >
          <h3>54 Nations</h3>
          <p>United in Sovereignty</p>
        </motion.div>

        <motion.div
          className={styles.stat}
          initial={{ y: 100, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.7 }}
        >
          <h3>{totalGoldValue.toFixed(2)} BTNG</h3>
          <p>Total Sovereign Value</p>
        </motion.div>

        <motion.div
          className={styles.stat}
          initial={{ x: 100, opacity: 0 }}
          animate={{ x: 0, opacity: 1 }}
          transition={{ delay: 0.9 }}
        >
          <h3>100% Gold Backed</h3>
          <p>African Sovereign Wealth</p>
        </motion.div>
      </div>
    </motion.div>
  );

  const renderCountriesView = () => (
    <motion.div
      className={styles.countriesView}
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.5 }}
    >
      <h2>All 54 African Nations</h2>
      <div className={styles.countriesGrid}>
        {countries.map((country, index) => (
          <motion.div
            key={country.code}
            className={styles.countryCard}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.02 }}
            whileHover={{ scale: 1.05 }}
            onClick={() => handleCountryClick(country)}
          >
            <div className={styles.countryFlag}>{country.flag}</div>
            <div className={styles.countryInfo}>
              <h4>{country.name}</h4>
              <p>Population: {country.population.toLocaleString()}</p>
              <p>Gold Reserves: {country.goldReserves} tonnes</p>
              <p>Sovereign Value: {country.sovereignValue.toFixed(4)} BTNG</p>
            </div>
          </motion.div>
        ))}
      </div>
    </motion.div>
  );

  const renderSovereignView = () => (
    <motion.div
      className={styles.sovereignView}
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.6 }}
    >
      <h2>BTNG Sovereign Gold Standard</h2>
      <div className={styles.sovereignContent}>
        <div className={styles.goldStandard}>
          <h3>1 BTNG = 1 Gram of Pure African Gold</h3>
          <p>Backed by the collective sovereign wealth of 54 African nations</p>
        </div>

        <div className={styles.valueProposition}>
          <div className={styles.propositionItem}>
            <h4>🇩🇿 Sovereign Backing</h4>
            <p>Each coin represents real African gold reserves and economic potential</p>
          </div>
          <div className={styles.propositionItem}>
            <h4>🌍 Pan-African Unity</h4>
            <p>54 nations united under one sovereign digital gold standard</p>
          </div>
          <div className={styles.propositionItem}>
            <h4>💰 Economic Sovereignty</h4>
            <p>Breaking free from foreign currency dependencies</p>
          </div>
          <div className={styles.propositionItem}>
            <h4>🔒 Blockchain Security</h4>
            <p>Immutable, transparent, and decentralized gold transactions</p>
          </div>
        </div>

        <div className={styles.goldMetrics}>
          <div className={styles.metric}>
            <span className={styles.metricValue}>{totalGoldReserves.toFixed(1)}</span>
            <span className={styles.metricLabel}>Tonnes Total Gold</span>
          </div>
          <div className={styles.metric}>
            <span className={styles.metricValue}>{totalPopulation.toLocaleString()}</span>
            <span className={styles.metricLabel}>People Represented</span>
          </div>
          <div className={styles.metric}>
            <span className={styles.metricValue}>{totalGoldValue.toFixed(2)}</span>
            <span className={styles.metricLabel}>BTNG Sovereign Value</span>
          </div>
        </div>
      </div>
    </motion.div>
  );

  return (
    <div className={`${styles.btng54AfricaGoldCoin} ${className || ''}`}>
      <div className={styles.header}>
        <h1>BTNG 54 Africa Gold Coin</h1>
        <p>The Unified Sovereign Digital Gold Standard for Africa</p>
      </div>

      <div className={styles.viewSelector}>
        <button
          className={currentView === 'unified' ? styles.active : ''}
          onClick={() => setCurrentView('unified')}
        >
          Unified Coin
        </button>
        <button
          className={currentView === 'countries' ? styles.active : ''}
          onClick={() => setCurrentView('countries')}
        >
          54 Nations
        </button>
        <button
          className={currentView === 'sovereign' ? styles.active : ''}
          onClick={() => setCurrentView('sovereign')}
        >
          Sovereign Standard
        </button>
      </div>

      <AnimatePresence mode="wait">
        {currentView === 'unified' && renderUnifiedView()}
        {currentView === 'countries' && renderCountriesView()}
        {currentView === 'sovereign' && renderSovereignView()}
      </AnimatePresence>

      <AnimatePresence>
        {selectedCountry && (
          <motion.div
            className={styles.countryModal}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={() => setSelectedCountry(null)}
          >
            <motion.div
              className={styles.modalContent}
              initial={{ scale: 0.8, y: 20 }}
              animate={{ scale: 1, y: 0 }}
              exit={{ scale: 0.8, y: 20 }}
              onClick={(e) => e.stopPropagation()}
            >
              <div className={styles.modalHeader}>
                <span className={styles.modalFlag}>{selectedCountry.flag}</span>
                <h3>{selectedCountry.name}</h3>
              </div>
              <div className={styles.modalBody}>
                <p><strong>Population:</strong> {selectedCountry.population.toLocaleString()}</p>
                <p><strong>Gold Reserves:</strong> {selectedCountry.goldReserves} tonnes</p>
                <p><strong>Sovereign BTNG Value:</strong> {selectedCountry.sovereignValue.toFixed(6)}</p>
                <p><strong>Country Code:</strong> {selectedCountry.code}</p>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default BTNG54AfricaGoldCoin;