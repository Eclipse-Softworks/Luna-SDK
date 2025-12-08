import { http, HttpResponse } from 'msw';
import {
    mockUser,
    mockUsers,
    mockProject,
    mockProjects,
    mockBucket,
    mockBuckets,
    mockFile,
    mockListResponse,
    mockErrorResponse,
    mockAuthErrorResponse,
} from './fixtures';

const BASE_URL = 'https://api.eclipse.dev';

// Helper to check auth header
const requireAuth = (request: Request) => {
    const authHeader = request.headers.get('Authorization');
    console.log('REQ HEADERS:', Object.fromEntries(request.headers.entries()));
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
        return HttpResponse.json(mockAuthErrorResponse, { status: 401 });
    }
    return null;
};

export const handlers = [
    // ==================== USERS ====================

    // List users
    http.get(`${BASE_URL}/v1/users`, ({ request }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        const url = new URL(request.url);
        const limit = parseInt(url.searchParams.get('limit') || '10');
        const cursor = url.searchParams.get('cursor');

        const users = cursor ? mockUsers.slice(1) : mockUsers.slice(0, limit);
        return HttpResponse.json(mockListResponse(users, mockUsers.length > limit));
    }),

    // Get user by ID
    http.get(`${BASE_URL}/v1/users/:id`, ({ request, params }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        const { id } = params;
        const user = mockUsers.find(u => u.id === id);

        if (!user) {
            return HttpResponse.json(mockErrorResponse, { status: 404 });
        }

        return HttpResponse.json(user);
    }),

    // Create user
    http.post(`${BASE_URL}/v1/users`, async ({ request }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        const body = await request.json() as Record<string, unknown>;
        const newUser = {
            ...mockUser,
            id: `usr_${Date.now()}`,
            name: body.name as string,
            email: body.email as string,
            createdAt: new Date(),
            updatedAt: new Date(),
        };

        return HttpResponse.json(newUser, { status: 201 });
    }),

    // Update user
    http.patch(`${BASE_URL}/v1/users/:id`, async ({ request, params }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        const { id } = params;
        const user = mockUsers.find(u => u.id === id);

        if (!user) {
            return HttpResponse.json(mockErrorResponse, { status: 404 });
        }

        const body = await request.json() as Record<string, unknown>;
        return HttpResponse.json({ ...user, ...body, updatedAt: new Date() });
    }),

    // Delete user
    http.delete(`${BASE_URL}/v1/users/:id`, ({ request, params }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        const { id } = params;
        const user = mockUsers.find(u => u.id === id);

        if (!user) {
            return HttpResponse.json(mockErrorResponse, { status: 404 });
        }

        return new HttpResponse(null, { status: 204 });
    }),

    // ==================== PROJECTS ====================

    // List projects
    http.get(`${BASE_URL}/v1/projects`, ({ request }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        return HttpResponse.json(mockListResponse(mockProjects));
    }),

    // Get project by ID
    http.get(`${BASE_URL}/v1/projects/:id`, ({ request, params }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        const { id } = params;
        const project = mockProjects.find(p => p.id === id);

        if (!project) {
            return HttpResponse.json(mockErrorResponse, { status: 404 });
        }

        return HttpResponse.json(project);
    }),

    // Create project
    http.post(`${BASE_URL}/v1/projects`, async ({ request }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        const body = await request.json() as Record<string, unknown>;
        const newProject = {
            ...mockProject,
            id: `prj_${Date.now()}`,
            name: body.name as string,
            description: body.description as string | undefined,
            createdAt: new Date(),
            updatedAt: new Date(),
        };

        return HttpResponse.json(newProject, { status: 201 });
    }),

    // Delete project
    http.delete(`${BASE_URL}/v1/projects/:id`, ({ request, params }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        const { id } = params;
        const project = mockProjects.find(p => p.id === id);

        if (!project) {
            return HttpResponse.json(mockErrorResponse, { status: 404 });
        }

        return new HttpResponse(null, { status: 204 });
    }),

    // ==================== STORAGE ====================

    // List buckets
    http.get(`${BASE_URL}/v1/storage/buckets`, ({ request }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        return HttpResponse.json(mockListResponse(mockBuckets));
    }),

    // Get bucket
    http.get(`${BASE_URL}/v1/storage/buckets/:id`, ({ request, params }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        const { id } = params;
        const bucket = mockBuckets.find(b => b.id === id);

        if (!bucket) {
            return HttpResponse.json(mockErrorResponse, { status: 404 });
        }

        return HttpResponse.json(bucket);
    }),

    // List files in bucket
    http.get(`${BASE_URL}/v1/storage/buckets/:bucketId/files`, ({ request }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        return HttpResponse.json(mockListResponse([mockFile]));
    }),

    // Upload file
    http.post(`${BASE_URL}/v1/storage/buckets/:bucketId/upload`, ({ request }) => {
        const authError = requireAuth(request);
        if (authError) return authError;

        return HttpResponse.json({
            ...mockFile,
            id: `file_${Date.now()}`,
            createdAt: new Date(),
        }, { status: 201 });
    }),

    // ==================== HEALTH CHECK ====================

    http.get(`${BASE_URL}/health`, () => {
        return HttpResponse.json({ status: 'healthy', version: '1.0.0' });
    }),
];

// Error simulation handlers (for testing error scenarios)
export const errorHandlers = {
    rateLimited: http.get(`${BASE_URL}/v1/users`, () => {
        return HttpResponse.json(
            { error: { message: 'Rate limit exceeded', code: 'RATE_LIMIT_EXCEEDED', status: 429 } },
            { status: 429, headers: { 'Retry-After': '60' } }
        );
    }),

    serverError: http.get(`${BASE_URL}/v1/users`, () => {
        return HttpResponse.json(
            { error: { message: 'Internal server error', code: 'SERVER_ERROR', status: 500 } },
            { status: 500 }
        );
    }),

    networkError: http.get(`${BASE_URL}/v1/users`, () => {
        return HttpResponse.error();
    }),
};
