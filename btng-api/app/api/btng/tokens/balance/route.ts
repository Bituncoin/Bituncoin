import { NextRequest, NextResponse } from 'next/server';
import { requireAuth } from '@/lib/auth/jwt';
import { getGoldTokenBalance } from '@/lib/fabric/chaincode';

import { JWTPayload } from '@/lib/auth/jwt';

async function handleGetBalance(request: NextRequest, user: JWTPayload) {
  try {
    const balanceResult = await getGoldTokenBalance(user.sub);

    if (!balanceResult.success) {
      return NextResponse.json(
        { error: "Failed to retrieve token balance", details: balanceResult.error },
        { status: 500 }
      );
    }

    return NextResponse.json({
      status: "ok",
      user_id: user.sub,
      balance: balanceResult.result?.balance || 0,
      last_updated: balanceResult.result?.lastUpdated || null,
      message: "Gold token balance retrieved successfully"
    }, {
      status: 200,
      headers: {
        'X-BTNG-Platform': 'Sovereign',
        'X-Service-Status': 'Operational',
        'X-Auth-User': user.sub
      }
    });

  } catch (error: any) {
    console.error("BTNG token balance endpoint error:", error.message);
    return NextResponse.json(
      { error: "Internal server error.", details: error.message },
      { status: 500 }
    );
  }
}

export const GET = requireAuth(handleGetBalance);