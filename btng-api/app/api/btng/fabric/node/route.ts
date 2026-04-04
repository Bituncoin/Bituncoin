import { NextRequest, NextResponse } from 'next/server';
import connectToDatabase from '@/lib/mongodb';

const FABRIC_NODE_CONFIG = {
  nodeId: "nd-JKUD2ATMA5A4FAHHTNKWALO2K4",
  memberId: "m-SP4QE6LJU5H5ZBO7BGL7RJ4QNQ",
  organization: "btng-root-member",
  network: "btng712-fabric-network",
  sovereignId: "BTNG-SOVEREIGN-ROOT-001",
  type: "peer",
  instanceType: "bc.t3.small",
  availabilityZone: "us-east-1a",
  stateDb: "LevelDB",
  status: "Available",
  arn: "arn:aws:managedblockchain:us-east-1:050946999466:nodes/nd-JKUD2ATMA5A4FAHHTNKWALO2K4",
  created: "2026-03-01T20:08:13Z",
  capabilities: ["endorsement", "validation", "gossip", "events"],
  endpoints: {
    peer: "grpcs://nd-jkud2atma5a4fahhtnkwalo2k4.m-sp4qe6lju5h5zbo7bgl7rj4qnq.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com:30003",
    events: "grpcs://nd-jkud2atma5a4fahhtnkwalo2k4.m-sp4qe6lju5h5zbo7bgl7rj4qnq.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com:30003",
    ca: "https://ca.m-sp4qe6lju5h5zbo7bgl7rj4qnq.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com:30002",
    ghanaAnchorPeer: "grpcs://154.161.183.158:38982"
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
      activeChannels: ["btng712-fabric-network"],
      connectedPeers: 2,
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
