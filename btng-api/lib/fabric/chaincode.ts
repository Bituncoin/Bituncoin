import { JWTPayload } from '@/lib/auth/jwt';

export interface ChaincodeOperation {
  function: string;
  args: string[];
}

export interface ChaincodeResponse {
  success: boolean;
  transactionId?: string;
  result?: any;
  error?: string;
}

/**
 * Mint gold tokens on Fabric network when gold price is updated
 */
export async function mintGoldTokensOnPriceUpdate(
  goldPriceRecord: any,
  user: JWTPayload
): Promise<ChaincodeResponse> {
  try {
    // Calculate token amount based on gold price (1 gram = 1 token for simplicity)
    const tokenAmount = Math.floor(goldPriceRecord.base_price_gram * 100); // Convert to cents for precision

    const operation: ChaincodeOperation = {
      function: 'Mint',
      args: [
        user.sub, // recipient (admin user)
        tokenAmount.toString(),
        `BTNG-${Date.now()}`, // unique token ID
        goldPriceRecord.timestamp.toString(),
        goldPriceRecord.base_price_gram.toString()
      ]
    };

    // Call Fabric chaincode API
    const response = await fetch(`${process.env.FABRIC_API_URL || 'http://localhost:3003'}/api/btng/fabric/chaincode`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${process.env.FABRIC_API_TOKEN || 'btng-internal-token'}`
      },
      body: JSON.stringify({
        chaincodeName: 'btng-gold-token',
        channelName: 'btng-channel',
        operation
      })
    });

    if (!response.ok) {
      throw new Error(`Fabric API error: ${response.status} ${response.statusText}`);
    }

    const result = await response.json();

    console.log('BTNG Gold Token Minted:', {
      user: user.sub,
      amount: tokenAmount,
      transactionId: result.transactionId,
      timestamp: goldPriceRecord.timestamp
    });

    return {
      success: true,
      transactionId: result.transactionId,
      result: result
    };

  } catch (error: any) {
    console.error('BTNG Fabric chaincode mint error:', error.message);
    return {
      success: false,
      error: error.message
    };
  }
}

/**
 * Get gold token balance for a user
 */
export async function getGoldTokenBalance(userId: string): Promise<ChaincodeResponse> {
  try {
    const operation: ChaincodeOperation = {
      function: 'GetBalance',
      args: [userId]
    };

    const response = await fetch(`${process.env.FABRIC_API_URL || 'http://localhost:3003'}/api/btng/fabric/chaincode`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${process.env.FABRIC_API_TOKEN || 'btng-internal-token'}`
      },
      body: JSON.stringify({
        chaincodeName: 'btng-gold-token',
        channelName: 'btng-channel',
        operation
      })
    });

    if (!response.ok) {
      throw new Error(`Fabric API error: ${response.status} ${response.statusText}`);
    }

    const result = await response.json();

    return {
      success: true,
      result: result
    };

  } catch (error: any) {
    console.error('BTNG Fabric balance query error:', error.message);
    return {
      success: false,
      error: error.message
    };
  }
}

/**
 * Transfer gold tokens between users
 */
export async function transferGoldTokens(
  fromUserId: string,
  toUserId: string,
  amount: number,
  user: JWTPayload
): Promise<ChaincodeResponse> {
  try {
    const operation: ChaincodeOperation = {
      function: 'Transfer',
      args: [
        fromUserId,
        toUserId,
        amount.toString(),
        `TX-${Date.now()}`
      ]
    };

    const response = await fetch(`${process.env.FABRIC_API_URL || 'http://localhost:3003'}/api/btng/fabric/chaincode`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${process.env.FABRIC_API_TOKEN || 'btng-internal-token'}`
      },
      body: JSON.stringify({
        chaincodeName: 'btng-gold-token',
        channelName: 'btng-channel',
        operation
      })
    });

    if (!response.ok) {
      throw new Error(`Fabric API error: ${response.status} ${response.statusText}`);
    }

    const result = await response.json();

    console.log('BTNG Gold Token Transfer:', {
      from: fromUserId,
      to: toUserId,
      amount: amount,
      transactionId: result.transactionId,
      initiatedBy: user.sub
    });

    return {
      success: true,
      transactionId: result.transactionId,
      result: result
    };

  } catch (error: any) {
    console.error('BTNG Fabric transfer error:', error.message);
    return {
      success: false,
      error: error.message
    };
  }
}