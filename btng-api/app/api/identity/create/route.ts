import { NextRequest, NextResponse } from 'next/server'
import { generateGoldCardNumber, calculateTrustScore, verifyProfile } from '@/lib/identity'
import { TrustProfile } from '@/types/trust-union'

export async function POST(request: NextRequest) {
  try {
    const body = await request.json()
    const { holderName, country, verificationType = 'basic' } = body

    // Generate new Gold Card
    const goldCardNumber = generateGoldCardNumber()
    
    // Create trust profile
    const profile: TrustProfile = {
      id: `tp_${Date.now()}`,
      holderName,
      goldCardNumber,
      issueDate: new Date().toISOString(),
      trustScore: 500, // Initial score
      verificationLevel: verificationType,
      activeCountries: [country]
    }

    // Verify profile via Trust Union Protocol
    const verified = await verifyProfile(profile.id)

    if (!verified) {
      return NextResponse.json(
        { error: 'Profile verification failed' },
        { status: 400 }
      )
    }

    // Calculate initial trust score
    profile.trustScore = calculateTrustScore(profile)

    return NextResponse.json({
      success: true,
      profile,
      goldCard: {
        number: goldCardNumber,
        holderName,
        issueDate: profile.issueDate,
        trustScore: profile.trustScore
      }
    }, {
      status: 201,
      headers: {
        'X-BTNG-Profile-Id': profile.id
      }
    })

  } catch (error) {
    console.error('Profile creation error:', error)
    return NextResponse.json(
      { error: 'Failed to create profile' },
      { status: 500 }
    )
  }
}
