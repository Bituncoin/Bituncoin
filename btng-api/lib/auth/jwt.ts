import { NextRequest, NextResponse } from 'next/server';
import jwt from 'jsonwebtoken';

// JWT Secret - In production, use environment variable
const JWT_SECRET = process.env.JWT_SECRET || 'a-string-secret-at-least-256-bits-long';

export interface JWTPayload {
  sub: string;
  name: string;
  admin: boolean;
  iat: number;
  exp: number;
}

export function authenticateToken(request: NextRequest): {
  isValid: boolean;
  user?: JWTPayload;
  error?: string;
} {
  try {
    // Get token from Authorization header
    const authHeader = request.headers.get('authorization');
    const token = authHeader?.replace('Bearer ', '');

    if (!token) {
      return { isValid: false, error: 'Access token required' };
    }

    // Verify the token
    const decoded = jwt.verify(token, JWT_SECRET, {
      algorithms: ['HS256']
    }) as JWTPayload;

    return { isValid: true, user: decoded };
  } catch (error: any) {
    if (error.name === 'TokenExpiredError') {
      return { isValid: false, error: 'Token expired' };
    } else if (error.name === 'JsonWebTokenError') {
      return { isValid: false, error: 'Invalid token' };
    } else {
      return { isValid: false, error: 'Token verification failed' };
    }
  }
}

export function requireAuth(handler: (request: NextRequest, user: JWTPayload) => Promise<NextResponse>) {
  return async (request: NextRequest): Promise<NextResponse> => {
    const auth = authenticateToken(request);

    if (!auth.isValid) {
      return NextResponse.json(
        { error: auth.error },
        { status: 401 }
      );
    }

    return handler(request, auth.user!);
  };
}

export function requireAdmin(handler: (request: NextRequest, user: JWTPayload) => Promise<NextResponse>) {
  return async (request: NextRequest): Promise<NextResponse> => {
    const auth = authenticateToken(request);

    if (!auth.isValid) {
      return NextResponse.json(
        { error: auth.error },
        { status: 401 }
      );
    }

    if (!auth.user!.admin) {
      return NextResponse.json(
        { error: 'Admin access required' },
        { status: 403 }
      );
    }

    return handler(request, auth.user!);
  };
}