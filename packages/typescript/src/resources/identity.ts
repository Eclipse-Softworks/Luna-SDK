import type { HttpClient } from '../http/client.js';
import type { Group, GroupList, GroupCreate, ListResponse } from '../types/index.js';

export class GroupsResource {
    private readonly httpClient: HttpClient;
    private readonly basePath = '/v1/identity/groups';

    constructor(httpClient: HttpClient) {
        this.httpClient = httpClient;
    }

    /**
     * List all groups
     */
    async list(): Promise<GroupList> {
        const response = await this.httpClient.request<GroupList>({
            method: 'GET',
            path: this.basePath,
        });
        return response.data;
    }

    /**
     * Get group by ID
     */
    async get(id: string): Promise<Group> {
        const response = await this.httpClient.request<Group>({
            method: 'GET',
            path: `${this.basePath}/${id}`,
        });
        return response.data;
    }

    /**
     * Create a new group
     */
    async create(params: GroupCreate): Promise<Group> {
        const response = await this.httpClient.request<Group>({
            method: 'POST',
            path: this.basePath,
            body: params,
        });
        return response.data;
    }
}

export class IdentityResource {
    public readonly groups: GroupsResource;

    constructor(httpClient: HttpClient) {
        this.groups = new GroupsResource(httpClient);
    }
}
