import { describe, it, expect } from 'vitest';
import { ApiKeyAuth } from '../src/auth/api-key.js';
import { TokenAuth } from '../src/auth/token.js';

describe('ApiKeyAuth', () => {
    const validApiKey = 'lk_test_12345678901234567890123456789012';

    it('should create auth with valid API key', () => {
        const auth = new ApiKeyAuth(validApiKey);
        expect(auth).toBeInstanceOf(ApiKeyAuth);
    });

    it('should throw error for empty API key', () => {
        expect(() => new ApiKeyAuth('')).toThrow('API key is required');
    });

    it('should throw error for invalid API key format', () => {
        expect(() => new ApiKeyAuth('invalid-key')).toThrow('Invalid API key format');
        expect(() => new ApiKeyAuth('lk_invalid_key')).toThrow('Invalid API key format');
        expect(() => new ApiKeyAuth('lk_test_short')).toThrow('Invalid API key format');
    });

    it('should return correct headers', async () => {
        const auth = new ApiKeyAuth(validApiKey);
        const headers = await auth.getHeaders();

        expect(headers).toEqual({
            'Authorization': `Bearer ${validApiKey}`,
        });
    });

    it('should not need refresh', () => {
        const auth = new ApiKeyAuth(validApiKey);
        expect(auth.needsRefresh()).toBe(false);
    });
});

describe('TokenAuth', () => {
    const accessToken = 'test-access-token';
    const refreshToken = 'test-refresh-token';

    it('should create auth with access token', () => {
        const auth = new TokenAuth({ accessToken });
        expect(auth).toBeInstanceOf(TokenAuth);
    });

    it('should throw error for empty access token', () => {
        expect(() => new TokenAuth({ accessToken: '' })).toThrow('Access token is required');
    });

    it('should return correct headers', async () => {
        const auth = new TokenAuth({ accessToken });
        const headers = await auth.getHeaders();

        expect(headers).toEqual({
            Authorization: `Bearer ${accessToken}`,
        });
    });

    it('should not need refresh without expiry', () => {
        const auth = new TokenAuth({ accessToken });
        expect(auth.needsRefresh()).toBe(false);
    });

    it('should update tokens', () => {
        const auth = new TokenAuth({ accessToken, refreshToken });

        auth.updateTokens({
            accessToken: 'new-access-token',
            refreshToken: 'new-refresh-token',
            expiresAt: new Date(Date.now() + 3600000),
        });

        // Token was updated (can verify via getHeaders)
        expect(auth.needsRefresh()).toBe(false);
    });
});
