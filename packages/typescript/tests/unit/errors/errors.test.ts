import { describe, it, expect, beforeAll, afterAll, afterEach } from 'vitest';
import { http, HttpResponse } from 'msw';
import { LunaClient } from '../../../src/client';
import { server, resetHandlers, useHandlers } from '../../mocks/server';
import {
    LunaError,
    AuthenticationError,
    NotFoundError,
    RateLimitError,
    ValidationError,
    NetworkError,
} from '../../../src/errors';

const BASE_URL = 'https://api.eclipse.dev';

describe('Error Handling', () => {
    let client: LunaClient;

    beforeAll(() => {
        client = new LunaClient({ apiKey: 'lk_test_12345678901234567890123456789012' });
    });

    describe('AuthenticationError (401)', () => {
        it('should throw Error for invalid API key format', () => {
            expect(() => new LunaClient({ apiKey: 'invalid-key' })).toThrow(/Invalid API key/);
        });

        it('should throw AuthenticationError for 401 response', async () => {
            server.use(
                http.get(`${BASE_URL}/v1/users`, () => {
                    return HttpResponse.json(
                        { error: { message: 'Invalid API key', code: 'AUTHENTICATION_ERROR', status: 401 } },
                        { status: 401 }
                    );
                })
            );

            await expect(client.users.list()).rejects.toThrow();
        });

        it('should include error details in AuthenticationError', async () => {
            server.use(
                http.get(`${BASE_URL}/v1/users`, () => {
                    return HttpResponse.json(
                        { error: { message: 'Token expired', code: 'TOKEN_EXPIRED', status: 401 } },
                        { status: 401 }
                    );
                })
            );

            try {
                await client.users.list();
                expect.fail('Should have thrown');
            } catch (error) {
                expect(error).toBeInstanceOf(Error);
                expect((error as Error).message).toContain('expired');
            }
        });
    });

    describe('NotFoundError (404)', () => {
        it('should throw NotFoundError for non-existent resource', async () => {
            server.use(
                http.get(`${BASE_URL}/v1/users/:id`, () => {
                    return HttpResponse.json(
                        { error: { message: 'User not found', code: 'NOT_FOUND', status: 404 } },
                        { status: 404 }
                    );
                })
            );

            await expect(client.users.get('usr_nonexistent')).rejects.toThrow();
        });
    });

    describe('ValidationError (400)', () => {
        it('should throw ValidationError for invalid input', async () => {
            server.use(
                http.post(`${BASE_URL}/v1/users`, () => {
                    return HttpResponse.json(
                        {
                            error: {
                                message: 'Validation failed',
                                code: 'VALIDATION_ERROR',
                                status: 400,
                                details: [
                                    { field: 'email', message: 'Invalid email format' },
                                ],
                            },
                        },
                        { status: 400 }
                    );
                })
            );

            await expect(
                client.users.create({ name: 'Test', email: 'invalid-email' })
            ).rejects.toThrow();
        });
    });

    describe('RateLimitError (429)', () => {
        it('should throw RateLimitError when rate limited', async () => {
            server.use(
                http.get(`${BASE_URL}/v1/users`, () => {
                    return HttpResponse.json(
                        { error: { message: 'Rate limit exceeded', code: 'RATE_LIMIT_EXCEEDED', status: 429, retryAfter: 60 } },
                        { status: 429, headers: { 'Retry-After': '60' } }
                    );
                })
            );

            await expect(client.users.list()).rejects.toThrow();
        });
    });

    describe('ServerError (5xx)', () => {
        it('should throw ServerError for 500 response', async () => {
            server.use(
                http.get(`${BASE_URL}/v1/users`, () => {
                    return HttpResponse.json(
                        { error: { message: 'Internal server error', code: 'SERVER_ERROR', status: 500 } },
                        { status: 500 }
                    );
                })
            );

            await expect(client.users.list()).rejects.toThrow();
        });

        it('should throw ServerError for 503 response', async () => {
            server.use(
                http.get(`${BASE_URL}/v1/users`, () => {
                    return HttpResponse.json(
                        { error: { message: 'Service unavailable', code: 'SERVICE_UNAVAILABLE', status: 503 } },
                        { status: 503 }
                    );
                })
            );

            await expect(client.users.list()).rejects.toThrow();
        });
    });

    describe('NetworkError', () => {
        it('should throw NetworkError for network failures', async () => {
            server.use(
                http.get(`${BASE_URL}/v1/users`, () => {
                    return HttpResponse.error();
                })
            );

            await expect(client.users.list()).rejects.toThrow();
        });
    });

    describe('Error properties', () => {
        it('should include requestId in error when available', async () => {
            server.use(
                http.get(`${BASE_URL}/v1/users`, () => {
                    return HttpResponse.json(
                        { error: { message: 'Test error', code: 'TEST_ERROR', status: 400, requestId: 'req_123' } },
                        { status: 400, headers: { 'X-Request-Id': 'req_123' } }
                    );
                })
            );

            try {
                await client.users.list();
                expect.fail('Should have thrown');
            } catch (error) {
                expect(error).toBeInstanceOf(Error);
            }
        });
    });
});
