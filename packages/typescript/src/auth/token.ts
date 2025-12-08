import { AuthenticationError } from '../errors/base.js';
import { ErrorCode } from '../errors/codes.js';
import type { AuthProvider, TokenRefreshCallback } from './types.js';

export class TokenAuth implements AuthProvider {
    private accessToken: string;
    private refreshToken?: string;
    private expiresAt?: Date;
    private readonly onRefresh?: TokenRefreshCallback;
    private refreshPromise?: Promise<void>;

    constructor(config: {
        accessToken: string;
        refreshToken?: string;
        expiresAt?: Date;
        onRefresh?: TokenRefreshCallback;
    }) {
        this.accessToken = config.accessToken;
        this.refreshToken = config.refreshToken;
        this.expiresAt = config.expiresAt;
        this.onRefresh = config.onRefresh;
    }

    async getHeaders(): Promise<Record<string, string>> {
        if (this.needsRefresh()) {
            await this.refresh();
        }
        return { Authorization: `Bearer ${this.accessToken}` };
    }

    needsRefresh(): boolean {
        if (!this.expiresAt) return false;
        // Refresh 5 minutes before expiration
        return Date.now() + 300000 >= this.expiresAt.getTime();
    }

    async refresh(): Promise<void> {
        if (this.refreshPromise) return this.refreshPromise;
        if (!this.refreshToken) {
            throw new AuthenticationError({
                code: ErrorCode.AUTH_INVALID_KEY,
                message: 'No refresh token available',
                requestId: 'local',
            });
        }

        // Capture token locally to satisfy TS
        const refreshToken = this.refreshToken;
        this.refreshPromise = this.performRefresh(refreshToken);
        try {
            await this.refreshPromise;
        } finally {
            this.refreshPromise = undefined;
        }
    }

    private async performRefresh(token: string): Promise<void> {
        const response = await fetch('https://api.eclipse.dev/v1/auth/refresh', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: token }),
        });

        if (!response.ok) {
            throw new AuthenticationError({
                code: ErrorCode.AUTH_TOKEN_EXPIRED,
                message: `Refresh failed: ${response.status}`,
                requestId: response.headers.get('x-request-id') ?? 'unknown',
            });
        }

        interface RefreshResponse {
            access_token: string;
            refresh_token: string;
            expires_in: number;
        }

        const data = (await response.json()) as RefreshResponse;
        this.accessToken = data.access_token;
        this.refreshToken = data.refresh_token;
        this.expiresAt = new Date(Date.now() + data.expires_in * 1000);

        if (this.onRefresh) {
            await this.onRefresh({
                accessToken: this.accessToken,
                refreshToken: this.refreshToken,
                expiresAt: this.expiresAt,
            });
        }
    }
}
