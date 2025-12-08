/**
 * Utility functions for Luna SDK
 */

/**
 * Sleep for a specified number of milliseconds
 */
export function sleep(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * Generate a unique request ID
 */
export function generateRequestId(): string {
    const timestamp = Date.now().toString(36);
    const randomPart = Math.random().toString(36).substring(2, 10);
    return `req_${timestamp}${randomPart}`;
}

/**
 * Deep clone an object
 */
export function deepClone<T>(obj: T): T {
    return JSON.parse(JSON.stringify(obj)) as T;
}

/**
 * Check if a value is a plain object
 */
export function isPlainObject(value: unknown): value is Record<string, unknown> {
    return typeof value === 'object' && value !== null && !Array.isArray(value);
}

/**
 * Merge objects deeply
 */
export function deepMerge<T extends Record<string, unknown>>(
    target: T,
    ...sources: Partial<T>[]
): T {
    const result = { ...target };

    for (const source of sources) {
        for (const key in source) {
            const sourceValue = source[key];
            const targetValue = result[key];

            if (isPlainObject(sourceValue) && isPlainObject(targetValue)) {
                (result as Record<string, unknown>)[key] = deepMerge(
                    targetValue as Record<string, unknown>,
                    sourceValue as Record<string, unknown>
                );
            } else if (sourceValue !== undefined) {
                (result as Record<string, unknown>)[key] = sourceValue;
            }
        }
    }

    return result;
}

/**
 * Omit keys from an object
 */
export function omit<T extends Record<string, unknown>, K extends keyof T>(
    obj: T,
    keys: K[]
): Omit<T, K> {
    const result = { ...obj };
    for (const key of keys) {
        delete result[key];
    }
    return result;
}

/**
 * Pick keys from an object
 */
export function pick<T extends Record<string, unknown>, K extends keyof T>(
    obj: T,
    keys: K[]
): Pick<T, K> {
    const result = {} as Pick<T, K>;
    for (const key of keys) {
        if (key in obj) {
            result[key] = obj[key];
        }
    }
    return result;
}

/**
 * Check if a string is a valid URL
 */
export function isValidUrl(str: string): boolean {
    try {
        new URL(str);
        return true;
    } catch {
        return false;
    }
}

/**
 * Build URL with query parameters
 */
export function buildUrl(
    base: string,
    path: string,
    query?: Record<string, string | string[] | undefined>
): string {
    const url = new URL(path.startsWith('/') ? path : `/${path}`, base);

    if (query) {
        for (const [key, value] of Object.entries(query)) {
            if (value === undefined) continue;
            if (Array.isArray(value)) {
                for (const v of value) {
                    url.searchParams.append(key, v);
                }
            } else {
                url.searchParams.set(key, value);
            }
        }
    }

    return url.toString();
}

/**
 * Retry a function with exponential backoff
 */
export async function retry<T>(
    fn: () => Promise<T>,
    options: {
        maxRetries?: number;
        initialDelayMs?: number;
        maxDelayMs?: number;
        backoffMultiplier?: number;
        shouldRetry?: (error: unknown) => boolean;
    } = {}
): Promise<T> {
    const {
        maxRetries = 3,
        initialDelayMs = 500,
        maxDelayMs = 30000,
        backoffMultiplier = 2,
        shouldRetry = (): boolean => true,
    } = options;

    let lastError: unknown;

    for (let attempt = 0; attempt <= maxRetries; attempt++) {
        try {
            return await fn();
        } catch (error) {
            lastError = error;

            if (attempt >= maxRetries || !shouldRetry(error)) {
                throw error;
            }

            const delay = Math.min(initialDelayMs * Math.pow(backoffMultiplier, attempt), maxDelayMs);
            const jitter = delay * 0.1 * (Math.random() * 2 - 1);
            await sleep(Math.round(delay + jitter));
        }
    }

    throw lastError;
}
