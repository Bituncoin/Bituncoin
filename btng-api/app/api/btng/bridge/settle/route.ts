import { NextRequest, NextResponse } from 'next/server';
import { requireAuth } from '@/lib/auth/jwt';
import connectToDatabase from '@/lib/mongodb';
import { getGoldPrice } from '@/lib/gold-price/model';
import { transferGoldTokens } from '@/lib/fabric/chaincode';

export type TradeBridgeSettlement = {
    id: string;
    fromMerchant: string;
    toMerchant: string;
    fromNation: string;
    toNation: string;
    amount: number;
    currency: string;
    goldEquivalent: number;
    status: 'pending' | 'settled' | 'failed';
    timestamp: string;
    transactionId?: string;
};

async function handleBridgeSettlement(request: NextRequest, user: any) {
    try {
        await connectToDatabase();

        const body = await request.json();
        const {
            fromMerchant,
            toMerchant,
            fromNation,
            toNation,
            amount,
            currency,
            invoiceId
        } = body;

        // Validate required fields
        if (!fromMerchant || !toMerchant || !amount || !currency) {
            return NextResponse.json(
                { error: 'Missing required settlement fields' },
                { status: 400 }
            );
        }

        // Get current gold price for conversion
        const goldPrice = await getGoldPrice();
        if (!goldPrice) {
            return NextResponse.json(
                { error: 'Unable to retrieve gold price for settlement' },
                { status: 500 }
            );
        }

        // Calculate gold equivalent
        let goldEquivalent: number;
        if (currency === 'BTNG') {
            goldEquivalent = amount; // Already in gold tokens
        } else {
            // Convert fiat to gold equivalent
            const pricePerGram = goldPrice.price_gram;
            goldEquivalent = amount / pricePerGram;
        }

        // Create settlement record
        const settlement: TradeBridgeSettlement = {
            id: `bridge-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
            fromMerchant,
            toMerchant,
            fromNation: fromNation || 'Unknown',
            toNation: toNation || 'Unknown',
            amount,
            currency,
            goldEquivalent,
            status: 'pending',
            timestamp: new Date().toISOString()
        };

        // Execute gold token transfer
        const transferResult = await transferGoldTokens(
            fromMerchant,
            toMerchant,
            goldEquivalent
        );

        if (transferResult.success) {
            settlement.status = 'settled';
            settlement.transactionId = transferResult.transactionId;

            // Store settlement record
            await storeSettlementRecord(settlement);

            return NextResponse.json({
                status: 'ok',
                settlement,
                message: `Pan-African settlement completed: ${amount} ${currency} transferred in gold`,
                goldTransferred: goldEquivalent,
                transactionId: transferResult.transactionId
            }, {
                status: 200,
                headers: {
                    'X-BTNG-Bridge': 'Pan-African',
                    'X-Settlement-ID': settlement.id,
                    'X-Gold-Equivalent': goldEquivalent.toString()
                }
            });
        } else {
            settlement.status = 'failed';
            await storeSettlementRecord(settlement);

            return NextResponse.json({
                error: 'Settlement failed',
                settlement,
                details: transferResult.error
            }, { status: 500 });
        }

    } catch (error: any) {
        console.error('BTNG Pan-African Trade Bridge error:', error.message);
        return NextResponse.json(
            { error: 'Internal server error', details: error.message },
            { status: 500 }
        );
    }
}

async function storeSettlementRecord(settlement: TradeBridgeSettlement): Promise<void> {
    // Implementation would store in MongoDB
    console.log('Settlement record stored:', settlement);
}

export const POST = requireAuth(handleBridgeSettlement);