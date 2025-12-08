import { setupServer } from 'msw/node';
import { handlers } from './handlers';

// Create and export the mock server
export const server = setupServer(...handlers);

// Server lifecycle helpers
export const startServer = () => {
    server.listen({ onUnhandledRequest: 'error' });
};

export const stopServer = () => {
    server.close();
};

export const resetHandlers = () => {
    server.resetHandlers();
};

// Helper to add custom handlers for specific tests
export const useHandlers = (...customHandlers: Parameters<typeof server.use>) => {
    server.use(...customHandlers);
};
