import { NextResponse } from 'next/server'

export async function GET() {
    const packagePayload = {
        packageName: 'BTNG Merchant App Deployment Package',
        version: '1.0.0',
        generatedAt: new Date().toISOString(),
        targetVendors: 10,
        distribution: {
            androidApkUrl: process.env.BTNG_MERCHANT_ANDROID_APK_URL || 'https://download.btng.global/merchant/android/latest.apk',
            webAppUrl: process.env.BTNG_MERCHANT_WEB_URL || 'https://merchant.btng.global',
            checksum: process.env.BTNG_MERCHANT_PACKAGE_CHECKSUM || 'pending-signature',
            installGuideUrl: process.env.BTNG_MERCHANT_INSTALL_GUIDE_URL || 'https://docs.btng.global/merchant-install'
        },
        sovereignConfig: {
            parentNode: 'http://74.118.126.72:64799',
            heartbeatIntervalMs: 15000,
            requestFormat: '35-digit-gold-request',
            oraclePath: '/api/gold-coin',
            paymentSettlementPath: '/api/btng/wallet/send'
        },
        rolloutChecklist: [
            'Install app on merchant Android devices',
            'Validate MoMo settlement account binding',
            'Run one live 35-digit payment request test',
            'Confirm ledger update on Watchtower dashboard'
        ]
    }

    return NextResponse.json(packagePayload, {
        headers: {
            'Cache-Control': 'no-store',
            'Content-Disposition': 'attachment; filename="btng-merchant-deployment-package.json"'
        }
    })
}
