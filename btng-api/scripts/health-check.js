/**
 * BTNG Platform Health Check Script
 * Monitors platform services and operational status
 */

const services = [
  { name: 'Platform Core', endpoint: '/api/health' },
  { name: 'Identity Service', endpoint: '/api/identity/create' },
  { name: 'Wallet Service', endpoint: '/api/wallet/transaction' }
]

async function checkHealth() {
  console.log('🔍 BTNG Platform Health Check')
  console.log('=' .repeat(50))
  
  const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000'
  
  for (const service of services) {
    try {
      const response = await fetch(`${baseUrl}${service.endpoint}`, {
        method: service.endpoint === '/api/health' ? 'GET' : 'HEAD'
      })
      
      const status = response.ok ? '✅ Healthy' : '⚠️  Degraded'
      console.log(`${service.name}: ${status}`)
      
    } catch (error) {
      console.log(`${service.name}: ❌ Down`)
    }
  }
  
  console.log('=' .repeat(50))
  console.log(`Timestamp: ${new Date().toISOString()}`)
}

checkHealth().catch(console.error)
