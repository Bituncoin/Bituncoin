import { GET as getImperialReserve } from '../api/dashboard/imperial-reserve/route';

// Mock dependencies
jest.mock('@/lib/mongodb', () => ({
    __esModule: true,
    default: jest.fn(),
}));

jest.mock('@/lib/gold-price/model', () => ({
    getGoldPrice: jest.fn().mockResolvedValue({
        price_ounce: 5185.76,
        price_gram: 1832.31,
        timestamp: new Date().toISOString()
    }),
}));

jest.mock('@/lib/countries', () => ({
    getAllCountries: jest.fn().mockReturnValue(new Map([
        ['GH', { name: 'Ghana', code: 'GH' }],
        ['NG', { name: 'Nigeria', code: 'NG' }],
        ['KE', { name: 'Kenya', code: 'KE' }],
        ['ZA', { name: 'South Africa', code: 'ZA' }],
        ['EG', { name: 'Egypt', code: 'EG' }]
    ])),
}));

describe('BTNG Imperial Reserve Dashboard', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    test('returns imperial reserve data with gold backing', async () => {
        const request = new Request('http://localhost/api/dashboard/imperial-reserve');
        const response = await getImperialReserve(request);
        const data = await response.json();

        expect(response.status).toBe(200);
        expect(data).toHaveProperty('totalNations');
        expect(data).toHaveProperty('activeNations');
        expect(data).toHaveProperty('totalGoldReserve');
        expect(data).toHaveProperty('goldBackingRatio');
        expect(data).toHaveProperty('merchantCount');
        expect(data).toHaveProperty('tradeVolume24h');
        expect(data).toHaveProperty('nations');
        expect(data).toHaveProperty('currentGoldPrice');
        expect(data).toHaveProperty('message');
        expect(data.message).toContain('54 Nations Backed by Gold');
        expect(data.status).toBe('ACTIVE');
    });

    test('includes proper headers for transparency', async () => {
        const request = new Request('http://localhost/api/dashboard/imperial-reserve');
        const response = await getImperialReserve(request);

        expect(response.headers.get('X-BTNG-Dashboard')).toBe('Imperial Reserve');
        expect(response.headers.get('X-Gold-Backing')).toBeDefined();
        expect(response.headers.get('X-Active-Nations')).toBeDefined();
        expect(response.headers.get('Cache-Control')).toContain('public');
    });
});