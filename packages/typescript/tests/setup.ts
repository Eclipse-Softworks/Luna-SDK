import { beforeAll, afterAll, afterEach } from 'vitest';
import { server, resetHandlers } from './mocks/server';

// Start server before all tests
beforeAll(() => {
    server.listen({ onUnhandledRequest: 'error' });
});

// Reset handlers after each test (important for test isolation)
afterEach(() => {
    resetHandlers();
});

// Close server after all tests
afterAll(() => {
    server.close();
});
