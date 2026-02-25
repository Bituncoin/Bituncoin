import { NextRequest, NextResponse } from 'next/server';
import connectToDatabase from '@/lib/mongodb';

const FABRIC_NODE_CONFIG = {
  nodeId: "nd-6HRNJ6OUIBGP3MV74YAW53NWYQ",
  organization: "btng-root-member",
  network: "btng-fabric-network",
  sovereignId: "BTNG-SOVEREIGN-ROOT-001",
  type: "peer",
  status: "active",
  capabilities: ["endorsement", "validation", "gossip", "events"],
  endpoints: {
    peer: "peer0.btng-root-member.btng-fabric-network.com:7051",
    events: "peer0.btng-root-member.btng-fabric-network.com:7053"
  }
};

export async function GET(request: NextRequest) {
  try {
    // Ensure database connection
    await connectToDatabase();

    const nodeStatus = {
      ...FABRIC_NODE_CONFIG,
      timestamp: new Date().toISOString(),
      uptime: process.uptime(),
      memory: process.memoryUsage(),
      version: "2.5.0",
      ledgerHeight: 150,
      activeChannels: ["btng-sovereign-channel"],
      connectedPeers: 1,
      pendingTransactions: 0,
      health: {
        status: "healthy",
        lastBlockTime: new Date(Date.now() - 30000).toISOString(),
        consensusStatus: "active"
      }
    };

    return NextResponse.json(nodeStatus, {
      status: 200,
      headers: {
        'X-BTNG-Platform': 'Sovereign',
        'X-Fabric-Network': FABRIC_NODE_CONFIG.network,
        'X-Node-ID': FABRIC_NODE_CONFIG.nodeId,
        'X-Node-Status': FABRIC_NODE_CONFIG.status
      }
    });

  } catch (error: any) {
    console.error('BTNG fabric node status error:', error.message);
    return NextResponse.json(
      {
        error: "Internal server error.",
        details: error.message,
        node_status: "error"
      },
      { status: 500 }
    );
  }
}