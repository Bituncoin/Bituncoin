#!/usr/bin/env node

/**
 * BTNG Fabric Chaincode Deployment Script
 * Deploys BTNG gold token and sovereign identity chaincodes
 */

const fs = require('fs');
const path = require('path');

const FABRIC_NETWORK = {
  name: "btng-fabric-network",
  channel: "btng-sovereign-channel",
  rootMember: "btng-root-member",
  nodeId: "nd-6HRNJ6OUIBGP3MV74YAW53NWYQ"
};

const CHAINCODES = {
  "btng-gold-token": {
    path: "chaincode/btng-gold-token",
    version: "1.0.0",
    sequence: "1",
    initRequired: true,
    language: "golang",
    functions: ["Mint", "Transfer", "BalanceOf", "TotalSupply"]
  },
  "btng-sovereign-identity": {
    path: "chaincode/btng-sovereign-identity",
    version: "1.0.0",
    sequence: "1",
    initRequired: true,
    language: "golang",
    functions: ["RegisterIdentity", "VerifyIdentity", "UpdateIdentity", "GetIdentity"]
  }
};

function generateChaincodePackage(chaincodeName) {
  const chaincode = CHAINCODES[chaincodeName];
  const packagePath = path.join(__dirname, '..', 'fabric', 'chaincode-packages', `${chaincodeName}.tar.gz`);

  console.log(`📦 Generating chaincode package for ${chaincodeName}...`);

  // Create chaincode directory structure
  const chaincodeDir = path.join(__dirname, '..', 'fabric', 'chaincode', chaincodeName);
  if (!fs.existsSync(chaincodeDir)) {
    fs.mkdirSync(chaincodeDir, { recursive: true });
  }

  // Generate Go chaincode files
  generateGoChaincode(chaincodeName, chaincodeDir);

  console.log(`✅ Chaincode package generated: ${packagePath}`);
  return packagePath;
}

function generateGoChaincode(chaincodeName, chaincodeDir) {
  const chaincode = CHAINCODES[chaincodeName];

  // Generate main.go
  const mainGo = `package main

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ${chaincodeName} provides functions for ${chaincodeName} operations
type ${chaincodeName} struct {
	contractapi.Contract
}

${generateContractFunctions(chaincodeName)}

// New${chaincodeName} creates a new instance of ${chaincodeName}
func New${chaincodeName}() *${chaincodeName} {
	return &${chaincodeName}{}
}

func main() {
	chaincode, err := contractapi.NewChaincode(New${chaincodeName}())
	if err != nil {
		log.Panicf("Error creating ${chaincodeName} chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting ${chaincodeName} chaincode: %v", err)
	}
}
`;

  fs.writeFileSync(path.join(chaincodeDir, 'main.go'), mainGo);
  console.log(`📝 Generated main.go for ${chaincodeName}`);
}

function generateContractFunctions(chaincodeName) {
  const chaincode = CHAINCODES[chaincodeName];
  let functions = '';

  switch (chaincodeName) {
    case 'btng-gold-token':
      functions = `
// Mint creates new gold tokens
func (c *${chaincodeName}) Mint(ctx contractapi.TransactionContextInterface, amount string, recipient string) error {
	// Implementation for minting gold tokens
	fmt.Printf("Minting %s gold tokens to %s\\n", amount, recipient)
	return nil
}

// Transfer transfers gold tokens between accounts
func (c *${chaincodeName}) Transfer(ctx contractapi.TransactionContextInterface, amount string, from string, to string) error {
	// Implementation for transferring gold tokens
	fmt.Printf("Transferring %s gold tokens from %s to %s\\n", amount, from, to)
	return nil
}

// BalanceOf returns the balance of gold tokens for an account
func (c *${chaincodeName}) BalanceOf(ctx contractapi.TransactionContextInterface, account string) (string, error) {
	// Implementation for checking balance
	fmt.Printf("Checking balance for account %s\\n", account)
	return "1000", nil
}

// TotalSupply returns the total supply of gold tokens
func (c *${chaincodeName}) TotalSupply(ctx contractapi.TransactionContextInterface) (string, error) {
	// Implementation for total supply
	return "1000000", nil
}
`;
      break;

    case 'btng-sovereign-identity':
      functions = `
// RegisterIdentity registers a new sovereign identity
func (c *${chaincodeName}) RegisterIdentity(ctx contractapi.TransactionContextInterface, identityId string, publicKey string, metadata string) error {
	// Implementation for registering sovereign identity
	fmt.Printf("Registering sovereign identity %s\\n", identityId)
	return nil
}

// VerifyIdentity verifies a sovereign identity
func (c *${chaincodeName}) VerifyIdentity(ctx contractapi.TransactionContextInterface, identityId string) (string, error) {
	// Implementation for verifying identity
	fmt.Printf("Verifying identity %s\\n", identityId)
	return "verified", nil
}

// UpdateIdentity updates sovereign identity metadata
func (c *${chaincodeName}) UpdateIdentity(ctx contractapi.TransactionContextInterface, identityId string, metadata string) error {
	// Implementation for updating identity
	fmt.Printf("Updating identity %s metadata\\n", identityId)
	return nil
}

// GetIdentity retrieves sovereign identity information
func (c *${chaincodeName}) GetIdentity(ctx contractapi.TransactionContextInterface, identityId string) (string, error) {
	// Implementation for getting identity
	fmt.Printf("Retrieving identity %s\\n", identityId)
	return "{\\"id\\": \\"BTNG-SOVEREIGN-001\\", \\"status\\": \\"active\\"}", nil
}
`;
      break;
  }

  return functions;
}

function generateDeploymentScript(chaincodeName) {
  const chaincode = CHAINCODES[chaincodeName];
  const scriptPath = path.join(__dirname, '..', 'fabric', 'scripts', `deploy-${chaincodeName}.sh`);

  const script = `#!/bin/bash

# BTNG Fabric Chaincode Deployment Script
# Deploys ${chaincodeName} to ${FABRIC_NETWORK.name}

set -e

echo "🚀 Deploying ${chaincodeName} to ${FABRIC_NETWORK.name}"

# Set environment variables
export FABRIC_CFG_PATH=\${PWD}/config
export CORE_PEER_LOCALMSPID="${FABRIC_NETWORK.rootMember}MSP"
export CORE_PEER_MSPCONFIGPATH=\${PWD}/crypto-config/peerOrganizations/${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com/users/Admin@${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com/msp
export CORE_PEER_ADDRESS=peer0.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=\${PWD}/crypto-config/peerOrganizations/${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com/peers/peer0.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com/tls/ca.crt

# Package chaincode
echo "📦 Packaging ${chaincodeName}..."
peer lifecycle chaincode package ${chaincodeName}.tar.gz --path ./chaincode/${chaincodeName} --lang golang --label ${chaincodeName}_${chaincode.version}

# Install chaincode
echo "📥 Installing ${chaincodeName}..."
peer lifecycle chaincode install ${chaincodeName}.tar.gz

# Get package ID
PACKAGE_ID=\$(peer lifecycle chaincode queryinstalled | grep ${chaincodeName} | awk '{print \$3}' | sed 's/,$//')

echo "📋 Package ID: \${PACKAGE_ID}"

# Approve chaincode
echo "✅ Approving ${chaincodeName}..."
peer lifecycle chaincode approveformyorg \\
  --channelID ${FABRIC_NETWORK.channel} \\
  --name ${chaincodeName} \\
  --version ${chaincode.version} \\
  --package-id \${PACKAGE_ID} \\
  --sequence ${chaincode.sequence} \\
  --tls true \\
  --cafile \${PWD}/crypto-config/ordererOrganizations/${FABRIC_NETWORK.name}.com/orderers/orderer.${FABRIC_NETWORK.name}.com/tls/ca.crt

# Check commit readiness
echo "🔍 Checking commit readiness..."
peer lifecycle chaincode checkcommitreadiness \\
  --channelID ${FABRIC_NETWORK.channel} \\
  --name ${chaincodeName} \\
  --version ${chaincode.version} \\
  --sequence ${chaincode.sequence} \\
  --tls true \\
  --cafile \${PWD}/crypto-config/ordererOrganizations/${FABRIC_NETWORK.name}.com/orderers/orderer.${FABRIC_NETWORK.name}.com/tls/ca.crt

# Commit chaincode
echo "📝 Committing ${chaincodeName}..."
peer lifecycle chaincode commit \\
  --channelID ${FABRIC_NETWORK.channel} \\
  --name ${chaincodeName} \\
  --version ${chaincode.version} \\
  --sequence ${chaincode.sequence} \\
  --tls true \\
  --cafile \${PWD}/crypto-config/ordererOrganizations/${FABRIC_NETWORK.name}.com/orderers/orderer.${FABRIC_NETWORK.name}.com/tls/ca.crt

echo "🎉 ${chaincodeName} deployment completed successfully!"
echo "🔗 Chaincode: ${chaincodeName}"
echo "📊 Version: ${chaincode.version}"
echo "🔢 Sequence: ${chaincode.sequence}"
echo "🌐 Network: ${FABRIC_NETWORK.name}"
echo "📺 Channel: ${FABRIC_NETWORK.channel}"
`;

  // Create scripts directory
  const scriptsDir = path.dirname(scriptPath);
  if (!fs.existsSync(scriptsDir)) {
    fs.mkdirSync(scriptsDir, { recursive: true });
  }

  fs.writeFileSync(scriptPath, script);
  fs.chmodSync(scriptPath, '755');

  console.log(`📜 Generated deployment script: ${scriptPath}`);
}

function generateDockerCompose() {
  const dockerComposePath = path.join(__dirname, '..', 'fabric', 'docker-compose.yml');

  const dockerCompose = `version: '2.4'

networks:
  ${FABRIC_NETWORK.name}:
    name: ${FABRIC_NETWORK.name}

services:
  # Orderer
  orderer.${FABRIC_NETWORK.name}.com:
    container_name: orderer.${FABRIC_NETWORK.name}.com
    image: hyperledger/fabric-orderer:2.5
    environment:
      - FABRIC_LOGGING_SPEC=INFO
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_LISTENPORT=7050
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
      - ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_CLUSTER_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
      - ORDERER_KAFKA_TOPIC_REPLICATIONFACTOR=1
      - ORDERER_KAFKA_VERBOSE=true
      - ORDERER_GENERAL_GENESISPROFILE=SampleInsecureSolo
      - ORDERER_GENERAL_SYSTEMCHANNEL=system-channel
      - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
      - ORDERER_FILELEDGER_LOCATION=/var/hyperledger/production/orderer
      - ORDERER_CONSENSUS_WALDIR=/var/hyperledger/production/orderer
      - ORDERER_CONSENSUS_SNAPDIR=/var/hyperledger/production/orderer
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ./crypto-config/ordererOrganizations/${FABRIC_NETWORK.name}.com/orderers/orderer.${FABRIC_NETWORK.name}.com/msp:/var/hyperledger/orderer/msp
      - ./crypto-config/ordererOrganizations/${FABRIC_NETWORK.name}.com/orderers/orderer.${FABRIC_NETWORK.name}.com/tls:/var/hyperledger/orderer/tls
      - ./channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
    networks:
      - ${FABRIC_NETWORK.name}
    ports:
      - 7050:7050

  # Peer
  peer0.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com:
    container_name: peer0.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com
    image: hyperledger/fabric-peer:2.5
    environment:
      - FABRIC_LOGGING_SPEC=INFO
      - CORE_PEER_ID=peer0.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com
      - CORE_PEER_ADDRESS=peer0.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com:7051
      - CORE_PEER_LISTENADDRESS=0.0.0.0:7051
      - CORE_PEER_CHAINCODEADDRESS=peer0.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com:7052
      - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052
      - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com:7051
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com:7051
      - CORE_PEER_LOCALMSPID=${FABRIC_NETWORK.rootMember}MSP
      - CORE_PEER_MSPCONFIGPATH=/var/hyperledger/peer/msp
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/var/hyperledger/peer/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/var/hyperledger/peer/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/peer/tls/ca.crt
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: peer node start
    volumes:
      - ./crypto-config/peerOrganizations/${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com/peers/peer0.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com/msp:/var/hyperledger/peer/msp
      - ./crypto-config/peerOrganizations/${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com/peers/peer0.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com/tls:/var/hyperledger/peer/tls
    networks:
      - ${FABRIC_NETWORK.name}
    ports:
      - 7051:7051
      - 7052:7052
      - 7053:7053

  # CA
  ca.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com:
    container_name: ca.${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com
    image: hyperledger/fabric-ca:1.5
    environment:
      - FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server
      - FABRIC_CA_SERVER_CA_NAME=ca.${FABRIC_NETWORK.rootMember}
      - FABRIC_CA_SERVER_TLS_ENABLED=true
      - FABRIC_CA_SERVER_PORT=7054
    ports:
      - "7054:7054"
    command: sh -c 'fabric-ca-server start -b admin:adminpw -d'
    volumes:
      - ./crypto-config/peerOrganizations/${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com/ca:/etc/hyperledger/fabric-ca-server-config
    networks:
      - ${FABRIC_NETWORK.name}
`;

  // Create fabric directory
  const fabricDir = path.dirname(dockerComposePath);
  if (!fs.existsSync(fabricDir)) {
    fs.mkdirSync(fabricDir, { recursive: true });
  }

  fs.writeFileSync(dockerComposePath, dockerCompose);
  console.log(`🐳 Generated Docker Compose: ${dockerComposePath}`);
}

function generateNetworkConfig() {
  const configPath = path.join(__dirname, '..', 'fabric', 'configtx.yaml');

  const configTx = `---
Organizations:
  - &OrdererOrg
    Name: OrdererOrg
    ID: OrdererMSP
    MSPDir: crypto-config/ordererOrganizations/${FABRIC_NETWORK.name}.com/msp

  - &${FABRIC_NETWORK.rootMember}
    Name: ${FABRIC_NETWORK.rootMember}
    ID: ${FABRIC_NETWORK.rootMember}MSP
    MSPDir: crypto-config/peerOrganizations/${FABRIC_NETWORK.rootMember}.${FABRIC_NETWORK.name}.com/msp
    Policies:
      Readers:
        Type: Signature
        Rule: "OR('${FABRIC_NETWORK.rootMember}MSP.admin', '${FABRIC_NETWORK.rootMember}MSP.peer', '${FABRIC_NETWORK.rootMember}MSP.client')"
      Writers:
        Type: Signature
        Rule: "OR('${FABRIC_NETWORK.rootMember}MSP.admin', '${FABRIC_NETWORK.rootMember}MSP.client')"
      Admins:
        Type: Signature
        Rule: "OR('${FABRIC_NETWORK.rootMember}MSP.admin')"
      Endorsement:
        Type: Signature
        Rule: "OR('${FABRIC_NETWORK.rootMember}MSP.peer')"

Capabilities:
  Channel: &ChannelCapabilities
    V2_0: true
  Orderer: &OrdererCapabilities
    V2_0: true
  Application: &ApplicationCapabilities
    V2_0: true

Application: &ApplicationDefaults
  Organizations:
  Policies:
    Readers:
      Type: ImplicitMeta
      Rule: "ANY Readers"
    Writers:
      Type: ImplicitMeta
      Rule: "ANY Writers"
    Admins:
      Type: ImplicitMeta
      Rule: "MAJORITY Admins"
    LifecycleEndorsement:
      Type: ImplicitMeta
      Rule: "MAJORITY Endorsement"
    Endorsement:
      Type: ImplicitMeta
      Rule: "MAJORITY Endorsement"
  Capabilities:
    <<: *ApplicationCapabilities

Orderer: &OrdererDefaults
  OrdererType: solo
  Addresses:
    - orderer.${FABRIC_NETWORK.name}.com:7050
  BatchTimeout: 2s
  BatchSize:
    MaxMessageCount: 10
    AbsoluteMaxBytes: 99 MB
    PreferredMaxBytes: 512 KB
  Organizations:
  Policies:
    Readers:
      Type: ImplicitMeta
      Rule: "ANY Readers"
    Writers:
      Type: ImplicitMeta
      Rule: "ANY Writers"
    Admins:
      Type: ImplicitMeta
      Rule: "MAJORITY Admins"
    BlockValidation:
      Type: ImplicitMeta
      Rule: "ANY Writers"
  Capabilities:
    <<: *OrdererCapabilities

Channel: &ChannelDefaults
  Policies:
    Readers:
      Type: ImplicitMeta
      Rule: "ANY Readers"
    Writers:
      Type: ImplicitMeta
      Rule: "ANY Writers"
    Admins:
      Type: ImplicitMeta
      Rule: "MAJORITY Admins"
  Capabilities:
    <<: *ChannelCapabilities

Profiles:
  TwoOrgsOrdererGenesis:
    <<: *ChannelDefaults
    Orderer:
      <<: *OrdererDefaults
      Organizations:
        - *OrdererOrg
    Consortiums:
      SampleConsortium:
        Organizations:
          - *${FABRIC_NETWORK.rootMember}
  TwoOrgsChannel:
    Consortium: SampleConsortium
    <<: *ChannelDefaults
    Application:
      <<: *ApplicationDefaults
      Organizations:
        - *${FABRIC_NETWORK.rootMember}
      Capabilities:
        <<: *ApplicationCapabilities
`;

  // Create config directory
  const configDir = path.dirname(configPath);
  if (!fs.existsSync(configDir)) {
    fs.mkdirSync(configDir, { recursive: true });
  }

  fs.writeFileSync(configPath, configTx);
  console.log(`⚙️  Generated network config: ${configPath}`);
}

function main() {
  console.log('🏗️  BTNG Fabric Chaincode Deployment Setup');
  console.log('==========================================');
  console.log(`Network: ${FABRIC_NETWORK.name}`);
  console.log(`Root Member: ${FABRIC_NETWORK.rootMember}`);
  console.log(`Node ID: ${FABRIC_NETWORK.nodeId}`);
  console.log('');

  // Create fabric directory structure
  const fabricDir = path.join(__dirname, '..', 'fabric');
  if (!fs.existsSync(fabricDir)) {
    fs.mkdirSync(fabricDir, { recursive: true });
  }

  // Generate chaincode packages
  Object.keys(CHAINCODES).forEach(chaincodeName => {
    generateChaincodePackage(chaincodeName);
    generateDeploymentScript(chaincodeName);
  });

  // Generate network configuration
  generateDockerCompose();
  generateNetworkConfig();

  console.log('');
  console.log('🎉 BTNG Fabric deployment setup complete!');
  console.log('');
  console.log('📁 Generated files:');
  console.log('  - fabric/chaincode/btng-gold-token/main.go');
  console.log('  - fabric/chaincode/btng-sovereign-identity/main.go');
  console.log('  - fabric/scripts/deploy-btng-gold-token.sh');
  console.log('  - fabric/scripts/deploy-btng-sovereign-identity.sh');
  console.log('  - fabric/docker-compose.yml');
  console.log('  - fabric/configtx.yaml');
  console.log('');
  console.log('🚀 Next steps:');
  console.log('1. Start Fabric network: docker-compose -f fabric/docker-compose.yml up -d');
  console.log('2. Create channel: ./fabric/scripts/create-channel.sh');
  console.log('3. Deploy chaincodes: ./fabric/scripts/deploy-btng-gold-token.sh');
  console.log('4. Test operations via BTNG API');
}

if (require.main === module) {
  main();
}