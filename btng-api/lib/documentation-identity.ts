import { createHash, createVerify } from 'crypto'
import { readFileSync } from 'fs'
import { resolve } from 'path'

export type DocumentationFileEntry = {
    path: string
    sha256: string
    size: number
}

export type DocumentationIdentity = {
    versionHash: string
    signature: string
    publicKey: string
    timestamp: number
    source: string
    hashAlgorithm: string
    signatureAlgorithm: string
    files: DocumentationFileEntry[]
}

export type DocumentationIdentityVerification = {
    available: boolean
    verified: boolean
    hashValid: boolean
    signatureValid: boolean
    error?: string
    documentation?: DocumentationIdentity
}

function sortKeysRecursively(value: unknown): unknown {
    if (Array.isArray(value)) {
        return value.map((entry) => sortKeysRecursively(entry))
    }

    if (value && typeof value === 'object') {
        const sortedObject: Record<string, unknown> = {}
        Object.keys(value as Record<string, unknown>)
            .sort()
            .forEach((key) => {
                sortedObject[key] = sortKeysRecursively((value as Record<string, unknown>)[key])
            })
        return sortedObject
    }

    return value
}

function toCanonicalJson(value: unknown): string {
    return JSON.stringify(sortKeysRecursively(value))
}

function sha256Hex(value: string): string {
    return createHash('sha256').update(value).digest('hex')
}

function verifySignature(versionHash: string, signature: string, publicKey: string): boolean {
    const verifier = createVerify('SHA256')
    verifier.update(versionHash)
    verifier.end()
    return verifier.verify(publicKey, signature, 'base64')
}

function loadDocumentationIdentity(): DocumentationIdentity | null {
    const candidatePaths = [
        resolve(process.cwd(), 'documentation.identity.json'),
        resolve(process.cwd(), '..', 'documentation.identity.json')
    ]

    for (const path of candidatePaths) {
        try {
            const raw = readFileSync(path, 'utf8')
            const parsed = JSON.parse(raw)
            if (parsed?.documentation?.versionHash) {
                return parsed.documentation as DocumentationIdentity
            }
        } catch {
            // continue
        }
    }

    return null
}

export function verifyDocumentationIdentity(): DocumentationIdentityVerification {
    const documentation = loadDocumentationIdentity()

    if (!documentation) {
        return {
            available: false,
            verified: false,
            hashValid: false,
            signatureValid: false,
            error: 'documentation.identity.json not found.'
        }
    }

    try {
        const core = {
            source: documentation.source,
            hashAlgorithm: documentation.hashAlgorithm,
            files: documentation.files || []
        }

        const computedVersionHash = `sha256-${sha256Hex(toCanonicalJson(core))}`
        const hashValid = computedVersionHash === documentation.versionHash
        const signatureValid = verifySignature(documentation.versionHash, documentation.signature, documentation.publicKey)

        return {
            available: true,
            verified: hashValid && signatureValid,
            hashValid,
            signatureValid,
            documentation
        }
    } catch (error: any) {
        return {
            available: true,
            verified: false,
            hashValid: false,
            signatureValid: false,
            error: error?.message || 'Unknown documentation identity verification error.',
            documentation
        }
    }
}
