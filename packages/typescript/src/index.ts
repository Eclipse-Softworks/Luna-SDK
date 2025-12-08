// Luna SDK - TypeScript Client
// Eclipse Softworks Platform API

export { LunaClient } from './client.js';
export type { ClientConfig, ClientOptions } from './client.js';

// Auth
export { ApiKeyAuth, TokenAuth } from './auth/index.js';
export type { AuthProvider, TokenPair, TokenRefreshCallback } from './auth/index.js';

// Resources
export { UsersResource } from './resources/users.js';
export { ProjectsResource } from './resources/projects.js';

// Errors
export {
    LunaError,
    AuthenticationError,
    AuthorizationError,
    ValidationError,
    RateLimitError,
    NetworkError,
    NotFoundError,
    ConflictError,
    ServerError,
} from './errors/index.js';
export { ErrorCode } from './errors/codes.js';

// Types
export type {
    User,
    UserCreate,
    UserUpdate,
    UserList,
    Project,
    ProjectCreate,
    ProjectUpdate,
    ProjectList,
    PaginationParams,
    ListResponse,
} from './types/index.js';

// HTTP
export type { RequestConfig, Response } from './http/index.js';

// Telemetry
export type { Logger, LogLevel, TelemetryConfig } from './telemetry/index.js';
