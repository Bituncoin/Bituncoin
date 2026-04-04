import { JWTPayload } from '@/lib/auth/jwt';

export interface FabricWalletResponse {
    success: boolean;
    transactionId?: string;
    wallet_id?: string;
    balance?: number;
    type?: string;
    amount?: number;
    error?: string;
    onChain: boolean;
}

const FABRIC_API_URL = process.env.FABRIC_API_URL || 'http://localhost:3003';
const FABRIC_API_TOKEN = process.env.FABRIC_API_TOKEN || 'btng-internal-token';

async function fabricInvoke(
    chaincodeName: string,
    channelName: string,
    fn: string,
    args: string[]
): Promise<any> {
    const response = await fetch(`${FABRIC_API_URL}/api/btng/fabric/chaincode`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${FABRIC_API_TOKEN}`,
        },
        body: JSON.stringify({
            chaincodeName,
            channelName,
            operation: { function: fn, args },
        }),
    });

    if (!response.ok) {
        throw new Error(`Fabric API error: ${response.status} ${response.statusText}`);
    }

    return response.json();
}

/**
 * Mint BTNG tokens on-chain via Fabric btng-wallet chaincode
 */
export async function fabricMint(
    walletId: string,
    amount: number,
    goldPriceUsd: number = 0
): Promise<FabricWalletResponse> {
    try {
        const result = await fabricInvoke(
            'btng-wallet',
            'btng712-fabric-network',
            'Mint',
            [walletId, amount.toString(), goldPriceUsd.toString()]
        );

        return {
            success: true,
            transactionId: result.transactionId || result.tx_id,
            wallet_id: walletId,
            balance: result.balance_after,
            type: 'mint',
            amount,
            onChain: true,
        };
    } catch (error: any) {
        console.error('Fabric mint error:', error.message);
        return { success: false, error: error.message, onChain: false };
    }
}

/**
 * Melt BTNG tokens on-chain via Fabric btng-wallet chaincode
 */
export async function fabricMelt(
    walletId: string,
    amount: number,
    goldPriceUsd: number = 0
): Promise<FabricWalletResponse> {
    try {
        const result = await fabricInvoke(
            'btng-wallet',
            'btng712-fabric-network',
            'Melt',
            [walletId, amount.toString(), goldPriceUsd.toString()]
        );

        return {
            success: true,
            transactionId: result.transactionId || result.tx_id,
            wallet_id: walletId,
            balance: result.balance_after,
            type: 'melt',
            amount,
            onChain: true,
        };
    } catch (error: any) {
        console.error('Fabric melt error:', error.message);
        return { success: false, error: error.message, onChain: false };
    }
}

/**
 * Get wallet balance from on-chain state
 */
export async function fabricGetBalance(walletId: string): Promise<FabricWalletResponse> {
    try {
        const result = await fabricInvoke(
            'btng-wallet',
            'btng712-fabric-network',
            'GetBalance',
            [walletId]
        );

        return {
            success: true,
            wallet_id: walletId,
            balance: result.balance || 0,
            onChain: true,
        };
    } catch (error: any) {
        console.error('Fabric balance query error:', error.message);
        return { success: false, error: error.message, onChain: false };
    }
}

/**
 * Transfer BTNG tokens on-chain between wallets
 */
export async function fabricTransfer(
    fromWallet: string,
    toWallet: string,
    amount: number
): Promise<FabricWalletResponse> {
    try {
        const result = await fabricInvoke(
            'btng-wallet',
            'btng712-fabric-network',
            'Transfer',
            [fromWallet, toWallet, amount.toString()]
        );

        return {
            success: true,
            transactionId: result.transactionId || result.tx_id,
            wallet_id: fromWallet,
            balance: result.balance_after,
            type: 'transfer',
            amount,
            onChain: true,
        };
    } catch (error: any) {
        console.error('Fabric transfer error:', error.message);
        return { success: false, error: error.message, onChain: false };
    }
}
