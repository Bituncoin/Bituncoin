import { NextRequest, NextResponse } from 'next/server';
import { requireAuth } from '@/lib/auth/jwt';
import connectToDatabase from '@/lib/mongodb';

const FABRIC_NETWORK_CONFIG = {
  name: "btng-fabric-network",
  rootMember: "btng-root-member",
  nodeId: "nd-6HRNJ6OUIBGP3MV74YAW53NWYQ",
  sovereignId: "BTNG-SOVEREIGN-ROOT-001",
  status: "active",
  capabilities: ["endorsement", "validation", "gossip", "events"]
};

export async function GET(request: NextRequest) {
  try {
    // Ensure database connection
    await connectToDatabase();

    const networkInfo = {
      network: FABRIC_NETWORK_CONFIG,
      timestamp: new Date().toISOString(),
      status: "operational",
      endpoints: {
        peer: "peer0.btng-root-member.btng-fabric-network.com:7051",
        orderer: "orderer.btng-fabric-network.com:7050",
        ca: "ca.btng-root-member.btng-fabric-network.com:7054"
      },
      channels: ["btng-sovereign-channel"],
      chaincodes: ["btng-gold-token", "btng-sovereign-identity"]
    };

    return NextResponse.json(networkInfo, {
      status: 200,
      headers: {
        'X-BTNG-Platform': 'Sovereign',
        'X-Fabric-Network': FABRIC_NETWORK_CONFIG.name,
        'X-Node-ID': FABRIC_NETWORK_CONFIG.nodeId
      }
    });

  } catch (error: any) {
    console.error('BTNG fabric network status error:', error.message);
    return NextResponse.json(
      {
        error: "Internal server error.",
        details: error.message,
        network_status: "error"
      },
      { status: 500 }
    );
  }
}

async function handleFabricNetworkUpdate(request: NextRequest, user: any) {
  try {
    // Ensure database connection
    await connectToDatabase();

    const body = await request.json();
    const { action, channel, chaincode, transaction } = body;

    // Validate required fields
    if (!action) {
      return NextResponse.json(
        { error: "Action required (deploy, invoke, query)" },
        { status: 400 }
      );
    }

    // Simulate fabric network operations
    const fabricResponse = {
      action: action,
      channel: channel || "btng-sovereign-channel",
      chaincode: chaincode || "btng-gold-token",
      transaction: transaction,
      nodeId: FABRIC_NETWORK_CONFIG.nodeId,
      sovereignId: FABRIC_NETWORK_CONFIG.sovereignId,
      timestamp: new Date().toISOString(),
      status: "success",
      result: {
        txId: `btng-tx-${Date.now()}`,
        blockNumber: Math.floor(Math.random() * 1000),
        endorsementCount: 1,
        validationStatus: "committed"
      }
    };

    console.log('BTNG Fabric Network Operation by:', user.sub, {
      action: action,
      channel: channel,
      chaincode: chaincode,
      nodeId: FABRIC_NETWORK_CONFIG.nodeId
    });

    return NextResponse.json({
      status: "ok",
      operation: fabricResponse,
      message: `Fabric network ${action} operation completed successfully`
    }, {
      status: 200,
      headers: {
        'X-BTNG-Platform': 'Sovereign',
        'X-Fabric-Network': FABRIC_NETWORK_CONFIG.name,
        'X-Node-ID': FABRIC_NETWORK_CONFIG.nodeId,
        'X-Operation': action
      }
    });

  } catch (error: any) {
    console.error("BTNG fabric network operation error:", error.message);
    return NextResponse.json(
      { error: "Internal server error.", details: error.message },
      { status: 500 }
    );
  }
}

export const POST = requireAuth(handleFabricNetworkUpdate);