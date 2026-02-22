#!/usr/bin/env node

/**
 * BTNG .env Updater
 * Helps update .env file with real credentials
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import readline from 'readline';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const envPath = path.join(__dirname, '..', '.env');

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

function askQuestion(question) {
  return new Promise((resolve) => {
    rl.question(question, (answer) => {
      resolve(answer.trim());
    });
  });
}

async function updateEnv() {
  console.log('🔧 BTNG .env File Updater\n');

  // Read current .env
  let envContent = '';
  if (fs.existsSync(envPath)) {
    envContent = fs.readFileSync(envPath, 'utf8');
  }

  console.log('Current .env status:');
  console.log(envContent);
  console.log('');

  // Get new values
  console.log('Enter your real credentials (press Enter to keep current value):\n');

  const infuraId = await askQuestion('Infura Project ID: ');
  const privateKey = await askQuestion('Private Key (without 0x): ');
  const etherscanKey = await askQuestion('Etherscan API Key (optional): ');

  // Update content
  let newContent = envContent;

  if (infuraId) {
    newContent = newContent.replace(
      /INFURA_PROJECT_ID=.*/,
      `INFURA_PROJECT_ID=${infuraId}`
    );
  }

  if (privateKey) {
    newContent = newContent.replace(
      /PRIVATE_KEY=.*/,
      `PRIVATE_KEY=${privateKey}`
    );
  }

  if (etherscanKey) {
    newContent = newContent.replace(
      /ETHERSCAN_API_KEY=.*/,
      `ETHERSCAN_API_KEY=${etherscanKey}`
    );
  }

  // Write back
  fs.writeFileSync(envPath, newContent);

  console.log('\n✅ .env file updated!');
  console.log('New content:');
  console.log(newContent);

  rl.close();
}

if (import.meta.url === `file://${process.argv[1]}`) {
  updateEnv();
}