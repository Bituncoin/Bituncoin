'use client';

import { useEffect } from 'react';

export default function BTNGDemoClient() {
  useEffect(() => {
    // Wait for the external script to load
    const setupEventListeners = () => {
      // Navigation buttons
      const navButtons = document.querySelectorAll('.nav-btn[data-section]');
      navButtons.forEach((btn) => {
        btn.addEventListener('click', function(this: HTMLElement) {
          const section = this.getAttribute('data-section');
          if (section && typeof (window as any).showSection === 'function') {
            (window as any).showSection(section, this);
          }
        });
      });

      // Copy address button
      const copyBtn = document.getElementById('copy-address-btn');
      if (copyBtn) {
        copyBtn.addEventListener('click', () => {
          if (typeof (window as any).copyAddress === 'function') {
            (window as any).copyAddress();
          }
        });
      }

      // Toggle QR button
      const toggleQRBtn = document.getElementById('toggle-qr-btn');
      if (toggleQRBtn) {
        toggleQRBtn.addEventListener('click', () => {
          if (typeof (window as any).toggleWalletQR === 'function') {
            (window as any).toggleWalletQR();
          }
        });
      }

      // Send transaction button
      const sendTxBtn = document.getElementById('send-transaction-btn');
      if (sendTxBtn) {
        sendTxBtn.addEventListener('click', () => {
          if (typeof (window as any).sendTransaction === 'function') {
            (window as any).sendTransaction();
          }
        });
      }

      // Mining toggle button
      const mineBtn = document.getElementById('mine-btn');
      if (mineBtn) {
        mineBtn.addEventListener('click', () => {
          if (typeof (window as any).toggleMining === 'function') {
            (window as any).toggleMining();
          }
        });
      }

      // Search button
      const searchBtn = document.getElementById('search-btn');
      if (searchBtn) {
        searchBtn.addEventListener('click', () => {
          if (typeof (window as any).doSearch === 'function') {
            (window as any).doSearch();
          }
        });
      }

      // Search input enter key
      const searchInput = document.getElementById('exp-search');
      if (searchInput) {
        searchInput.addEventListener('keydown', (e: KeyboardEvent) => {
          if (e.key === 'Enter' && typeof (window as any).doSearch === 'function') {
            (window as any).doSearch();
          }
        });
      }
    };

    // Try to setup immediately
    setupEventListeners();

    // Also setup after a delay to ensure external script is loaded
    const timeoutId = setTimeout(setupEventListeners, 1000);

    return () => {
      clearTimeout(timeoutId);
    };
  }, []);

  return null;
}
