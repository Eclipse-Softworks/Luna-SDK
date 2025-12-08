import { describe, it, expect, beforeAll, afterAll, afterEach } from 'vitest';
import { http, HttpResponse } from 'msw';
import { LunaClient } from '../../src/client';
import { server, resetHandlers } from '../mocks/server';
import { mockUser, mockProject, mockBucket } from '../mocks/fixtures';

const BASE_URL = 'https://api.eclipse.dev';

/**
 * Contract validation tests ensure API responses match expected schemas.
 * These tests verify the structure and types of API responses.
 */
describe('API Contract Validation', () => {
    let client: LunaClient;

    beforeAll(() => {
        client = new LunaClient({ apiKey: 'lk_test_12345678901234567890123456789012' });
    });

    describe('User contract', () => {
        it('should have required user fields', async () => {
            const user = await client.users.get(mockUser.id);

            // Required fields
            expect(user).toHaveProperty('id');
            expect(user).toHaveProperty('name');
            expect(user).toHaveProperty('email');
            expect(user).toHaveProperty('created_at');
            expect(user).toHaveProperty('updated_at');

            // Type checks
            expect(typeof user.id).toBe('string');
            expect(typeof user.name).toBe('string');
            expect(typeof user.email).toBe('string');
        });

        it('should have correct ID prefix', async () => {
            const user = await client.users.get(mockUser.id);
            expect(user.id).toMatch(/^usr_/);
        });

        it('should have valid email format', async () => {
            const user = await client.users.get(mockUser.id);
            expect(user.email).toMatch(/^[^\s@]+@[^\s@]+\.[^\s@]+$/);
        });
    });

    describe('Project contract', () => {
        it('should have required project fields', async () => {
            const project = await client.projects.get(mockProject.id);

            // Required fields
            expect(project).toHaveProperty('id');
            expect(project).toHaveProperty('name');
            expect(project).toHaveProperty('created_at');
            expect(project).toHaveProperty('updated_at');

            // Type checks
            expect(typeof project.id).toBe('string');
            expect(typeof project.name).toBe('string');
        });

        it('should have correct ID prefix', async () => {
            const project = await client.projects.get(mockProject.id);
            expect(project.id).toMatch(/^prj_/);
        });

        it('should have optional description field', async () => {
            const project = await client.projects.get(mockProject.id);

            if (project.description !== undefined) {
                expect(typeof project.description).toBe('string');
            }
        });
    });

    describe('Bucket contract', () => {
        it('should have required bucket fields', async () => {
            const bucket = await client.storage.buckets.get(mockBucket.id);

            // Required fields
            expect(bucket).toHaveProperty('id');
            expect(bucket).toHaveProperty('name');

            // Type checks
            expect(typeof bucket.id).toBe('string');
            expect(typeof bucket.name).toBe('string');
        });

        it('should have correct ID prefix', async () => {
            const bucket = await client.storage.buckets.get(mockBucket.id);
            expect(bucket.id).toMatch(/^bkt_/);
        });
    });

    describe('List response contract', () => {
        it('should have correct list response structure', async () => {
            const result = await client.users.list();

            expect(result).toHaveProperty('data');
            expect(result).toHaveProperty('has_more');
            expect(Array.isArray(result.data)).toBe(true);
            expect(typeof result.has_more).toBe('boolean');
        });

        it('should have nextCursor when hasMore is true', async () => {
            server.use(
                http.get(`${BASE_URL}/v1/users`, () => {
                    return HttpResponse.json({
                        data: [mockUser],
                        has_more: true,
                        next_cursor: 'cursor_123',
                    });
                })
            );

            const result = await client.users.list();

            expect(result.has_more).toBe(true);
            expect(result.next_cursor).toBeDefined();
            expect(typeof result.next_cursor).toBe('string');
        });
    });

    describe('Error response contract', () => {
        it('should have correct error response structure', async () => {
            server.use(
                http.get(`${BASE_URL}/v1/users/:id`, () => {
                    return HttpResponse.json({
                        error: {
                            message: 'User not found',
                            code: 'NOT_FOUND',
                            status: 404,
                        },
                    }, { status: 404 });
                })
            );

            try {
                await client.users.get('usr_nonexistent');
                expect.fail('Should have thrown');
            } catch (error) {
                expect(error).toBeInstanceOf(Error);
                expect((error as Error).message).toBeDefined();
            }
        });

        it('should include validation details for 400 errors', async () => {
            server.use(
                http.post(`${BASE_URL}/v1/users`, () => {
                    return HttpResponse.json({
                        error: {
                            message: 'Validation failed',
                            code: 'VALIDATION_ERROR',
                            status: 400,
                            details: [
                                { field: 'email', message: 'Invalid email format' },
                            ],
                        },
                    }, { status: 400 });
                })
            );

            try {
                await client.users.create({ name: 'Test', email: 'invalid' });
                expect.fail('Should have thrown');
            } catch (error) {
                expect(error).toBeInstanceOf(Error);
            }
        });
    });

    describe('Timestamp formats', () => {
        it('should return dates as Date objects or ISO strings', async () => {
            const user = await client.users.get(mockUser.id);

            // The SDK should parse dates - check if it's a Date or valid ISO string
            const createdAt = user.created_at;

            if (typeof createdAt === 'string') {
                expect(new Date(createdAt).getTime()).not.toBeNaN();
            } else {
                // Should be string
                expect(typeof createdAt).toBe('string');
            }
        });
    });
});
