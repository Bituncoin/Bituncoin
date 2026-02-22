import GoldPriceDisplay from '@/components/GoldPriceDisplay';

export default function GoldPriceDemo() {
  return (
    <div style={{
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #0a0a0a 0%, #1a1a1a 100%)',
      padding: '40px 20px',
      fontFamily: 'system-ui, -apple-system, sans-serif'
    }}>
      <div style={{ maxWidth: '800px', margin: '0 auto' }}>
        <header style={{
          textAlign: 'center',
          marginBottom: '40px',
          color: 'white'
        }}>
          <h1 style={{
            fontSize: '2.5em',
            margin: '0 0 10px 0',
            background: 'linear-gradient(45deg, #FFD700, #FFA500)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            backgroundClip: 'text'
          }}>
            🇰🇪 BTNG Sovereign Gold Standard
          </h1>
          <p style={{
            fontSize: '1.2em',
            opacity: 0.8,
            margin: 0
          }}>
            Real-time Gold Price Monitoring
          </p>
        </header>

        <GoldPriceDisplay />

        <div style={{
          marginTop: '40px',
          padding: '20px',
          background: 'rgba(255, 255, 255, 0.05)',
          borderRadius: '12px',
          color: 'white',
          fontSize: '0.9em'
        }}>
          <h3 style={{ margin: '0 0 15px 0', color: '#FFD700' }}>
            🔧 Demo Features
          </h3>
          <ul style={{ margin: 0, paddingLeft: '20px' }}>
            <li>Real-time price polling every 10 seconds</li>
            <li>Multi-currency support (when available)</li>
            <li>Connection status monitoring</li>
            <li>Error handling and recovery</li>
            <li>Responsive design for all devices</li>
          </ul>

          <h4 style={{ margin: '20px 0 10px 0', color: '#FFD700' }}>
            🚀 API Endpoints Tested
          </h4>
          <ul style={{ margin: 0, paddingLeft: '20px' }}>
            <li><code>GET /api/health</code> - System health check</li>
            <li><code>POST /api/auth/login</code> - JWT authentication</li>
            <li><code>POST /api/btng/gold/price</code> - Price updates (admin only)</li>
            <li><code>GET /api/btng/gold/price/status</code> - Current price & status</li>
          </ul>
        </div>
      </div>
    </div>
  );
}