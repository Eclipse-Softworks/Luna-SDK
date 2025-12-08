import type { HttpClient } from '../http/client.js';
import { Paginator } from './pagination.js';
import type {
    User,
    UserCreate,
    UserUpdate,
    UserList,
    PaginationParams,
} from '../types/index.js';

/**
 * Users resource for managing user accounts
 */
export class UsersResource {
    private readonly httpClient: HttpClient;
    private readonly basePath = '/v1/users';

    constructor(httpClient: HttpClient) {
        this.httpClient = httpClient;
    }

    /**
     * List all users with pagination
     */
    async list(params?: PaginationParams): Promise<UserList> {
        const response = await this.httpClient.request<UserList>({
            method: 'GET',
            path: this.basePath,
            query: {
                limit: params?.limit?.toString(),
                cursor: params?.cursor,
            },
        });

        return response.data;
    }

    /**
     * Iterate over all users automatically handling pagination
     */
    iterate(params?: PaginationParams): Paginator<User> {
        return new Paginator<User>(
            async (cursor) => {
                const page = await this.list({ ...params, cursor });
                return page;
            }
        );
    }

    /**
     * Get a user by ID
     */
    async get(userId: string): Promise<User> {
        this.validateUserId(userId);

        const response = await this.httpClient.request<User>({
            method: 'GET',
            path: `${this.basePath}/${userId}`,
        });

        return response.data;
    }

    /**
     * Create a new user
     */
    async create(data: UserCreate): Promise<User> {
        this.validateUserCreate(data);

        const response = await this.httpClient.request<User>({
            method: 'POST',
            path: this.basePath,
            body: data,
        });

        return response.data;
    }

    /**
     * Update an existing user
     */
    async update(userId: string, data: UserUpdate): Promise<User> {
        this.validateUserId(userId);

        const response = await this.httpClient.request<User>({
            method: 'PATCH',
            path: `${this.basePath}/${userId}`,
            body: data,
        });

        return response.data;
    }

    /**
     * Delete a user
     */
    async delete(userId: string): Promise<void> {
        this.validateUserId(userId);

        await this.httpClient.request<void>({
            method: 'DELETE',
            path: `${this.basePath}/${userId}`,
        });
    }

    /**
     * Validate user ID format
     */
    private validateUserId(userId: string): void {
        if (!userId) {
            throw new Error('User ID is required');
        }
        if (!/^usr_[a-zA-Z0-9]+$/.test(userId)) {
            throw new Error('Invalid user ID format. Expected: usr_<id>');
        }
    }

    /**
     * Validate user creation data
     */
    private validateUserCreate(data: UserCreate): void {
        if (!data.email) {
            throw new Error('Email is required');
        }
        if (!data.name) {
            throw new Error('Name is required');
        }
    }
}
