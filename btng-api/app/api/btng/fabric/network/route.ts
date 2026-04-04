import { NextRequest, NextResponse } from 'next/server';
import { requireAuth } from '@/lib/auth/jwt';
import connectToDatabase from '@/lib/mongodb';

const FABRIC_NETWORK_CONFIG = {
  name: "btng712-fabric-network",
  rootMember: "btng-root-member",
  memberId: "m-SP4QE6LJU5H5ZBO7BGL7RJ4QNQ",
  nodeId: "nd-JKUD2ATMA5A4FAHHTNKWALO2K4",
  sovereignId: "BTNG-SOVEREIGN-ROOT-001",
  status: "Available",
  instanceType: "bc.t3.small",
  availabilityZone: "us-east-1a",
  stateDb: "LevelDB",
  framework: "Hyperledger Fabric",
  frameworkVersion: "2.x",
  region: "us-east-1",
  arn: "arn:aws:managedblockchain:us-east-1:050946999466:nodes/nd-JKUD2ATMA5A4FAHHTNKWALO2K4",
  created: "2026-03-01T20:08:13Z",
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
        peer: "grpcs://nd-jkud2atma5a4fahhtnkwalo2k4.m-sp4qe6lju5h5zbo7bgl7rj4qnq.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com:30003",
        peerEvent: "grpcs://nd-jkud2atma5a4fahhtnkwalo2k4.m-sp4qe6lju5h5zbo7bgl7rj4qnq.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com:30003",
        orderer: "grpcs://orderer.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com:30001",
        ca: "https://ca.m-sp4qe6lju5h5zbo7bgl7rj4qnq.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com:30002",
        ghanaAnchorPeer: "grpcs://154.161.183.158:38982"
      },
      tlsCert: "/opt/home/managedblockchain-tls-chain.pem",
      channels: ["btng712-fabric-network"],
      chaincodes: ["btng-wallet", "btng-gold-token", "btng-sovereign-identity", "btng-liquidity-alr"]
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
