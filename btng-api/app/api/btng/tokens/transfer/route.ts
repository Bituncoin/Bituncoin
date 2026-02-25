import { NextRequest, NextResponse } from 'next/server';
import { requireAuth } from '@/lib/auth/jwt';
import { transferGoldTokens } from '@/lib/fabric/chaincode';

import { JWTPayload } from '@/lib/auth/jwt';

async function handleTransferTokens(request: NextRequest, user: JWTPayload) {
  try {
    const body = await request.json();
    const { to_user_id, amount } = body;

    // Validate required fields
    if (!to_user_id || !amount) {
      return NextResponse.json(
        { error: "Missing required fields: to_user_id and amount" },
        { status: 400 }
      );
    }

    // Validate amount is positive number
    const tokenAmount = parseInt(amount);
    if (isNaN(tokenAmount) || tokenAmount <= 0) {
      return NextResponse.json(
        { error: "Amount must be a positive integer" },
        { status: 400 }
      );
    }

    // Prevent self-transfer
    if (to_user_id === user.sub) {
      return NextResponse.json(
        { error: "Cannot transfer tokens to yourself" },
        { status: 400 }
      );
    }

    const transferResult = await transferGoldTokens(user.sub, to_user_id, tokenAmount, user);

    if (!transferResult.success) {
      return NextResponse.json(
        { error: "Token transfer failed", details: transferResult.error },
        { status: 500 }
      );
    }

    return NextResponse.json({
      status: "ok",
      from_user: user.sub,
      to_user: to_user_id,
      amount: tokenAmount,
      transaction_id: transferResult.transactionId,
      message: "Gold tokens transferred successfully"
    }, {
      status: 200,
      headers: {
        'X-BTNG-Platform': 'Sovereign',
        'X-Service-Status': 'Operational',
        'X-Auth-User': user.sub,
        'X-Transaction-ID': transferResult.transactionId || ''
      }
    });

  } catch (error: any) {
    console.error("BTNG token transfer endpoint error:", error.message);
    return NextResponse.json(
      { error: "Internal server error.", details: error.message },
      { status: 500 }
    );
  }
}

export const POST = requireAuth(handleTransferTokens);