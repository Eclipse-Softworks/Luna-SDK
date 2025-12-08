import type { RetryConfig } from './types.js';
import type { LunaError } from '../errors/base.js';

/**
 * Default retry configuration
 */
export const DEFAULT_RETRY_CONFIG: RetryConfig = {
    maxRetries: 3,
    initialDelayMs: 500,
    maxDelayMs: 30_000,
    backoffMultiplier: 2,
    retryableStatuses: [408, 429, 500, 502, 503, 504],
};

/**
 * Handles retry logic with exponential backoff
 */
export class RetryHandler {
    private readonly config: RetryConfig;

    constructor(config: Partial<RetryConfig> = {}) {
        this.config = { ...DEFAULT_RETRY_CONFIG, ...config };
    }

    /**
     * Determine if an error/status is retryable
     */
    isRetryable(status: number, error?: LunaError): boolean {
        if (error?.retryable === false) {
            return false;
        }
        return this.config.retryableStatuses.includes(status);
    }

    /**
     * Calculate delay for next retry attempt
     */
    getDelayMs(attempt: number, retryAfterHeader?: string): number {
        // Respect Retry-After header if present
        if (retryAfterHeader) {
            const retryAfter = parseInt(retryAfterHeader, 10);
            if (!isNaN(retryAfter)) {
                return retryAfter * 1000;
            }
        }

        // Exponential backoff with jitter
        const delay = Math.min(
            this.config.initialDelayMs * Math.pow(this.config.backoffMultiplier, attempt),
            this.config.maxDelayMs
        );

        // Add random jitter (±10%)
        const jitter = delay * 0.1 * (Math.random() * 2 - 1);

        return Math.round(delay + jitter);
    }

    /**
     * Check if more retries are allowed
     */
    shouldRetry(attempt: number): boolean {
        return attempt < this.config.maxRetries;
    }

    /**
     * Wait for the calculated delay
     */
    async wait(attempt: number, retryAfterHeader?: string): Promise<void> {
        const delayMs = this.getDelayMs(attempt, retryAfterHeader);
        await new Promise((resolve) => setTimeout(resolve, delayMs));
    }
}
