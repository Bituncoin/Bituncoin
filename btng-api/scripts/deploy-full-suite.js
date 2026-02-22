#!/usr/bin/env node

/**
 * BTNG Full Suite Deployment with Dependency Management
 * Deploys all contracts in correct order with automatic dependency resolution
 */

const { ethers } = require("hardhat");
const { BTNG_DEPLOY_ORDER, BTNG_POST_DEPLOY_UPDATES, BTNG_OWNERSHIP_TRANSFERS } = require("./deploy-config.js");

async function main() {
  console.log("🚀 Deploying BTNG Gold System (Full Suite)...");

  // Get deployer account
  const [deployer] = await ethers.getSigners();
  console.log("Deploying contracts with account:", deployer.address);
  console.log("Account balance:", (await deployer.provider.getBalance(deployer.address)).toString());

  const addresses = {};
  const contracts = {};

  // Phase 1: Deploy contracts in dependency order
  console.log("\n📦 Phase 1: Deploying contracts...");

  for (const deployConfig of BTNG_DEPLOY_ORDER) {
    console.log(`\n🔨 Deploying ${deployConfig.name}...`);
    console.log(`   Description: ${deployConfig.description}`);

    const factory = await ethers.getContractFactory(deployConfig.contract);

    // Resolve constructor arguments
    let constructorArgs = [];
    for (const arg of deployConfig.constructorArgs) {
      if (arg === 'ethers.ZeroAddress') {
        constructorArgs.push(ethers.ZeroAddress);
      } else if (arg.startsWith('addresses.')) {
        const addressKey = arg.split('.')[1];
        if (!addresses[addressKey]) {
          throw new Error(`Missing dependency: ${addressKey} not deployed yet`);
        }
        constructorArgs.push(addresses[addressKey]);
      } else {
        constructorArgs.push(arg);
      }
    }

    console.log(`   Constructor args:`, constructorArgs);

    const contract = await factory.deploy(...constructorArgs);
    await contract.waitForDeployment();

    const address = await contract.getAddress();
    addresses[deployConfig.name.toLowerCase()] = address;
    contracts[deployConfig.name] = contract;

    console.log(`   ✅ ${deployConfig.name} deployed to: ${address}`);
  }

  // Phase 2: Post-deployment updates
  console.log("\n🔄 Phase 2: Applying post-deployment updates...");

  for (const update of BTNG_POST_DEPLOY_UPDATES) {
    console.log(`\n📝 Updating ${update.contract}...`);
    console.log(`   Function: ${update.function}`);
    console.log(`   Description: ${update.description}`);

    const contract = contracts[update.contract];
    if (!contract) {
      throw new Error(`Contract ${update.contract} not found`);
    }

    // Resolve update arguments
    let updateArgs = [];
    for (const arg of update.args) {
      if (arg.startsWith('addresses.')) {
        const addressKey = arg.split('.')[1];
        if (!addresses[addressKey]) {
          throw new Error(`Missing address for update: ${addressKey}`);
        }
        updateArgs.push(addresses[addressKey]);
      } else {
        updateArgs.push(arg);
      }
    }

    console.log(`   Args:`, updateArgs);

    const tx = await contract[update.function](...updateArgs);
    await tx.wait();

    console.log(`   ✅ ${update.contract} updated successfully`);
  }

  // Phase 3: Transfer ownerships
  console.log("\n👑 Phase 3: Transferring ownerships...");

  for (const contractName of BTNG_OWNERSHIP_TRANSFERS) {
    const contract = contracts[contractName];
    if (contract && contract.transferOwnership) {
      console.log(`   Transferring ${contractName} ownership to deployer...`);
      try {
        const tx = await contract.transferOwnership(deployer.address);
        await tx.wait();
        console.log(`   ✅ ${contractName} ownership transferred`);
      } catch (error) {
        console.log(`   ⚠️  ${contractName} ownership transfer failed (may not be needed):`, error.message);
      }
    }
  }

  // Phase 4: Verification
  console.log("\n🔍 Phase 4: Verifying deployment...");

  const token = contracts.BTNGGoldToken;
  const custody = contracts.BTNGCustody;
  const oracle = contracts.BTNGGoldOracle;

  if (token) {
    const tokenName = await token.name();
    const tokenSymbol = await token.symbol();
    const custodyContract = await token.custodyContract();
    console.log(`   Token Name: ${tokenName}`);
    console.log(`   Token Symbol: ${tokenSymbol}`);
    console.log(`   Custody Contract: ${custodyContract}`);
  }

  if (oracle) {
    const goldPrice = await oracle.getGoldPrice();
    const totalReserves = await oracle.getTotalAfricanReserves();
    const oracleOwner = await oracle.owner();
    console.log(`   Gold Price: ${goldPrice.toString()}`);
    console.log(`   Total Reserves: ${totalReserves.toString()}`);
    console.log(`   Oracle Owner: ${oracleOwner}`);
  }

  console.log("\n🎉 BTNG Gold System deployed successfully!");
  console.log("📋 Contract Addresses:");

  for (const [name, address] of Object.entries(addresses)) {
    console.log(`   - ${name}: ${address}`);
  }

  // Save deployment info
  const deploymentInfo = {
    network: (await ethers.provider.getNetwork()).name,
    addresses: addresses,
    deployer: deployer.address,
    deployedAt: new Date().toISOString(),
    type: "full-suite",
    phase: "complete"
  };

  console.log("\n💾 Deployment Info:", JSON.stringify(deploymentInfo, null, 2));

  return deploymentInfo;
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("❌ Full suite deployment failed:", error);
    console.error("Error details:", error.message);
    process.exit(1);
  });