import type { HttpClient } from '../http/client.js';
import type {
    Project,
    ProjectCreate,
    ProjectUpdate,
    ProjectList,
    PaginationParams,
} from '../types/index.js';

/**
 * Projects resource for managing projects
 *
 * @example
 * ```typescript
 * // List projects
 * const projects = await client.projects.list({ limit: 10 });
 *
 * // Get a project
 * const project = await client.projects.get('prj_123');
 *
 * // Create a project
 * const newProject = await client.projects.create({
 *   name: 'My Project',
 *   description: 'A new project',
 * });
 *
 * // Update a project
 * const updated = await client.projects.update('prj_123', { name: 'Updated' });
 *
 * // Delete a project
 * await client.projects.delete('prj_123');
 * ```
 */
export class ProjectsResource {
    private readonly httpClient: HttpClient;
    private readonly basePath = '/v1/projects';

    constructor(httpClient: HttpClient) {
        this.httpClient = httpClient;
    }

    /**
     * List all projects with pagination
     */
    async list(params?: PaginationParams): Promise<ProjectList> {
        const response = await this.httpClient.request<ProjectList>({
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
     * Get a project by ID
     */
    async get(projectId: string): Promise<Project> {
        this.validateProjectId(projectId);

        const response = await this.httpClient.request<Project>({
            method: 'GET',
            path: `${this.basePath}/${projectId}`,
        });

        return response.data;
    }

    /**
     * Create a new project
     */
    async create(data: ProjectCreate): Promise<Project> {
        this.validateProjectCreate(data);

        const response = await this.httpClient.request<Project>({
            method: 'POST',
            path: this.basePath,
            body: data,
        });

        return response.data;
    }

    /**
     * Update an existing project
     */
    async update(projectId: string, data: ProjectUpdate): Promise<Project> {
        this.validateProjectId(projectId);

        const response = await this.httpClient.request<Project>({
            method: 'PATCH',
            path: `${this.basePath}/${projectId}`,
            body: data,
        });

        return response.data;
    }

    /**
     * Delete a project
     */
    async delete(projectId: string): Promise<void> {
        this.validateProjectId(projectId);

        await this.httpClient.request<void>({
            method: 'DELETE',
            path: `${this.basePath}/${projectId}`,
        });
    }

    /**
     * Validate project ID format
     */
    private validateProjectId(projectId: string): void {
        if (!projectId) {
            throw new Error('Project ID is required');
        }
        if (!/^prj_[a-zA-Z0-9]+$/.test(projectId)) {
            throw new Error('Invalid project ID format. Expected: prj_<id>');
        }
    }

    /**
     * Validate project creation data
     */
    private validateProjectCreate(data: ProjectCreate): void {
        if (!data.name) {
            throw new Error('Name is required');
        }
    }
}
