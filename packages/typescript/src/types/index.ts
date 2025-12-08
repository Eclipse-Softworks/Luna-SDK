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

/**
 * Residence resource
 */
export interface Residence {
    id: string;
    name: string;
    slug: string;
    address: string;
    description?: string;

    // Filters & Attributes
    is_nsfas_accredited: boolean;
    min_price: number;
    max_price: number;
    currency_code: string;
    gender_policy: 'mixed' | 'male' | 'female';

    // Location & Relations
    location: {
        latitude: number;
        longitude: number;
        suburb?: string;
        city?: string;
    };
    campus_ids: string[];

    // Social
    rating: number;
    review_count: number;

    images: string[];
    amenities: string[];
}

/**
 * Search/Filter parameters for residences
 */
export interface ResidenceSearch extends PaginationParams {
    query?: string;
    nsfas?: boolean;
    min_price?: number;
    max_price?: number;
    gender?: 'male' | 'female' | 'mixed';
    campus_id?: string;
    radius?: number;
    min_rating?: number;
}

/**
 * Campus resource
 */
export interface Campus {
    id: string;
    name: string;
    location: {
        latitude: number;
        longitude: number;
    };
}

/**
 * List of campuses
 */
export type CampusList = ListResponse<Campus>;

/**
 * Paginated list of residences
 */
export type ResidenceList = ListResponse<Residence>;

/**
 * Group resource
 */
export interface Group {
    id: string;
    name: string;
    description?: string;
    permissions: string[];
    member_ids: string[];
}

/**
 * Paginated list of groups
 */
export type GroupList = ListResponse<Group>;

/**
 * Parameters for creating a group
 */
export interface GroupCreate {
    name: string;
    description?: string;
    permissions?: string[];
    member_ids?: string[];
}

/**
 * Storage Bucket
 */
export interface Bucket {
    id: string;
    name: string;
    region: string;
}

/**
 * List of buckets
 */
export type BucketList = ListResponse<Bucket>;

/**
 * File Object
 */
export interface FileObject {
    id: string;
    bucket_id: string;
    key: string;
    size: number;
    content_type: string;
    url: string;
}

/**
 * AI Completion Request
 */
export interface CompletionRequest {
    model: string;
    messages: { role: 'system' | 'user' | 'assistant'; content: string }[];
    temperature?: number;
}

/**
 * AI Completion Response
 */
export interface CompletionResponse {
    id: string;
    choices: { index: number; message: { role: string; content: string } }[];
}

/**
 * Automation Workflow
 */
export interface Workflow {
    id: string;
    name: string;
    trigger_type: 'event' | 'schedule' | 'webhook';
    is_active: boolean;
}

/**
 * Workflow List
 */
export type WorkflowList = ListResponse<Workflow>;

/**
 * Workflow Run
 */
export interface WorkflowRun {
    id: string;
    workflow_id: string;
    status: 'pending' | 'running' | 'completed' | 'failed';
    started_at: string;
}
