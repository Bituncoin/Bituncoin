import { GET as getMerchantDeploymentPackage } from '../api/merchant/deployment-package/route';

describe('Merchant Deployment Package', () => {
    test('returns complete deployment package with sovereign config', async () => {
        const response = await getMerchantDeploymentPackage();
        expect(response.status).toBe(200);

        const data = await response.json();
        expect(data).toHaveProperty('packageName', 'BTNG Merchant App Deployment Package');
        expect(data).toHaveProperty('version', '1.0.0');
        expect(data).toHaveProperty('generatedAt');
        expect(data).toHaveProperty('targetVendors', 10);
        expect(data).toHaveProperty('distribution');
        expect(data).toHaveProperty('sovereignConfig');
        expect(data).toHaveProperty('rolloutChecklist');

        // Check sovereign config
        expect(data.sovereignConfig).toHaveProperty('parentNode', 'http://74.118.126.72:64799');
        expect(data.sovereignConfig).toHaveProperty('heartbeatIntervalMs', 15000);
        expect(data.sovereignConfig).toHaveProperty('requestFormat', '35-digit-gold-request');

        // Check rollout checklist
        expect(data.rolloutChecklist).toHaveLength(4);
        expect(data.rolloutChecklist[0]).toContain('Install app');
        expect(data.rolloutChecklist[3]).toContain('Watchtower dashboard');
    });
});