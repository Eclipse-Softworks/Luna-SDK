import { describe, it, expect, vi, beforeEach } from 'vitest';
import { LunaClient } from '../src/client.js';
import { LunaError, AuthenticationError, NotFoundError } from '../src/errors/index.js';

describe('LunaClient', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    describe('constructor', () => {
        it('should create client with API key', () => {
            const client = new LunaClient({
                apiKey: 'lk_test_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa',
            });

            expect(client).toBeInstanceOf(LunaClient);
            expect(client.version).toBe('1.0.0');
        });

        it('should create client with access token', () => {
            const client = new LunaClient({
                accessToken: 'test-access-token',
            });

            expect(client).toBeInstanceOf(LunaClient);
        });

        it('should throw error when no auth provided', () => {
            expect(() => new LunaClient({})).toThrow('Either apiKey or accessToken must be provided');
        });

        it('should use custom base URL', () => {
            const client = new LunaClient({
                apiKey: 'lk_test_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa',
                baseUrl: 'https://api.staging.eclipse.dev',
            });

            expect(client).toBeInstanceOf(LunaClient);
        });
    });

    describe('resources', () => {
        it('should expose users resource', () => {
            const client = new LunaClient({
                apiKey: 'lk_test_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa',
            });

            expect(client.users).toBeDefined();
            expect(typeof client.users.list).toBe('function');
            expect(typeof client.users.get).toBe('function');
            expect(typeof client.users.create).toBe('function');
            expect(typeof client.users.update).toBe('function');
            expect(typeof client.users.delete).toBe('function');
        });

        it('should expose projects resource', () => {
            const client = new LunaClient({
                apiKey: 'lk_test_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa',
            });

            expect(client.projects).toBeDefined();
            expect(typeof client.projects.list).toBe('function');
            expect(typeof client.projects.get).toBe('function');
            expect(typeof client.projects.create).toBe('function');
            expect(typeof client.projects.update).toBe('function');
            expect(typeof client.projects.delete).toBe('function');
        });
    });
});

describe('LunaError', () => {
    it('should create error with all properties', () => {
        const error = new LunaError({
            code: 'LUNA_ERR_TEST',
            message: 'Test error',
            status: 500,
            requestId: 'req_123',
            details: { key: 'value' },
        });

        expect(error.code).toBe('LUNA_ERR_TEST');
        expect(error.message).toBe('Test error');
        expect(error.status).toBe(500);
        expect(error.requestId).toBe('req_123');
        expect(error.details).toEqual({ key: 'value' });
        expect(error.docsUrl).toBe('https://docs.eclipse.dev/luna/errors#LUNA_ERR_TEST');
    });

    it('should serialize to JSON', () => {
        const error = new LunaError({
            code: 'LUNA_ERR_TEST',
            message: 'Test error',
            status: 500,
            requestId: 'req_123',
        });

        const json = error.toJSON();
        expect(json.code).toBe('LUNA_ERR_TEST');
        expect(json.message).toBe('Test error');
        expect(json.request_id).toBe('req_123');
    });

    it('should identify retryable errors', () => {
        const serverError = new LunaError({
            code: 'LUNA_ERR_SERVER_INTERNAL',
            message: 'Internal error',
            status: 500,
            requestId: 'req_123',
        });

        const authError = new AuthenticationError({
            code: 'LUNA_ERR_AUTH_INVALID_KEY',
            message: 'Invalid key',
            requestId: 'req_123',
        });

        expect(serverError.retryable).toBe(true);
        expect(authError.retryable).toBe(false);
    });
});

describe('Error classes', () => {
    it('should create AuthenticationError with 401 status', () => {
        const error = new AuthenticationError({
            code: 'LUNA_ERR_AUTH_INVALID_KEY',
            message: 'Invalid API key',
            requestId: 'req_123',
        });

        expect(error.status).toBe(401);
        expect(error.name).toBe('AuthenticationError');
    });

    it('should create NotFoundError with 404 status', () => {
        const error = new NotFoundError({
            code: 'LUNA_ERR_RESOURCE_NOT_FOUND',
            message: 'User not found',
            requestId: 'req_123',
        });

        expect(error.status).toBe(404);
        expect(error.name).toBe('NotFoundError');
    });
});
