import { handleGoldPriceUpdate } from '../api/btng/gold/price/route';

// Mock the dependencies
jest.mock('@/lib/gold-price/model', () => ({
    saveGoldPrice: jest.fn(),
    setLatestGoldPrice: jest.fn(),
}));

jest.mock('@/lib/mongodb', () => ({
    __esModule: true,
    default: jest.fn(),
}));

jest.mock('@/lib/fabric/chaincode', () => ({
    mintGoldTokensOnPriceUpdate: jest.fn(),
}));

describe('BTNG Gold Price Update', () => {
    const mockUser = { sub: 'admin-user' };

    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('returns 400 for missing base_price_gram', async () => {
        const mockRequest = {
            json: jest.fn().mockResolvedValue({
                base_price_ounce: 2000,
                base_price_kilo: 2000000,
            }),
        } as any;

        const response = await handleGoldPriceUpdate(mockRequest, mockUser);
        expect(response.status).toBe(400);
        const data = await response.json();
        expect(data.error).toContain('Missing required base gold price fields');
    });

    it('returns 400 for negative base_price_gram', async () => {
        const mockRequest = {
            json: jest.fn().mockResolvedValue({
                base_price_gram: -50,
                base_price_ounce: 2000,
                base_price_kilo: 2000000,
            }),
        } as any;

        const response = await handleGoldPriceUpdate(mockRequest, mockUser);
        expect(response.status).toBe(400);
        const data = await response.json();
        expect(data.error).toContain('Base gold prices must be positive numbers');
    });

    it('successfully updates gold price with valid data', async () => {
        const mockRequest = {
            json: jest.fn().mockResolvedValue({
                base_price_gram: 50,
                base_price_ounce: 2000,
                base_price_kilo: 2000000,
                currencies: ['EUR', 'GBP'],
                fx_rates: { EUR: 0.85, GBP: 0.75 },
                bid: 49.5,
                ask: 50.5,
                spread: 1.0,
                timestamp: new Date().toISOString(),
            }),
        } as any;

        const { saveGoldPrice, setLatestGoldPrice } = require('@/lib/gold-price/model');
        const { mintGoldTokensOnPriceUpdate } = require('@/lib/fabric/chaincode');

        saveGoldPrice.mockResolvedValue({ _id: 'mock-id' });
        setLatestGoldPrice.mockResolvedValue({});
        mintGoldTokensOnPriceUpdate.mockResolvedValue({ success: true, transactionId: 'tx-123' });

        const response = await handleGoldPriceUpdate(mockRequest, mockUser);
        expect(response.status).toBe(200);
        const data = await response.json();
        expect(data.status).toBe('ok');
        expect(data.token_minted).toBe(true);
        expect(data.transaction_id).toBe('tx-123');
        expect(data.message).toContain('and tokens minted');

        expect(response.headers.get('X-BTNG-Platform')).toBe('Sovereign');
        expect(response.headers.get('X-Service-Status')).toBe('Operational');
        expect(response.headers.get('X-Auth-User')).toBe(mockUser.sub);
        expect(response.headers.get('X-Token-Minted')).toBe('true');
    });

    it('succeeds even when Fabric minting fails', async () => {
        const mockRequest = {
            json: jest.fn().mockResolvedValue({
                base_price_gram: 55,
                base_price_ounce: 2100,
                base_price_kilo: 2100000,
            }),
        } as any;

        const { saveGoldPrice, setLatestGoldPrice } = require('@/lib/gold-price/model');
        const { mintGoldTokensOnPriceUpdate } = require('@/lib/fabric/chaincode');

        saveGoldPrice.mockResolvedValue({ _id: 'mock-id' });
        setLatestGoldPrice.mockResolvedValue({});
        mintGoldTokensOnPriceUpdate.mockResolvedValue({ success: false, error: 'fabric unavailable' });

        const response = await handleGoldPriceUpdate(mockRequest, mockUser);
        expect(response.status).toBe(200);
        const data = await response.json();
        expect(data.status).toBe('ok');
        expect(data.token_minted).toBe(false);
        expect(data.message).toContain('token minting failed');
        expect(response.headers.get('X-Token-Minted')).toBe('false');
    });
});