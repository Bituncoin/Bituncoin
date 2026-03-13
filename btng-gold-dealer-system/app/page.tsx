'use client'

import { btng } from '../../lib/btng';

export default async function DealerDashboard() {
  const price = await btng.goldPriceLatest();
  const { goldPrice } = price;

  const buyPrice = goldPrice * 1.02; // 2% spread
  const sellPrice = goldPrice * 0.98; // 2% spread

  const vat = buyPrice * 0.15; // 15% VAT

  return (
    <div className="p-8 max-w-6xl mx-auto">
      <h1 className="text-4xl font-bold mb-8">BTNG Gold Dealer Terminal</h1>
      
      <div className="grid md:grid-cols-3 gap-6 mb-8">
        <div className="bg-gradient-to-r from-yellow-400 to-amber-500 p-6 rounded-xl text-white">
          <h2 className="text-2xl font-semibold">Spot Price</h2>
          <p className="text-3xl font-bold">${goldPrice.toFixed(2)}/g</p>
        </div>
        
        <div className="bg-gradient-to-r from-green-400 to-emerald-500 p-6 rounded-xl text-white">
          <h2 className="text-2xl font-semibold">Buy Price</h2>
          <p className="text-3xl font-bold">${buyPrice.toFixed(2)}/g</p>
          <p>+ VAT: GHS {vat.toFixed(2)}</p>
        </div>
        
        <div className="bg-gradient-to-r from-red-400 to-rose-500 p-6 rounded-xl text-white">
          <h2 className="text-2xl font-semibold">Sell Price</h2>
          <p className="text-3xl font-bold">${sellPrice.toFixed(2)}/g</p>
        </div>
      </div>
      
      <div className="grid md:grid-cols-2 gap-6">
        <div className="p-6 bg-white rounded-xl shadow-lg border">
          <h3 className="text-xl font-semibold mb-4">Quick Buy</h3>
          <input type="number" placeholder="Grams" className="w-full p-3 border rounded-lg mb-4" />
          <button className="w-full bg-green-500 text-white py-3 rounded-lg font-semibold">
            Generate Buy Quote
          </button>
        </div>
        
        <div className="p-6 bg-white rounded-xl shadow-lg border">
          <h3 className="text-xl font-semibold mb-4">Quick Sell</h3>
          <input type="number" placeholder="Grams" className="w-full p-3 border rounded-lg mb-4" />
          <button className="w-full bg-red-500 text-white py-3 rounded-lg font-semibold">
            Generate Sell Quote
          </button>
        </div>
      </div>
      
      <div className="mt-8 p-6 bg-blue-50 rounded-xl border">
        <h3 className="text-lg font-semibold mb-2">Recent Transactions</h3>
        <div className="text-sm text-gray-600">
          No transactions today. Ready for first customer.
        </div>
      </div>
    </div>
  );
}
