/**
 * Pagination parameters for list requests
 */
export interface PaginationParams {
    /** Maximum number of results (1-100, default: 20) */
    limit?: number;
    /** Cursor for pagination */
    cursor?: string;
}

/**
 * Generic list response with pagination
 */
export interface ListResponse<T> {
    data: T[];
    has_more: boolean;
    next_cursor?: string;
}

/**
 * User resource
 */
export interface User {
    id: string;
    email: string;
    name: string;
    avatar_url?: string;
    created_at: string;
    updated_at: string;
}

/**
 * Parameters for creating a user
 */
export interface UserCreate {
    email: string;
    name: string;
    avatar_url?: string;
}

/**
 * Parameters for updating a user
 */
export interface UserUpdate {
    name?: string;
    avatar_url?: string;
}

/**
 * Paginated list of users
 */
export type UserList = ListResponse<User>;

/**
 * Project resource
 */
export interface Project {
    id: string;
    name: string;
    description?: string;
    owner_id: string;
    created_at: string;
    updated_at: string;
}

/**
 * Parameters for creating a project
 */
export interface ProjectCreate {
    name: string;
    description?: string;
}

/**
 * Parameters for updating a project
 */
export interface ProjectUpdate {
    name?: string;
    description?: string;
}

/**
 * Paginated list of projects
 */
export type ProjectList = ListResponse<Project>;
