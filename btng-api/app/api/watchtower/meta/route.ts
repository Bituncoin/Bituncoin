import { NextResponse } from 'next/server'
import { verifyDocumentationIdentity } from '../../../../lib/documentation-identity'

export const runtime = 'nodejs'

export async function GET() {
    const documentationIdentity = verifyDocumentationIdentity()

    return NextResponse.json(
        {
            ok: true,
            watchtower: {
                role: 'documentation-authority-node',
                network: 'BTNG Sovereign Constellation',
                generatedAt: new Date().toISOString()
            },
            documentation: {
                available: documentationIdentity.available,
                verified: documentationIdentity.verified,
                hashValid: documentationIdentity.hashValid,
                signatureValid: documentationIdentity.signatureValid,
                versionHash: documentationIdentity.documentation?.versionHash || null,
                timestamp: documentationIdentity.documentation?.timestamp || null,
                source: documentationIdentity.documentation?.source || 'BTNG Sovereign Documentation Suite',
                signatureAlgorithm: documentationIdentity.documentation?.signatureAlgorithm || 'ES256/SHA256',
                fileCount: documentationIdentity.documentation?.files?.length || 0,
                error: documentationIdentity.error
            }
        },
        {
            headers: {
                'Cache-Control': 'no-store'
            }
        }
    )
}
