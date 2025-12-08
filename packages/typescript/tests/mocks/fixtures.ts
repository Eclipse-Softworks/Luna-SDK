import { User, Project, Bucket, FileObject } from '../../src/types';

// User fixtures
export const mockUser: User = {
    id: 'usr_123456789',
    name: 'John Doe',
    email: 'john@example.com',
    avatar_url: 'https://example.com/avatar.jpg',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
};

export const mockUsers: User[] = [
    mockUser,
    {
        id: 'usr_987654321',
        name: 'Jane Smith',
        email: 'jane@example.com',
        avatar_url: undefined,
        created_at: '2024-01-02T00:00:00Z',
        updated_at: '2024-01-02T00:00:00Z',
    },
];

export const mockUserCreate = {
    name: 'New User',
    email: 'newuser@example.com',
};

// Project fixtures
export const mockProject: Project = {
    id: 'prj_123456789',
    name: 'Test Project',
    description: 'A test project for unit tests',
    owner_id: 'usr_123456789',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
};

export const mockProjects: Project[] = [
    mockProject,
    {
        id: 'prj_987654321',
        name: 'Another Project',
        description: 'Another test project',
        owner_id: 'usr_987654321',
        created_at: '2024-01-02T00:00:00Z',
        updated_at: '2024-01-02T00:00:00Z',
    },
];

export const mockProjectCreate = {
    name: 'New Project',
    description: 'A new project',
};

// Storage fixtures
export const mockBucket: Bucket = {
    id: 'bkt_123456789',
    name: 'test-bucket',
    region: 'us-east-1',
};

export const mockBuckets: Bucket[] = [
    mockBucket,
    {
        id: 'bkt_987654321',
        name: 'public-bucket',
        region: 'eu-west-1',
    },
];

export const mockFile: FileObject = {
    id: 'file_123456789',
    bucket_id: 'bkt_123456789',
    key: 'test-file.pdf',
    size: 1024,
    content_type: 'application/pdf',
    url: 'https://storage.eclipse.dev/bkt_123456789/test-file.pdf',
};

// API response fixtures
export const mockListResponse = <T>(data: T[], hasMore = false, nextCursor?: string) => ({
    data,
    has_more: hasMore,
    next_cursor: nextCursor,
});

// Error fixtures
export const mockErrorResponse = {
    error: {
        message: 'Resource not found',
        code: 'NOT_FOUND',
        status: 404,
    },
};

export const mockValidationErrorResponse = {
    error: {
        message: 'Validation failed',
        code: 'VALIDATION_ERROR',
        status: 400,
        details: [
            { field: 'email', message: 'Invalid email format' },
        ],
    },
};

export const mockRateLimitErrorResponse = {
    error: {
        message: 'Rate limit exceeded',
        code: 'RATE_LIMIT_EXCEEDED',
        status: 429,
        retryAfter: 60,
    },
};

export const mockAuthErrorResponse = {
    error: {
        message: 'Invalid API key',
        code: 'AUTHENTICATION_ERROR',
        status: 401,
    },
};
