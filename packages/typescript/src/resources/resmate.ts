import type { HttpClient } from '../http/client.js';
import type {
    Residence,
    ResidenceSearch,
    ResidenceList,
    CampusList,
} from '../types/index.js';
import { Paginator } from './pagination.js';

export class ResidencesResource {
    private readonly httpClient: HttpClient;
    private readonly basePath = '/v1/resmate/residences';

    constructor(httpClient: HttpClient) {
        this.httpClient = httpClient;
    }

    /**
     * List and search residences
     */
    async list(params?: ResidenceSearch): Promise<ResidenceList> {
        const query: Record<string, string | undefined> = {
            limit: params?.limit?.toString(),
            cursor: params?.cursor,
            query: params?.query,
            nsfas: params?.nsfas?.toString(),
            min_price: params?.min_price?.toString(),
            max_price: params?.max_price?.toString(),
            gender: params?.gender,
            campus_id: params?.campus_id,
            radius: params?.radius?.toString(),
            min_rating: params?.min_rating?.toString(),
        };

        const response = await this.httpClient.request<ResidenceList>({
            method: 'GET',
            path: this.basePath,
            query,
        });

        return response.data;
    }

    /**
     * Iterate over residences
     */
    iterate(params?: ResidenceSearch): Paginator<Residence> {
        return new Paginator<Residence>(
            async (cursor) => {
                return this.list({ ...params, cursor });
            }
        );
    }

    /**
     * Get residence by ID
     */
    async get(id: string): Promise<Residence> {
        const response = await this.httpClient.request<Residence>({
            method: 'GET',
            path: `${this.basePath}/${id}`,
        });

        return response.data;
    }
}

export class CampusesResource {
    private readonly httpClient: HttpClient;
    private readonly basePath = '/v1/resmate/campuses';

    constructor(httpClient: HttpClient) {
        this.httpClient = httpClient;
    }

    /**
     * List all campuses
     */
    async list(): Promise<CampusList> {
        const response = await this.httpClient.request<CampusList>({
            method: 'GET',
            path: this.basePath,
        });
        return response.data;
    }
}

export class ResMateResource {
    public readonly residences: ResidencesResource;
    public readonly campuses: CampusesResource;

    constructor(httpClient: HttpClient) {
        this.residences = new ResidencesResource(httpClient);
        this.campuses = new CampusesResource(httpClient);
    }
}
