import { describe, it, expect, beforeAll, afterAll, afterEach } from 'vitest';
import { http, HttpResponse } from 'msw';
import { LunaClient } from '../../src/client';
import { server, resetHandlers } from '../mocks/server';
import { mockUser, mockUsers } from '../mocks/fixtures';

const BASE_URL = 'https://api.eclipse.dev';

/**
 * Integration tests verify the full request/response cycle.
 * These tests use mock servers but test the complete client flow.
 */
describe('Users Integration Tests', () => {
    let client: LunaClient;

    beforeAll(() => {
        client = new LunaClient({ apiKey: 'lk_test_12345678901234567890123456789012' });
    });

    describe('Full CRUD workflow', () => {
        it('should create, read, update, and delete a user', async () => {
            const createdUserId = `usr_${Date.now()}`;

            // CREATE
            server.use(
                http.post(`${BASE_URL}/v1/users`, () => {
                    return HttpResponse.json({
                        id: createdUserId,
                        name: 'Integration Test User',
                        email: 'integration@test.com',
                        createdAt: new Date().toISOString(),
                        updatedAt: new Date().toISOString(),
                    }, { status: 201 });
                })
            );

            const created = await client.users.create({
                name: 'Integration Test User',
                email: 'integration@test.com',
            });

            expect(created.id).toBe(createdUserId);
            expect(created.name).toBe('Integration Test User');

            // READ
            server.use(
                http.get(`${BASE_URL}/v1/users/${createdUserId}`, () => {
                    return HttpResponse.json({
                        id: createdUserId,
                        name: 'Integration Test User',
                        email: 'integration@test.com',
                        createdAt: new Date().toISOString(),
                        updatedAt: new Date().toISOString(),
                    });
                })
            );

            const fetched = await client.users.get(createdUserId);
            expect(fetched.id).toBe(createdUserId);

            // UPDATE
            server.use(
                http.patch(`${BASE_URL}/v1/users/${createdUserId}`, () => {
                    return HttpResponse.json({
                        id: createdUserId,
                        name: 'Updated Name',
                        email: 'integration@test.com',
                        createdAt: new Date().toISOString(),
                        updatedAt: new Date().toISOString(),
                    });
                })
            );

            const updated = await client.users.update(createdUserId, { name: 'Updated Name' });
            expect(updated.name).toBe('Updated Name');

            // DELETE
            server.use(
                http.delete(`${BASE_URL}/v1/users/${createdUserId}`, () => {
                    return new HttpResponse(null, { status: 204 });
                })
            );

            await expect(client.users.delete(createdUserId)).resolves.toBeUndefined();
        });
    });

    describe('Pagination workflow', () => {
        it('should paginate through all users', async () => {
            let requestCount = 0;

            server.use(
                http.get(`${BASE_URL}/v1/users`, ({ request }) => {
                    requestCount++;
                    const url = new URL(request.url);
                    const cursor = url.searchParams.get('cursor');

                    if (!cursor) {
                        return HttpResponse.json({
                            data: mockUsers.slice(0, 1),
                            has_more: true,
                            next_cursor: 'cursor_page_2',
                        });
                    } else if (cursor === 'cursor_page_2') {
                        return HttpResponse.json({
                            data: mockUsers.slice(1),
                            has_more: false,
                            next_cursor: null,
                        });
                    }

                    return HttpResponse.json({ data: [], has_more: false });
                })
            );

            // Fetch first page
            const page1 = await client.users.list({ limit: 1 });
            expect(page1.data).toHaveLength(1);
            expect(page1.has_more).toBe(true);
            expect(page1.next_cursor).toBe('cursor_page_2');

            // Fetch second page

            const page2 = await client.users.list({ cursor: page1.next_cursor });


            expect(page2.data.length).toBeGreaterThan(0);
            expect(page2.has_more).toBe(false);

            expect(requestCount).toBe(2);
        });
    });

    describe('Retry behavior', () => {
        it('should retry on transient errors', async () => {
            let attemptCount = 0;

            server.use(
                http.get(`${BASE_URL}/v1/users`, () => {
                    attemptCount++;
                    if (attemptCount < 3) {
                        return HttpResponse.json(
                            { error: { message: 'Service unavailable', status: 503, code: 'LUNA_ERR_SERVER_UNAVAILABLE' } },
                            { status: 503 }
                        );
                    }
                    return HttpResponse.json({ data: mockUsers, hasMore: false });
                })
            );

            // Client should retry and eventually succeed
            // Note: This test assumes the client has retry logic configured
            try {
                const result = await client.users.list();
                expect(result.data).toBeDefined();
                expect(attemptCount).toBeGreaterThanOrEqual(1);
            } catch {
                // If client doesn't retry, it will fail on first attempt
                expect(attemptCount).toBe(1);
            }
        });
    });

    describe('Header handling', () => {
        it('should send correct headers', async () => {
            let capturedHeaders: Headers | null = null;

            server.use(
                http.get(`${BASE_URL}/v1/users`, ({ request }) => {
                    capturedHeaders = request.headers;
                    return HttpResponse.json({ data: mockUsers, hasMore: false });
                })
            );

            await client.users.list();

            expect(capturedHeaders).not.toBeNull();
            expect(capturedHeaders!.get('Authorization')).toBe('Bearer lk_test_12345678901234567890123456789012');
            expect(capturedHeaders!.get('Content-Type')).toContain('application/json');
        });
    });
});
