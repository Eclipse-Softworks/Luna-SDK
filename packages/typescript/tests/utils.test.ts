import { describe, it, expect } from 'vitest';
import {
    sleep,
    generateRequestId,
    deepClone,
    isPlainObject,
    deepMerge,
    omit,
    pick,
    isValidUrl,
    buildUrl,
} from '../src/utils/index.js';

describe('utils', () => {
    describe('sleep', () => {
        it('should wait for specified duration', async () => {
            const start = Date.now();
            await sleep(50);
            const duration = Date.now() - start;

            expect(duration).toBeGreaterThanOrEqual(45);
            expect(duration).toBeLessThan(100);
        });
    });

    describe('generateRequestId', () => {
        it('should generate unique IDs', () => {
            const id1 = generateRequestId();
            const id2 = generateRequestId();

            expect(id1).toMatch(/^req_[a-z0-9]+$/);
            expect(id2).toMatch(/^req_[a-z0-9]+$/);
            expect(id1).not.toBe(id2);
        });
    });

    describe('deepClone', () => {
        it('should clone objects deeply', () => {
            const obj = { a: 1, b: { c: 2 } };
            const clone = deepClone(obj);

            expect(clone).toEqual(obj);
            expect(clone).not.toBe(obj);
            expect(clone.b).not.toBe(obj.b);
        });
    });

    describe('isPlainObject', () => {
        it('should return true for plain objects', () => {
            expect(isPlainObject({})).toBe(true);
            expect(isPlainObject({ a: 1 })).toBe(true);
        });

        it('should return false for non-objects', () => {
            expect(isPlainObject(null)).toBe(false);
            expect(isPlainObject(undefined)).toBe(false);
            expect(isPlainObject([])).toBe(false);
            expect(isPlainObject('string')).toBe(false);
            expect(isPlainObject(123)).toBe(false);
        });
    });

    describe('deepMerge', () => {
        it('should merge objects deeply', () => {
            const target = { a: 1, b: { c: 2 } };
            const source = { b: { d: 3 }, e: 4 };

            const result = deepMerge(target, source);

            expect(result).toEqual({ a: 1, b: { c: 2, d: 3 }, e: 4 });
        });
    });

    describe('omit', () => {
        it('should omit specified keys', () => {
            const obj = { a: 1, b: 2, c: 3 };
            const result = omit(obj, ['b']);

            expect(result).toEqual({ a: 1, c: 3 });
        });
    });

    describe('pick', () => {
        it('should pick specified keys', () => {
            const obj = { a: 1, b: 2, c: 3 };
            const result = pick(obj, ['a', 'c']);

            expect(result).toEqual({ a: 1, c: 3 });
        });
    });

    describe('isValidUrl', () => {
        it('should validate URLs', () => {
            expect(isValidUrl('https://example.com')).toBe(true);
            expect(isValidUrl('http://localhost:3000')).toBe(true);
            expect(isValidUrl('not-a-url')).toBe(false);
            expect(isValidUrl('')).toBe(false);
        });
    });

    describe('buildUrl', () => {
        it('should build URL with path', () => {
            const url = buildUrl('https://api.example.com', '/users');
            expect(url).toBe('https://api.example.com/users');
        });

        it('should build URL with query parameters', () => {
            const url = buildUrl('https://api.example.com', '/users', {
                limit: '10',
                offset: '0',
            });
            expect(url).toBe('https://api.example.com/users?limit=10&offset=0');
        });

        it('should handle array query parameters', () => {
            const url = buildUrl('https://api.example.com', '/users', {
                ids: ['1', '2', '3'],
            });
            expect(url).toContain('ids=1');
            expect(url).toContain('ids=2');
            expect(url).toContain('ids=3');
        });
    });
});
