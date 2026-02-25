import { GET } from '../api/health/route';

describe('Health endpoint', () => {
    it('returns 200 OK with operational status', async () => {
        const response = await GET();
        expect(response.status).toBe(200);

        const data = await response.json();
        expect(data.status).toBe('operational');
        expect(data).toHaveProperty('timestamp');
        expect(data).toHaveProperty('version');
        expect(data).toHaveProperty('services');
    });
});