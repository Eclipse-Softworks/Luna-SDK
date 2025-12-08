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
            throw new Error('API key is required');
        }
        if (!this.isValidApiKey(apiKey)) {
            throw new Error('Invalid API key format. Expected: lk_<env>_<key>');
        }
        this.apiKey = apiKey;
    }

    async getHeaders(): Promise<Record<string, string>> {
        return {
            'X-Luna-Api-Key': this.apiKey,
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
