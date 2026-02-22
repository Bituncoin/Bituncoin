import { NextRequest, NextResponse } from 'next/server';
import jwt from 'jsonwebtoken';

// JWT Secret - In production, use environment variable
const JWT_SECRET = process.env.JWT_SECRET || 'a-string-secret-at-least-256-bits-long';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const { username, password } = body;

    // Basic authentication (replace with your actual auth logic)
    // For demo purposes, accepting any username/password
    // In production, validate against database/users
    if (!username || !password) {
      return NextResponse.json(
        { error: 'Username and password required' },
        { status: 400 }
      );
    }

    // Create JWT payload
    const payload = {
      sub: username,
      name: username,
      admin: true, // Set based on user role
      iat: Math.floor(Date.now() / 1000),
      exp: Math.floor(Date.now() / 1000) + (24 * 60 * 60) // 24 hours
    };

    // Sign the token
    const token = jwt.sign(payload, JWT_SECRET, {
      algorithm: 'HS256'
    });

    return NextResponse.json({
      status: 'ok',
      token: token,
      user: {
        username: username,
        admin: payload.admin
      },
      message: 'Login successful'
    }, {
      status: 200,
      headers: {
        'X-BTNG-Platform': 'Sovereign',
        'X-Auth-Status': 'Authenticated'
      }
    });

  } catch (error: any) {
    console.error('BTNG login error:', error.message);
    return NextResponse.json(
      { error: 'Authentication failed', details: error.message },
      { status: 500 }
    );
  }
}