const { readFileSync } = require('fs')
const { resolve } = require('path')

function getDocumentationVersionHash() {
    try {
        const identityPath = resolve(process.cwd(), 'documentation.identity.json')
        const manifest = JSON.parse(readFileSync(identityPath, 'utf8'))
        return manifest?.documentation?.versionHash || 'unavailable'
    } catch {
        return 'unavailable'
    }
}

const documentationVersionHash = getDocumentationVersionHash()

module.exports = {
    apps: [
        {
            name: 'btng-gold-oracle',
            script: 'npm',
            args: 'run oracle:push-price:ps',
            env_file: '.env.oracle',
            env: {
                BTNG_DOC_VERSION_HASH: documentationVersionHash
            },
            documentation: {
                versionHash: documentationVersionHash,
                source: 'BTNG Sovereign Documentation Suite'
            },
            restart_delay: 5000,
            max_restarts: 10,
            log_date_format: 'YYYY-MM-DD HH:mm Z',
            error_file: './logs/oracle-error.log',
            out_file: './logs/oracle-out.log',
            merge_logs: true
        }
    ]
}
