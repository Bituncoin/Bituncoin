import { NextRequest, NextResponse } from 'next/server'
import { createValueProfile, addProofOfValue, verifyProof } from '@/lib/pov'
import { ProofOfValue } from '@/types/trust-union'

/**
 * Get proof-of-value credentials for a profile
 */
export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url)
    const profileId = searchParams.get('profileId')

    if (!profileId) {
      return NextResponse.json(
        { error: 'Profile ID is required' },
        { status: 400 }
      )
    }

    // TODO: Fetch from database
    // For now, return mock data
    const mockProfile = createValueProfile(profileId)

    return NextResponse.json({
      success: true,
      profile: mockProfile,
      proofOfValueCredentials: mockProfile.proofOfValueCredentials
    })

  } catch (error) {
    console.error('POV fetch error:', error)
    return NextResponse.json(
      { error: 'Failed to fetch proof-of-value credentials' },
      { status: 500 }
    )
  }
}

/**
 * Create a new proof-of-value credential
 */
export async function POST(request: NextRequest) {
  try {
    const body = await request.json()
    const { 
      trustProfileId, 
      type, 
      amount, 
      verifiedBy, 
      metadata 
    } = body

    // Validate required fields
    if (!trustProfileId || !type || !amount || !verifiedBy) {
      return NextResponse.json(
        { error: 'Missing required fields' },
        { status: 400 }
      )
    }

    // Validate POV type
    const validTypes = ['work', 'trade', 'trust', 'contribution']
    if (!validTypes.includes(type)) {
      return NextResponse.json(
        { error: 'Invalid proof-of-value type' },
        { status: 400 }
      )
    }

    // Create value profile or get existing
    const valueProfile = createValueProfile(trustProfileId)

    // Add proof-of-value
    const updatedProfile = addProofOfValue(
      valueProfile,
      {
        type,
        amount,
        verifiedBy,
        metadata: metadata || {}
      }
    )

    // Verify the proof
    const latestProof = updatedProfile.proofOfValueCredentials[
      updatedProfile.proofOfValueCredentials.length - 1
    ]
    const verified = await verifyProof(latestProof)

    if (!verified) {
      return NextResponse.json(
        { error: 'Proof-of-value verification failed' },
        { status: 400 }
      )
    }

    return NextResponse.json({
      success: true,
      proof: latestProof,
      totalValue: updatedProfile.totalValue,
      profileId: updatedProfile.profileId
    }, {
      status: 201,
      headers: {
        'X-POV-Id': latestProof.id
      }
    })

  } catch (error) {
    console.error('POV creation error:', error)
    return NextResponse.json(
      { error: 'Failed to create proof-of-value' },
      { status: 500 }
    )
  }
}

/**
 * Verify a proof-of-value credential
 */
export async function PUT(request: NextRequest) {
  try {
    const body = await request.json()
    const { proofId, verificationData } = body

    if (!proofId) {
      return NextResponse.json(
        { error: 'Proof ID is required' },
        { status: 400 }
      )
    }

    // TODO: Implement cryptographic verification
    // For now, return mock verification
    const verified = true

    return NextResponse.json({
      success: true,
      proofId,
      verified,
      verificationTimestamp: new Date().toISOString()
    })

  } catch (error) {
    console.error('POV verification error:', error)
    return NextResponse.json(
      { error: 'Failed to verify proof-of-value' },
      { status: 500 }
    )
  }
}
