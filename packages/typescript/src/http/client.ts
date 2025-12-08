import type { RequestConfig, Response, HttpClientConfig } from './types.js';
import { RetryHandler } from './retry.js';
import { createError, NetworkError, LunaError } from '../errors/base.js';
import { ErrorCode } from '../errors/codes.js';

/**
 * HTTP client with authentication, retry, and telemetry
 */
export class HttpClient {
    private readonly baseUrl: string;
    private readonly timeout: number;
    private readonly authProvider: HttpClientConfig['authProvider'];
    private readonly logger: HttpClientConfig['logger'];
    private readonly retryHandler: RetryHandler;

    constructor(config: HttpClientConfig) {
        this.baseUrl = config.baseUrl.replace(/\/$/, ''); // Remove trailing slash
        this.timeout = config.timeout;
        this.authProvider = config.authProvider;
        this.logger = config.logger;
        this.retryHandler = new RetryHandler({ maxRetries: config.maxRetries });
    }

    /**
     * Make an HTTP request
     */
    async request<T>(config: RequestConfig): Promise<Response<T>> {
        const url = this.buildUrl(config.path, config.query);
        const requestId = this.generateRequestId();

        let lastError: LunaError | undefined;
        let attempt = 0;

        while (true) {
            try {
                const response = await this.executeRequest<T>(url, config, requestId);

                this.logger.info('HTTP request completed', {
                    request_id: requestId,
                    method: config.method,
                    path: config.path,
                    status: response.status,
                });

                return response;
            } catch (error) {
                lastError = error as LunaError;

                const status = lastError.status ?? 0;
                const shouldRetry =
                    this.retryHandler.shouldRetry(attempt) &&
                    this.retryHandler.isRetryable(status, lastError);

                if (!shouldRetry) {
                    this.logger.error('HTTP request failed', {
                        request_id: requestId,
                        method: config.method,
                        path: config.path,
                        error: lastError.code,
                        attempt,
                    });
                    throw lastError;
                }

                this.logger.warn('HTTP request failed, retrying', {
                    request_id: requestId,
                    method: config.method,
                    path: config.path,
                    status,
                    attempt,
                });

                await this.retryHandler.wait(
                    attempt,
                    (lastError as LunaError & { retryAfter?: number }).retryAfter?.toString()
                );
                attempt++;
            }
        }
    }

    /**
     * Execute a single HTTP request
     */
    private async executeRequest<T>(
        url: string,
        config: RequestConfig,
        requestId: string
    ): Promise<Response<T>> {
        const authHeaders = await this.authProvider.getHeaders();

        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
            Accept: 'application/json',
            'X-Request-Id': requestId,
            'User-Agent': `luna-sdk-typescript/1.0.0`,
            ...authHeaders,
            ...config.headers,
        };

        const controller = new AbortController();
        const timeoutId = setTimeout(
            () => controller.abort(),
            config.timeout ?? this.timeout
        );

        try {
            this.logger.debug('Sending HTTP request', {
                request_id: requestId,
                method: config.method,
                url,
            });

            const response = await fetch(url, {
                method: config.method,
                headers,
                body: config.body ? JSON.stringify(config.body) : undefined,
                signal: controller.signal,
            });

            clearTimeout(timeoutId);

            // Extract response headers
            const responseHeaders: Record<string, string> = {};
            response.headers.forEach((value, key) => {
                responseHeaders[key] = value;
            });

            // Get request ID from response (server may override)
            const serverRequestId =
                responseHeaders['x-request-id'] ?? requestId;

            // Parse response body
            const body = await this.parseResponseBody(response);

            // Handle error responses
            if (!response.ok) {
                const retryAfter = responseHeaders['retry-after']
                    ? parseInt(responseHeaders['retry-after'], 10)
                    : undefined;

                throw createError(
                    response.status,
                    body as { code?: string; message?: string; details?: Record<string, unknown> },
                    serverRequestId,
                    retryAfter
                );
            }

            return {
                data: body as T,
                status: response.status,
                headers: responseHeaders,
                requestId: serverRequestId,
            };
        } catch (error) {
            clearTimeout(timeoutId);

            // Handle abort/timeout
            if (error instanceof DOMException && error.name === 'AbortError') {
                throw new NetworkError({
                    code: ErrorCode.NETWORK_TIMEOUT,
                    message: 'Request timeout',
                    requestId,
                });
            }

            // Handle network errors
            if (error instanceof TypeError && error.message.includes('fetch')) {
                throw new NetworkError({
                    code: ErrorCode.NETWORK_CONNECTION,
                    message: 'Connection error',
                    requestId,
                });
            }

            // Re-throw Luna errors
            if (error instanceof LunaError) {
                throw error;
            }

            // Wrap unknown errors
            throw new NetworkError({
                code: ErrorCode.NETWORK_CONNECTION,
                message: (error as Error).message ?? 'Unknown error',
                requestId,
            });
        }
    }

    /**
     * Build full URL with query parameters
     */
    private buildUrl(
        path: string,
        query?: Record<string, string | string[] | undefined>
    ): string {
        const url = new URL(path.startsWith('/') ? path : `/${path}`, this.baseUrl);

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
     * Parse response body as JSON
     */
    private async parseResponseBody(response: globalThis.Response): Promise<unknown> {
        const contentType = response.headers.get('content-type');

        if (contentType?.includes('application/json')) {
            try {
                return await response.json();
            } catch {
                return null;
            }
        }

        return null;
    }

    /**
     * Generate unique request ID
     */
    private generateRequestId(): string {
        const timestamp = Date.now().toString(36);
        const random = Math.random().toString(36).substring(2, 10);
        return `req_${timestamp}${random}`;
    }
}
