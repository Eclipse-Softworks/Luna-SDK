/**
 * Configuration loader for Luna SDK
 */

import type { ClientConfig } from '../client.js';

/**
 * Environment variable names
 */
const ENV_VARS = {
    API_KEY: 'LUNA_API_KEY',
    ACCESS_TOKEN: 'LUNA_ACCESS_TOKEN',
    REFRESH_TOKEN: 'LUNA_REFRESH_TOKEN',
    BASE_URL: 'LUNA_BASE_URL',
    TIMEOUT: 'LUNA_TIMEOUT',
    MAX_RETRIES: 'LUNA_MAX_RETRIES',
    LOG_LEVEL: 'LUNA_LOG_LEVEL',
} as const;

/**
 * Default configuration values
 */
export const DEFAULT_CONFIG = {
    baseUrl: 'https://api.eclipse.dev',
    timeout: 30000,
    maxRetries: 3,
    logLevel: 'info' as const,
} as const;

/**
 * Load configuration from environment variables
 */
export function loadFromEnv(): Partial<ClientConfig> {
    const config: Partial<ClientConfig> = {};

    // Auth
    const apiKey = getEnv(ENV_VARS.API_KEY);
    if (apiKey) {
        config.apiKey = apiKey;
    }

    const accessToken = getEnv(ENV_VARS.ACCESS_TOKEN);
    if (accessToken) {
        config.accessToken = accessToken;
    }

    const refreshToken = getEnv(ENV_VARS.REFRESH_TOKEN);
    if (refreshToken) {
        config.refreshToken = refreshToken;
    }

    // Base URL
    const baseUrl = getEnv(ENV_VARS.BASE_URL);
    if (baseUrl) {
        config.baseUrl = baseUrl;
    }

    // Timeout
    const timeout = getEnv(ENV_VARS.TIMEOUT);
    if (timeout) {
        const parsed = parseInt(timeout, 10);
        if (!isNaN(parsed) && parsed > 0) {
            config.timeout = parsed;
        }
    }

    // Max retries
    const maxRetries = getEnv(ENV_VARS.MAX_RETRIES);
    if (maxRetries) {
        const parsed = parseInt(maxRetries, 10);
        if (!isNaN(parsed) && parsed >= 0) {
            config.maxRetries = parsed;
        }
    }

    // Log level
    const logLevel = getEnv(ENV_VARS.LOG_LEVEL);
    if (logLevel && isValidLogLevel(logLevel)) {
        config.logLevel = logLevel;
    }

    return config;
}

/**
 * Merge configuration with defaults
 */
export function mergeConfig(userConfig: Partial<ClientConfig>): ClientConfig {
    const envConfig = loadFromEnv();

    return {
        ...DEFAULT_CONFIG,
        ...envConfig,
        ...userConfig,
    } as ClientConfig;
}

/**
 * Get environment variable value (works in Node.js and edge runtimes)
 */
function getEnv(name: string): string | undefined {
    if (typeof process !== 'undefined' && process.env) {
        return process.env[name];
    }
    return undefined;
}

/**
 * Check if a string is a valid log level
 */
function isValidLogLevel(value: string): value is 'error' | 'warn' | 'info' | 'debug' | 'trace' {
    return ['error', 'warn', 'info', 'debug', 'trace'].includes(value);
}

/**
 * Validate configuration
 */
export function validateConfig(config: ClientConfig): void {
    if (!config.apiKey && !config.accessToken) {
        throw new Error('Either apiKey or accessToken must be provided');
    }

    if (config.timeout !== undefined && (config.timeout <= 0 || !Number.isFinite(config.timeout))) {
        throw new Error('timeout must be a positive number');
    }

    if (
        config.maxRetries !== undefined &&
        (config.maxRetries < 0 || !Number.isInteger(config.maxRetries))
    ) {
        throw new Error('maxRetries must be a non-negative integer');
    }

    if (config.baseUrl) {
        try {
            new URL(config.baseUrl);
        } catch {
            throw new Error('baseUrl must be a valid URL');
        }
    }
}
