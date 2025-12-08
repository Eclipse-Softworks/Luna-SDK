import type { HttpClient } from '../http/client.js';
import type {
    User,
    UserCreate,
    UserUpdate,
    UserList,
    PaginationParams,
} from '../types/index.js';

/**
 * Users resource for managing user accounts
 *
 * @example
 * ```typescript
 * // List users
 * const users = await client.users.list({ limit: 10 });
 *
 * // Get a user
 * const user = await client.users.get('usr_123');
 *
 * // Create a user
 * const newUser = await client.users.create({
 *   email: 'john@example.com',
 *   name: 'John Doe',
 * });
 *
 * // Update a user
 * const updated = await client.users.update('usr_123', { name: 'Jane Doe' });
 *
 * // Delete a user
 * await client.users.delete('usr_123');
 * ```
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
