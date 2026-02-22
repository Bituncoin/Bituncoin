/**
 * BTNG Deployment Order Configuration
 * Defines the dependency chain for full-suite deployment
 */

export const BTNG_DEPLOY_ORDER = [
  {
    name: 'BTNGGoldOracle',
    contract: 'BTNGGoldOracle',
    dependencies: [],
    constructorArgs: [],
    description: 'Gold price oracle with African reserve data'
  },
  {
    name: 'BTNGGoldToken',
    contract: 'BTNGGoldToken',
    dependencies: [], // Initially uses ZeroAddress, updated later
    constructorArgs: ['ethers.ZeroAddress'], // Placeholder custody address
    description: 'Sovereign gold-backed token'
  },
  {
    name: 'BTNGCustody',
    contract: 'BTNGCustody',
    dependencies: ['BTNGGoldToken', 'BTNGGoldOracle'],
    constructorArgs: ['addresses.token', 'addresses.oracle'],
    description: 'Secure custody vault for gold reserves'
  }
];

export const BTNG_POST_DEPLOY_UPDATES = [
  {
    contract: 'BTNGGoldToken',
    function: 'updateCustodyContract',
    args: ['addresses.custody'],
    description: 'Link token to custody contract'
  }
];

export const BTNG_OWNERSHIP_TRANSFERS = [
  'BTNGGoldOracle',
  'BTNGGoldToken',
  'BTNGCustody'
];