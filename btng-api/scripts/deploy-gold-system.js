const { ethers } = require("hardhat");

async function main() {
  console.log("🚀 Deploying BTNG Gold System...");

  // Get deployer account
  const [deployer] = await ethers.getSigners();
  console.log("Deploying contracts with account:", deployer.address);
  console.log("Account balance:", (await deployer.provider.getBalance(deployer.address)).toString());

  // Deploy Oracle
  console.log("📡 Deploying BTNG Gold Oracle...");
  const BTNGGoldOracle = await ethers.getContractFactory("BTNGGoldOracle");
  const oracle = await BTNGGoldOracle.deploy();
  await oracle.waitForDeployment();
  console.log("✅ Oracle deployed to:", await oracle.getAddress());

  // Deploy Token with placeholder custody
  console.log("🪙 Deploying BTNG Gold Token...");
  const BTNGGoldToken = await ethers.getContractFactory("BTNGGoldToken");
  const goldToken = await BTNGGoldToken.deploy(ethers.ZeroAddress);
  await goldToken.waitForDeployment();
  console.log("✅ Gold Token deployed to:", await goldToken.getAddress());

  // Deploy Custody with real token address
  console.log("🏦 Deploying BTNG Custody...");
  const BTNGCustody = await ethers.getContractFactory("BTNGCustody");
  const custody = await BTNGCustody.deploy(await goldToken.getAddress(), await oracle.getAddress());
  await custody.waitForDeployment();
  console.log("✅ Custody deployed to:", await custody.getAddress());

  // Update Token with custody address
  console.log("🔄 Updating Token contract with custody address...");
  await goldToken.updateCustodyContract(await custody.getAddress());

  // Transfer ownerships to deployer
  console.log("👑 Transferring ownerships...");
  await oracle.transferOwnership(deployer.address);
  await custody.transferOwnership(deployer.address);
  await goldToken.transferOwnership(deployer.address);

  // Setup initial custodians and feeders
  console.log("🔧 Setting up initial roles...");
  // Roles already set in constructors

  // Verify deployment
  console.log("🔍 Verifying deployment...");
  const tokenName = await goldToken.name();
  const tokenSymbol = await goldToken.symbol();
  const custodyContract = await goldToken.custodyContract();
  const oracleOwner = await oracle.owner();

  console.log("Token Name:", tokenName);
  console.log("Token Symbol:", tokenSymbol);
  console.log("Custody Contract:", custodyContract);
  console.log("Oracle Owner:", oracleOwner);

  // Test basic functionality
  console.log("🧪 Testing basic functionality...");
  const goldPrice = await oracle.getGoldPrice();
  console.log("Current Gold Price:", goldPrice.toString());

  const totalReserves = await oracle.getTotalAfricanReserves();
  console.log("Total African Gold Reserves:", totalReserves.toString());

  console.log("\n🎉 BTNG Gold System deployed successfully!");
  console.log("📋 Contract Addresses:");
  console.log("- BTNG Gold Token:", await goldToken.getAddress());
  console.log("- BTNG Custody:", await custody.getAddress());
  console.log("- BTNG Oracle:", await oracle.getAddress());

  // Save deployment info
  const deploymentInfo = {
    network: network.name,
    goldToken: await goldToken.getAddress(),
    custody: await custody.getAddress(),
    oracle: await oracle.getAddress(),
    deployer: deployer.address,
    deployedAt: new Date().toISOString(),
  };

  console.log("\n💾 Deployment Info:", JSON.stringify(deploymentInfo, null, 2));

  return deploymentInfo;
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("❌ Deployment failed:", error);
    process.exit(1);
  });