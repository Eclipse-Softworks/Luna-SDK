import { describe, it, expect, beforeAll, afterAll, afterEach } from 'vitest';
import { LunaClient } from '../../../src/client';
import { server, resetHandlers } from '../../mocks/server';
import { mockUser, mockUsers, mockUserCreate } from '../../mocks/fixtures';

describe('UsersResource', () => {
    let client: LunaClient;

    beforeAll(() => {
        client = new LunaClient({ apiKey: 'lk_test_12345678901234567890123456789012' });
    });

    describe('list()', () => {
        it('should return a list of users', async () => {
            const result = await client.users.list();

            expect(result.data).toHaveLength(mockUsers.length);
            expect(result.data[0].id).toBe(mockUsers[0].id);
            expect(result.data[0].name).toBe(mockUsers[0].name);
            expect(result.data[0].email).toBe(mockUsers[0].email);
        });

        it('should support pagination with limit', async () => {
            const result = await client.users.list({ limit: 1 });

            expect(result.data.length).toBeLessThanOrEqual(1);
        });

        it('should support pagination with cursor', async () => {
            const result = await client.users.list({ cursor: 'next_page_cursor' });

            expect(result.data).toBeDefined();
            expect(Array.isArray(result.data)).toBe(true);
        });
    });

    describe('get()', () => {
        it('should return a single user by ID', async () => {
            const result = await client.users.get(mockUser.id);

            expect(result.id).toBe(mockUser.id);
            expect(result.name).toBe(mockUser.name);
            expect(result.email).toBe(mockUser.email);
        });

        it('should throw NotFoundError for non-existent user', async () => {
            await expect(client.users.get('usr_nonexistent')).rejects.toThrow();
        });
    });

    describe('create()', () => {
        it('should create a new user', async () => {
            const result = await client.users.create(mockUserCreate);

            expect(result.name).toBe(mockUserCreate.name);
            expect(result.email).toBe(mockUserCreate.email);
            expect(result.id).toBeDefined();
            expect(result.id).toMatch(/^usr_/);
        });

        it('should return user with timestamps', async () => {
            const result = await client.users.create(mockUserCreate);

            expect(result.created_at).toBeDefined();
            expect(result.updated_at).toBeDefined();
        });
    });

    describe('update()', () => {
        it('should update an existing user', async () => {
            const updates = { name: 'Updated Name' };
            const result = await client.users.update(mockUser.id, updates);

            expect(result.name).toBe(updates.name);
            expect(result.id).toBe(mockUser.id);
        });

        it('should throw NotFoundError for non-existent user', async () => {
            await expect(
                client.users.update('usr_nonexistent', { name: 'Test' })
            ).rejects.toThrow();
        });
    });

    describe('delete()', () => {
        it('should delete an existing user', async () => {
            // Should not throw
            await expect(client.users.delete(mockUser.id)).resolves.toBeUndefined();
        });

        it('should throw NotFoundError for non-existent user', async () => {
            await expect(client.users.delete('usr_nonexistent')).rejects.toThrow();
        });
    });
});
