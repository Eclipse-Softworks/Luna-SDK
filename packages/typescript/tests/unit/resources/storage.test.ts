import { describe, it, expect, beforeAll, afterAll, afterEach } from 'vitest';
import { LunaClient } from '../../../src/client';
import { server, resetHandlers } from '../../mocks/server';
import { mockBucket, mockBuckets, mockFile } from '../../mocks/fixtures';

describe('StorageResource', () => {
    let client: LunaClient;

    beforeAll(() => {
        client = new LunaClient({ apiKey: 'lk_test_12345678901234567890123456789012' });
    });

    describe('buckets', () => {
        describe('list()', () => {
            it('should return a list of buckets', async () => {
                const result = await client.storage.buckets.list();

                expect(result.data).toHaveLength(mockBuckets.length);
                expect(result.data[0].id).toBe(mockBuckets[0].id);
                expect(result.data[0].name).toBe(mockBuckets[0].name);
            });
        });

        describe('get()', () => {
            it('should return a single bucket by ID', async () => {
                const result = await client.storage.buckets.get(mockBucket.id);

                expect(result.id).toBe(mockBucket.id);
                expect(result.name).toBe(mockBucket.name);
            });

            it('should throw NotFoundError for non-existent bucket', async () => {
                await expect(client.storage.buckets.get('bkt_nonexistent')).rejects.toThrow();
            });
        });
    });

    describe('files', () => {
        describe('list()', () => {
            it('should return files in a bucket', async () => {
                const result = await client.storage.files.list(mockBucket.id);

                expect(result.data).toBeDefined();
                expect(Array.isArray(result.data)).toBe(true);
            });
        });

        describe('upload()', () => {
            it('should upload a file to a bucket', async () => {
                const fileContent = Buffer.from('test file content');
                const result = await client.storage.buckets.upload(mockBucket.id, fileContent, 'test.txt');

                expect(result.id).toBeDefined();
                expect(result.id).toMatch(/^file_/);
            });
        });
    });
});
