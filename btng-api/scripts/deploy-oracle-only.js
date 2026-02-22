#!/usr/bin/env node

/**
 * BTNG Single Contract Deployment - Oracle First
 * Deploy just the BTNGGoldOracle for initial testing
 */

import hre from "hardhat";
import dotenv from "dotenv";

dotenv.config();

async function main() {
  console.log("🚀 Deploying BTNG Gold Oracle (Single Contract)...");

  // Get deployer account - manual setup for ES modules
  const provider = new hre.ethers.JsonRpcProvider("https://ethereum-sepolia.publicnode.com");
  const wallet = new hre.ethers.Wallet(`0x${process.env.PRIVATE_KEY}`, provider);
  console.log("Deploying with account:", wallet.address);
  console.log("Account balance:", (await provider.getBalance(wallet.address)).toString());

  // Deploy Oracle only
  console.log("📡 Deploying BTNG Gold Oracle...");
  const BTNGGoldOracle = await hre.ethers.getContractFactory("BTNGGoldOracle", wallet);
  const oracle = await BTNGGoldOracle.deploy();
  await oracle.waitForDeployment();
  const oracleAddress = await oracle.getAddress();
  console.log("✅ Oracle deployed to:", oracleAddress);

  // Basic functionality test
  console.log("🧪 Testing Oracle functionality...");
  const goldPrice = await oracle.getGoldPrice();
  console.log("Current Gold Price:", goldPrice.toString());

  const totalReserves = await oracle.getTotalAfricanReserves();
  console.log("Total African Gold Reserves:", totalReserves.toString());

  // Transfer ownership
  console.log("👑 Transferring ownership...");
  await oracle.transferOwnership(deployer.address);

  console.log("\n🎉 BTNG Gold Oracle deployed successfully!");
  console.log("📋 Contract Address:", oracleAddress);

  // Save deployment info
  const deploymentInfo = {
    network: (await hre.ethers.provider.getNetwork()).name,
    oracle: oracleAddress,
    deployer: deployer.address,
    deployedAt: new Date().toISOString(),
    type: "single-contract-oracle"
  };

  console.log("\n💾 Deployment Info:", JSON.stringify(deploymentInfo, null, 2));

  return deploymentInfo;
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("❌ Oracle deployment failed:", error);
    process.exit(1);
  });