import type { HttpClient } from '../http/client.js';
import type { WorkflowList, WorkflowRun } from '../types/index.js';

export class WorkflowsResource {
    private readonly httpClient: HttpClient;
    private readonly basePath = '/v1/automation/workflows';

    constructor(httpClient: HttpClient) {
        this.httpClient = httpClient;
    }

    /**
     * List all workflows
     */
    async list(): Promise<WorkflowList> {
        const response = await this.httpClient.request<WorkflowList>({
            method: 'GET',
            path: this.basePath,
        });
        return response.data;
    }

    /**
     * Trigger a workflow
     */
    async trigger(id: string, params: Record<string, unknown> = {}): Promise<WorkflowRun> {
        const response = await this.httpClient.request<WorkflowRun>({
            method: 'POST',
            path: `${this.basePath}/${id}/trigger`,
            body: params,
        });
        return response.data;
    }
}

export class AutomationResource {
    public readonly workflows: WorkflowsResource;

    constructor(httpClient: HttpClient) {
        this.workflows = new WorkflowsResource(httpClient);
    }
}
