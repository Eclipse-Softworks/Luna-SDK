import { ErrorCode, type ErrorCodeType } from './codes.js';

/**
 * Parameters for creating a Luna error
 */
export interface LunaErrorParams {
    code: ErrorCodeType | string;
    message: string;
    status: number;
    requestId: string;
    details?: Record<string, unknown>;
}

/**
 * Base error class for all Luna SDK errors
 */
export class LunaError extends Error {
    /** Error code (e.g., LUNA_ERR_AUTH_INVALID_KEY) */
    readonly code: string;
    /** HTTP status code */
    readonly status: number;
    /** Unique request identifier for debugging */
    readonly requestId: string;
    /** Additional error details */
    readonly details?: Record<string, unknown>;

    constructor(params: LunaErrorParams) {
        super(params.message);
        this.name = 'LunaError';
        this.code = params.code;
        this.status = params.status;
        this.requestId = params.requestId;
        this.details = params.details;

        // Maintains proper stack trace in V8
        if (Error.captureStackTrace) {
            Error.captureStackTrace(this, this.constructor);
        }
    }

    /** URL to documentation for this error */
    get docsUrl(): string {
        return `https://docs.eclipse.dev/luna/errors#${this.code}`;
    }

    /** Check if this error is retryable */
    get retryable(): boolean {
        return [
            ErrorCode.RATE_LIMIT_EXCEEDED,
            ErrorCode.NETWORK_TIMEOUT,
            ErrorCode.NETWORK_CONNECTION,
            ErrorCode.SERVER_INTERNAL,
            ErrorCode.SERVER_UNAVAILABLE,
        ].includes(this.code as ErrorCodeType);
    }

    toJSON(): Record<string, unknown> {
        return {
            code: this.code,
            message: this.message,
            status: this.status,
            request_id: this.requestId,
            details: this.details,
            docs_url: this.docsUrl,
        };
    }
}

/**
 * Authentication failed (invalid API key, expired token)
 */
export class AuthenticationError extends LunaError {
    constructor(params: Omit<LunaErrorParams, 'status'>) {
        super({ ...params, status: 401 });
        this.name = 'AuthenticationError';
    }
}

/**
 * Authorization failed (insufficient permissions)
 */
export class AuthorizationError extends LunaError {
    constructor(params: Omit<LunaErrorParams, 'status'>) {
        super({ ...params, status: 403 });
        this.name = 'AuthorizationError';
    }
}

/**
 * Validation failed for request parameters
 */
export class ValidationError extends LunaError {
    constructor(params: Omit<LunaErrorParams, 'status'>) {
        super({ ...params, status: 400 });
        this.name = 'ValidationError';
    }
}

/**
 * Rate limit exceeded
 */
export class RateLimitError extends LunaError {
    /** Seconds until rate limit resets */
    readonly retryAfter?: number;

    constructor(params: Omit<LunaErrorParams, 'status'> & { retryAfter?: number }) {
        super({ ...params, status: 429 });
        this.name = 'RateLimitError';
        this.retryAfter = params.retryAfter;
    }
}

/**
 * Network-related errors (timeout, connection)
 */
export class NetworkError extends LunaError {
    constructor(params: Omit<LunaErrorParams, 'status'> & { status?: number }) {
        super({ ...params, status: params.status ?? 0 });
        this.name = 'NetworkError';
    }
}

/**
 * Resource not found
 */
export class NotFoundError extends LunaError {
    constructor(params: Omit<LunaErrorParams, 'status'>) {
        super({ ...params, status: 404 });
        this.name = 'NotFoundError';
    }
}

/**
 * Resource conflict (e.g., duplicate creation)
 */
export class ConflictError extends LunaError {
    constructor(params: Omit<LunaErrorParams, 'status'>) {
        super({ ...params, status: 409 });
        this.name = 'ConflictError';
    }
}

/**
 * Server-side errors
 */
export class ServerError extends LunaError {
    constructor(params: Omit<LunaErrorParams, 'status'> & { status?: number }) {
        super({ ...params, status: params.status ?? 500 });
        this.name = 'ServerError';
    }
}

/**
 * Create appropriate error class from API response
 */
export function createError(
    status: number,
    body: { code?: string; message?: string; details?: Record<string, unknown> },
    requestId: string,
    retryAfter?: number
): LunaError {
    const params = {
        code: body.code ?? 'LUNA_ERR_UNKNOWN',
        message: body.message ?? 'An unknown error occurred',
        requestId,
        details: body.details,
    };

    switch (status) {
        case 400:
            return new ValidationError(params);
        case 401:
            return new AuthenticationError(params);
        case 403:
            return new AuthorizationError(params);
        case 404:
            return new NotFoundError(params);
        case 409:
            return new ConflictError(params);
        case 429:
            return new RateLimitError({ ...params, retryAfter });
        default:
            if (status >= 500) {
                return new ServerError({ ...params, status });
            }
            return new LunaError({ ...params, status });
    }
}
