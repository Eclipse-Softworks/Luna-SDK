import type { HttpClient } from '../http/client.js';
import type { CompletionRequest, CompletionResponse } from '../types/index.js';

export class AiResource {
    private readonly httpClient: HttpClient;
    private readonly basePath = '/v1/ai';

    constructor(httpClient: HttpClient) {
        this.httpClient = httpClient;
    }

    /**
     * Generate chat completions
     */
    async chatCompletions(params: CompletionRequest): Promise<CompletionResponse> {
        const response = await this.httpClient.request<CompletionResponse>({
            method: 'POST',
            path: `${this.basePath}/chat/completions`,
            body: params,
        });
        return response.data;
    }
}
