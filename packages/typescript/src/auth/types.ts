/**
 * Token pair returned from authentication
 */
export interface TokenPair {
    accessToken: string;
    refreshToken: string;
    expiresAt?: Date;
}

/**
 * Callback invoked when tokens are refreshed
 */
export type TokenRefreshCallback = (tokens: TokenPair) => void | Promise<void>;

/**
 * Interface for authentication providers
 */
export interface AuthProvider {
    /**
     * Get authorization headers for a request
     */
    getHeaders(): Promise<Record<string, string>>;

    /**
     * Check if credentials need refresh
     */
    needsRefresh(): boolean;

    /**
     * Refresh credentials if needed
     */
    refresh?(): Promise<void>;
}
