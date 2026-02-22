/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  poweredByHeader: false,
  compress: true,

  // Enable standalone output for Docker
  output: 'standalone',

  // BTNG Sovereign Configuration
  env: {
    BTNG_VERSION: '0.1.0',
    BTNG_ENVIRONMENT: process.env.NODE_ENV || 'development',
  },

  // Webpack configuration for path mapping
  webpack: (config) => {
    config.resolve.alias = {
      ...config.resolve.alias,
      '@': require('path').join(__dirname),
    };
    return config;
  },

  // Future expansion: API routes for trust-union protocol
  async headers() {
    return [
      {
        source: '/api/:path*',
        headers: [
          { key: 'X-BTNG-Platform', value: 'Sovereign' },
          { key: 'X-Content-Type-Options', value: 'nosniff' },
          { key: 'X-Frame-Options', value: 'DENY' },
        ],
      },
    ];
  },
}

module.exports = nextConfig
