import type { AuthProvider } from '../auth/types.js';
import type { Logger } from '../telemetry/index.js';
import type { TelemetryConfig } from '../telemetry/types.js';

/**
 * HTTP request configuration
 */
export interface RequestConfig {
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
    path: string;
    headers?: Record<string, string>;
    query?: Record<string, string | string[] | undefined>;
    body?: unknown;
    timeout?: number;
}

/**
 * HTTP response wrapper
 */
export interface Response<T = unknown> {
    data: T;
    status: number;
    headers: Record<string, string>;
    requestId: string;
}

/**
 * HTTP client configuration
 */
export interface HttpClientConfig {
    baseUrl: string;
    timeout: number;
    maxRetries: number;
    authProvider: AuthProvider;
    logger: Logger;
    telemetry?: TelemetryConfig;
}

/**
 * Retry configuration
 */
export interface RetryConfig {
    maxRetries: number;
    initialDelayMs: number;
    maxDelayMs: number;
    backoffMultiplier: number;
    retryableStatuses: number[];
}
