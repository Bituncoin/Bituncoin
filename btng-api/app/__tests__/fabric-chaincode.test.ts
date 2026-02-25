import { GET as getChaincodes, POST as invokeChaincode } from '../api/btng/fabric/chaincode/route';

// Mock dependencies
jest.mock('@/lib/auth/jwt', () => ({
    requireAuth: jest.fn((handler) => async (req: any) => {
        const mockUser = { sub: 'test-user' };
        return handler(req, mockUser);
    }),
}));

jest.mock('@/lib/mongodb', () => ({
    __esModule: true,
    default: jest.fn(),
}));

describe('Fabric Chaincode Testing (Blockchain Operations)', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    test('GET lists all supported chaincodes', async () => {
        const response = await getChaincodes();
        expect(response.status).toBe(200);

        const data = await response.json();
        expect(data).toHaveProperty('network');
        expect(data).toHaveProperty('chaincodes');
        expect(data).toHaveProperty('timestamp');
        expect(data.network.name).toBe('btng-fabric-network');
        expect(data.chaincodes).toHaveLength(2);
        expect(data.chaincodes[0]).toHaveProperty('name', 'btng-gold-token');
        expect(data.chaincodes[1]).toHaveProperty('name', 'btng-sovereign-identity');

        // Check headers
        expect(response.headers.get('X-BTNG-Platform')).toBe('Sovereign');
        expect(response.headers.get('X-Fabric-Network')).toBe('btng-fabric-network');
    });

    test('POST invokes gold token mint function', async () => {
        const mockRequest = {
            json: jest.fn().mockResolvedValue({
                chaincode: 'btng-gold-token',
                function: 'Mint',
                args: ['100', 'BTNG12345678901234567890123456789012345']
            }),
        } as any;

        const response = await invokeChaincode(mockRequest);
        expect(response.status).toBe(200);

        const data = await response.json();
        expect(data.status).toBe('ok');
        expect(data.operation.chaincode).toBe('btng-gold-token');
        expect(data.operation.function).toBe('Mint');
        expect(data.operation.result).toHaveProperty('success', true);
        expect(data.operation.transaction).toHaveProperty('txId');
        expect(data.operation.transaction).toHaveProperty('blockNumber');

        // Check headers
        expect(response.headers.get('X-Chaincode')).toBe('btng-gold-token');
        expect(response.headers.get('X-Function')).toBe('Mint');
    });

    test('POST rejects unsupported chaincode', async () => {
        const mockRequest = {
            json: jest.fn().mockResolvedValue({
                chaincode: 'unsupported-chaincode',
                function: 'Mint',
                args: ['100']
            }),
        } as any;

        const response = await invokeChaincode(mockRequest);
        expect(response.status).toBe(400);

        const data = await response.json();
        expect(data.error).toContain('Unsupported chaincode');
    });

    test('POST rejects unsupported function', async () => {
        const mockRequest = {
            json: jest.fn().mockResolvedValue({
                chaincode: 'btng-gold-token',
                function: 'UnsupportedFunction',
                args: ['100']
            }),
        } as any;

        const response = await invokeChaincode(mockRequest);
        expect(response.status).toBe(400);

        const data = await response.json();
        expect(data.error).toContain('Unsupported function');
    });
});