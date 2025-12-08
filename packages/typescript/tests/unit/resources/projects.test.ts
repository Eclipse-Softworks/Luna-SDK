import { describe, it, expect, beforeAll, afterAll, afterEach } from 'vitest';
import { LunaClient } from '../../../src/client';
import { server, resetHandlers } from '../../mocks/server';
import { mockProject, mockProjects, mockProjectCreate } from '../../mocks/fixtures';

describe('ProjectsResource', () => {
    let client: LunaClient;

    beforeAll(() => {
        client = new LunaClient({ apiKey: 'lk_test_12345678901234567890123456789012' });
    });

    describe('list()', () => {
        it('should return a list of projects', async () => {
            const result = await client.projects.list();

            expect(result.data).toHaveLength(mockProjects.length);
            expect(result.data[0].id).toBe(mockProjects[0].id);
            expect(result.data[0].name).toBe(mockProjects[0].name);
        });

        it('should support pagination parameters', async () => {
            const result = await client.projects.list({ limit: 5 });

            expect(result.data).toBeDefined();
            expect(Array.isArray(result.data)).toBe(true);
        });
    });

    describe('get()', () => {
        it('should return a single project by ID', async () => {
            const result = await client.projects.get(mockProject.id);

            expect(result.id).toBe(mockProject.id);
            expect(result.name).toBe(mockProject.name);
            expect(result.description).toBe(mockProject.description);
        });

        it('should throw NotFoundError for non-existent project', async () => {
            await expect(client.projects.get('prj_nonexistent')).rejects.toThrow();
        });
    });

    describe('create()', () => {
        it('should create a new project', async () => {
            const result = await client.projects.create(mockProjectCreate);

            expect(result.name).toBe(mockProjectCreate.name);
            expect(result.description).toBe(mockProjectCreate.description);
            expect(result.id).toBeDefined();
            expect(result.id).toMatch(/^prj_/);
        });

        it('should create project with optional description', async () => {
            const result = await client.projects.create({ name: 'Minimal Project' });

            expect(result.name).toBe('Minimal Project');
            expect(result.id).toBeDefined();
        });
    });

    describe('delete()', () => {
        it('should delete an existing project', async () => {
            await expect(client.projects.delete(mockProject.id)).resolves.toBeUndefined();
        });

        it('should throw NotFoundError for non-existent project', async () => {
            await expect(client.projects.delete('prj_nonexistent')).rejects.toThrow();
        });
    });
});
