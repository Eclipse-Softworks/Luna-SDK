import type { AuthProvider, TokenPair, TokenRefreshCallback } from './types.js';

interface TokenAuthConfig {
    accessToken: string;
    refreshToken?: string;
    expiresAt?: Date;
    onRefresh?: TokenRefreshCallback;
}

/**
 * Token-based authentication provider with automatic refresh
 *
 * @example
 * ```typescript
 * const auth = new TokenAuth({
 *   accessToken: session.accessToken,
 *   refreshToken: session.refreshToken,
 *   onRefresh: async (tokens) => {
 *     await saveToDatabase(tokens);
 *   },
 * });
 * ```
 */
export class TokenAuth implements AuthProvider {
    private accessToken: string;
    private refreshToken?: string;
    private expiresAt?: Date;
    private readonly onRefresh?: TokenRefreshCallback;
    private refreshPromise?: Promise<void>;

    constructor(config: TokenAuthConfig) {
        if (!config.accessToken) {
            throw new Error('Access token is required');
        }
        this.accessToken = config.accessToken;
        this.refreshToken = config.refreshToken;
        this.expiresAt = config.expiresAt;
        this.onRefresh = config.onRefresh;
    }

    async getHeaders(): Promise<Record<string, string>> {
        // Ensure token is valid before returning headers
        if (this.needsRefresh()) {
            await this.refresh();
        }

        return {
            Authorization: `Bearer ${this.accessToken}`,
        };
    }

    needsRefresh(): boolean {
        if (!this.expiresAt) {
            return false;
        }
        // Refresh if expiring within 5 minutes
        const bufferMs = 5 * 60 * 1000;
        return Date.now() + bufferMs >= this.expiresAt.getTime();
    }

    async refresh(): Promise<void> {
        // Prevent concurrent refresh requests
        if (this.refreshPromise) {
            return this.refreshPromise;
        }

        if (!this.refreshToken) {
            throw new Error('No refresh token available');
        }

        this.refreshPromise = this.performRefresh();

        try {
            await this.refreshPromise;
        } finally {
            this.refreshPromise = undefined;
        }
    }

    private async performRefresh(): Promise<void> {
        // Call refresh endpoint
        const response = await fetch('https://api.eclipse.dev/v1/auth/refresh', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                refresh_token: this.refreshToken,
            }),
        });

        if (!response.ok) {
            throw new Error(`Token refresh failed: ${response.status}`);
        }

        const data = (await response.json()) as {
            access_token: string;
            refresh_token: string;
            expires_in: number;
        };

        // Update tokens
        this.accessToken = data.access_token;
        this.refreshToken = data.refresh_token;
        this.expiresAt = new Date(Date.now() + data.expires_in * 1000);

        // Notify callback
        if (this.onRefresh) {
            const tokenPair: TokenPair = {
                accessToken: this.accessToken,
                refreshToken: this.refreshToken,
                expiresAt: this.expiresAt,
            };
            await this.onRefresh(tokenPair);
        }
    }

    /**
     * Manually update tokens (e.g., after re-authentication)
     */
    updateTokens(tokens: TokenPair): void {
        this.accessToken = tokens.accessToken;
        this.refreshToken = tokens.refreshToken;
        this.expiresAt = tokens.expiresAt;
    }
}
