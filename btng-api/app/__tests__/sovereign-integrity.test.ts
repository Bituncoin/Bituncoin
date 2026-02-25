import { GET as getWatchtowerNodes } from '../api/watchtower/nodes/route';
import { POST as createIdentity } from '../api/identity/create/route';

// Mock dependencies
jest.mock('@/lib/documentation-identity', () => ({
    verifyDocumentationIdentity: jest.fn().mockReturnValue({
        available: true,
        verified: true,
        hashValid: true,
        signatureValid: true
    }),
}));

// Mock fetch for node pulse checks
global.fetch = jest.fn().mockResolvedValue({
    ok: true,
    json: jest.fn().mockResolvedValue({ blockHeight: 12459 })
});

jest.mock('@/lib/identity', () => ({
    generateGoldCardNumber: jest.fn().mockReturnValue('BTNG1234567890123456789012345678912'),
    calculateTrustScore: jest.fn().mockReturnValue(500),
    verifyProfile: jest.fn().mockResolvedValue(true),
}));

describe('Sovereign Integrity Suite', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    describe('Watchtower Node Testing (The Auditor)', () => {
        test('Snapshot chaining maintains hash continuity', async () => {
            // First snapshot
            const res1 = await getWatchtowerNodes();
            expect(res1.status).toBe(200);
            const data1 = await res1.json();
            expect(data1).toHaveProperty('snapshot');
            expect(data1.snapshot).toHaveProperty('snapshotHash');
            expect(data1.snapshot).toHaveProperty('snapshotSequence');

            // Second snapshot should chain from first
            const res2 = await getWatchtowerNodes();
            expect(res2.status).toBe(200);
            const data2 = await res2.json();
            expect(data2.snapshot.previousSnapshotHash).toBe(data1.snapshot.snapshotHash);
            expect(data2.snapshot.snapshotSequence).toBe(data1.snapshot.snapshotSequence + 1);
        });

        test('Snapshot sequencing works correctly', async () => {
            const res1 = await getWatchtowerNodes();
            const res2 = await getWatchtowerNodes();

            const data1 = await res1.json();
            const data2 = await res2.json();

            // Second snapshot should have higher sequence number
            expect(data2.snapshot.snapshotSequence).toBeGreaterThan(data1.snapshot.snapshotSequence);
        });
    });

    describe('Identity Creation Testing (The Gatekeeper)', () => {
        test('35-digit address format validation', async () => {
            const mockRequest = {
                json: jest.fn().mockResolvedValue({
                    holderName: 'John Doe',
                    country: 'US',
                    verificationType: 'basic'
                }),
            } as any;

            const response = await createIdentity(mockRequest);
            expect(response.status).toBe(201);

            const data = await response.json();
            expect(data.success).toBe(true);
            expect(data.goldCard.number).toHaveLength(35);
            expect(data.goldCard.number).toMatch(/^BTNG/);
            expect(data.goldCard.holderName).toBe('John Doe');
            expect(data.goldCard).toHaveProperty('trustScore');
        });

        test('Profile verification failure handling', async () => {
            // Mock verification failure
            const { verifyProfile } = require('@/lib/identity');
            verifyProfile.mockResolvedValueOnce(false);

            const mockRequest = {
                json: jest.fn().mockResolvedValue({
                    holderName: 'Jane Doe',
                    country: 'CA'
                }),
            } as any;

            const response = await createIdentity(mockRequest);
            expect(response.status).toBe(400);

            const data = await response.json();
            expect(data.error).toBe('Profile verification failed');
        });
    });
});