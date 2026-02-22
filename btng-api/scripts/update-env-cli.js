#!/usr/bin/env node

/**
 * BTNG .env Updater - Command Line Version
 * Usage: node scripts/update-env-cli.js INFURA_API_KEY PRIVATE_KEY [ETHERSCAN_API_KEY]
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const envPath = path.join(__dirname, '..', '.env');
const envExamplePath = path.join(__dirname, '..', '.env.example');

function updateEnv(infuraApiKey, privateKey, etherscanKey = '') {
  let content;

  if (!fs.existsSync(envPath)) {
    // Create .env from .env.example if it doesn't exist
    if (!fs.existsSync(envExamplePath)) {
      console.error('❌ Neither .env nor .env.example found!');
      return false;
    }
    content = fs.readFileSync(envExamplePath, 'utf8');
    console.log('📄 Creating .env file from .env.example...');
  } else {
    content = fs.readFileSync(envPath, 'utf8');
  }

  // Validate inputs
  if (!infuraApiKey || infuraApiKey.length < 20) {
    console.error('❌ Invalid INFURA_API_KEY: Must be at least 20 characters');
    return false;
  }

  if (!privateKey || privateKey.length !== 64) {
    console.error('❌ Invalid PRIVATE_KEY: Must be exactly 64 characters (without 0x prefix)');
    console.log('💡 Get your private key from MetaMask: Settings → Security & Privacy → Export Private Key');
    return false;
  }

  // Remove 0x prefix if present
  privateKey = privateKey.replace(/^0x/, '');

  // Update values
  content = content.replace(/INFURA_API_KEY=.*/, `INFURA_API_KEY=${infuraApiKey}`);
  content = content.replace(/PRIVATE_KEY=.*/, `PRIVATE_KEY=${privateKey}`);
  if (etherscanKey) {
    content = content.replace(/ETHERSCAN_API_KEY=.*/, `ETHERSCAN_API_KEY=${etherscanKey}`);
  }

  fs.writeFileSync(envPath, content);

  console.log('✅ .env file updated successfully!');
  console.log('📁 File location:', envPath);
  console.log('\n🔐 Security Reminder:');
  console.log('- Never commit .env files to version control');
  console.log('- Keep your private key secure and never share it');
  console.log('- Use different keys for testnet and mainnet');

  return true;
}

// CLI usage
if (process.argv.length < 4) {
  console.log('🔧 BTNG .env Updater');
  console.log('Usage: node scripts/update-env-cli.js <INFURA_API_KEY> <PRIVATE_KEY> [ETHERSCAN_API_KEY]');
  console.log('');
  console.log('Required:');
  console.log('- INFURA_API_KEY: Your Infura API key (get from https://infura.io)');
  console.log('- PRIVATE_KEY: Your MetaMask private key (64 chars, no 0x prefix)');
  console.log('');
  console.log('Optional:');
  console.log('- ETHERSCAN_API_KEY: For contract verification (get from https://etherscan.io)');
  console.log('');
  console.log('Example:');
  console.log('node scripts/update-env-cli.js abc123def456... def456ghi789... ghi789jkl012...');
  process.exit(1);
}

const [,, infuraApiKey, privateKey, etherscanKey] = process.argv;

if (updateEnv(infuraApiKey, privateKey, etherscanKey)) {
  console.log('\n🚀 Ready for deployment! Next steps:');
  console.log('1. npm run check-env');
  console.log('2. npm run test-connection');
  console.log('3. npm run deploy:oracle');
}