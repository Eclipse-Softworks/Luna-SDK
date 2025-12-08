import { AuthenticationError } from '../errors/base.js';
import { ErrorCode } from '../errors/codes.js';
import type { AuthProvider } from './types.js';

/**
 * API Key authentication provider
 *
 * @example
 * ```typescript
 * const auth = new ApiKeyAuth('lk_live_xxxx');
 * const headers = await auth.getHeaders();
 * // { 'X-Luna-Api-Key': 'lk_live_xxxx' }
 * ```
 */
export class ApiKeyAuth implements AuthProvider {
    private readonly apiKey: string;

    constructor(apiKey: string) {
        if (!apiKey) {
            throw new AuthenticationError({
                code: ErrorCode.AUTH_INVALID_KEY,
                message: 'API key is required',
                requestId: 'local',
            });
        }
        if (!this.isValidApiKey(apiKey)) {
            throw new AuthenticationError({
                code: ErrorCode.AUTH_INVALID_KEY,
                message: 'Invalid API key format. Expected: lk_<env>_<key>',
                requestId: 'local',
            });
        }
        this.apiKey = apiKey;
    }

    async getHeaders(): Promise<Record<string, string>> {
        return {
            'Authorization': `Bearer ${this.apiKey}`,
        };
    }

    needsRefresh(): boolean {
        return false; // API keys don't expire
    }

    /**
     * Validate API key format
     * Format: lk_<environment>_<32 chars>
     */
    private isValidApiKey(key: string): boolean {
        const pattern = /^lk_(live|test|dev)_[a-zA-Z0-9]{32}$/;
        return pattern.test(key);
    }
}
