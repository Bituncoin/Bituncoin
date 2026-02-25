import { NextRequest, NextResponse } from 'next/server';
import connectToDatabase from '@/lib/mongodb';
import { getGoldPrice } from '@/lib/gold-price/model';
import { getAllCountries } from '@/lib/countries';

export type ImperialReserveData = {
    totalNations: number;
    activeNations: number;
    totalGoldReserve: number;
    goldBackingRatio: number;
    merchantCount: number;
    tradeVolume24h: number;
    nations: NationReserveData[];
    lastUpdated: string;
};

export type NationReserveData = {
    name: string;
    code: string;
    goldReserve: number;
    merchantCount: number;
    tradeVolume: number;
    uptime: number;
    status: 'online' | 'offline' | 'deploying';
    lastUpdated: string;
};

async function getImperialReserveData(): Promise<ImperialReserveData> {
    // Get all nations data
    const allCountries = getAllCountries();
    const nations: NationReserveData[] = [];

    let totalGoldReserve = 0;
    let totalMerchants = 0;
    let totalTradeVolume = 0;
    let activeNations = 0;

    for (const [code, country] of allCountries) {
        // Simulate nation data (in real implementation, this would come from database)
        const isActive = Math.random() > 0.1; // 90% of nations active
        const goldReserve = isActive ? Math.random() * 10000 + 5000 : 0;
        const merchantCount = isActive ? Math.floor(Math.random() * 100) + 10 : 0;
        const tradeVolume = isActive ? Math.random() * 100000 + 10000 : 0;
        const uptime = isActive ? 99.9 : 0;

        nations.push({
            name: country.name,
            code: code,
            goldReserve,
            merchantCount,
            tradeVolume,
            uptime,
            status: isActive ? 'online' : 'offline',
            lastUpdated: new Date().toISOString()
        });

        if (isActive) {
            activeNations++;
            totalGoldReserve += goldReserve;
            totalMerchants += merchantCount;
            totalTradeVolume += tradeVolume;
        }
    }

    // Calculate gold backing ratio (gold reserve vs trade volume)
    const goldBackingRatio = totalTradeVolume > 0 ? totalGoldReserve / totalTradeVolume : 0;

    return {
        totalNations: nations.length,
        activeNations,
        totalGoldReserve,
        goldBackingRatio,
        merchantCount: totalMerchants,
        tradeVolume24h: totalTradeVolume,
        nations,
        lastUpdated: new Date().toISOString()
    };
}

export async function GET(request: NextRequest) {
    try {
        await connectToDatabase();

        const reserveData = await getImperialReserveData();
        const goldPrice = await getGoldPrice();

        return NextResponse.json({
            ...reserveData,
            currentGoldPrice: goldPrice ? {
                usdPerOunce: goldPrice.price_ounce,
                usdPerGram: goldPrice.price_gram,
                lastUpdated: goldPrice.timestamp
            } : null,
            message: "BTNG Imperial Reserve: 54 Nations Backed by Gold",
            status: "ACTIVE"
        }, {
            status: 200,
            headers: {
                'X-BTNG-Dashboard': 'Imperial Reserve',
                'X-Gold-Backing': reserveData.goldBackingRatio.toFixed(4),
                'X-Active-Nations': reserveData.activeNations.toString(),
                'Cache-Control': 'public, max-age=300' // Cache for 5 minutes
            }
        });

    } catch (error: any) {
        console.error('BTNG Imperial Reserve Dashboard error:', error.message);
        return NextResponse.json(
            { error: 'Dashboard temporarily unavailable', details: error.message },
            { status: 503 }
        );
    }
}