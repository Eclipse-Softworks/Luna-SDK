import type { HttpClient } from '../http/client.js';
import type { Bucket, BucketList, FileObject, ListResponse } from '../types/index.js';

export class BucketsResource {
    private readonly httpClient: HttpClient;
    private readonly basePath = '/v1/storage/buckets';

    constructor(httpClient: HttpClient) {
        this.httpClient = httpClient;
    }

    /**
     * List all buckets
     */
    async list(): Promise<BucketList> {
        const response = await this.httpClient.request<BucketList>({
            method: 'GET',
            path: this.basePath,
        });
        return response.data;
    }

    async get(id: string): Promise<Bucket> {
        const response = await this.httpClient.request<Bucket>({
            method: 'GET',
            path: `${this.basePath}/${id}`,
        });
        return response.data;
    }

    /**
     * Upload a file to a bucket
     * Note: In a real implementation, this might handle Multipart forms.
     * For now, we'll assume a direct binary upload or presigned URL flow wrapper.
     * Following the spec: POST /v1/storage/buckets/:id/upload
     */
    async upload(bucketId: string, file: Blob | Buffer | unknown, metadata?: Record<string, string>): Promise<FileObject> {
        // Implementation detail: constructing FormData would happen here
        // For SDK simplicity in this exercise, we pass the body directly or as structured data
        const response = await this.httpClient.request<FileObject>({
            method: 'POST',
            path: `${this.basePath}/${bucketId}/upload`,
            body: { file, metadata },
        });
        return response.data;
    }
}

export class FilesResource {
    private readonly httpClient: HttpClient;
    private readonly basePath = '/v1/storage/files';

    constructor(httpClient: HttpClient) {
        this.httpClient = httpClient;
    }

    /**
     * List files in a bucket
     */
    async list(bucketId: string): Promise<ListResponse<FileObject>> {
        const response = await this.httpClient.request<ListResponse<FileObject>>({
            method: 'GET',
            path: `/v1/storage/buckets/${bucketId}/files`,
        });
        return response.data;
    }

    /**
     * Get download URL for a file
     */
    async getDownloadUrl(id: string): Promise<{ url: string }> {
        const response = await this.httpClient.request<{ url: string }>({
            method: 'GET',
            path: `${this.basePath}/${id}/download`,
        });
        return response.data;
    }
}

export class StorageResource {
    public readonly buckets: BucketsResource;
    public readonly files: FilesResource;

    constructor(httpClient: HttpClient) {
        this.buckets = new BucketsResource(httpClient);
        this.files = new FilesResource(httpClient);
    }
}
