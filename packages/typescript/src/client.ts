import { HttpClient } from './http/client.js';
import { ApiKeyAuth, TokenAuth, type AuthProvider } from './auth/index.js';
import { UsersResource } from './resources/users.js';
import { ProjectsResource } from './resources/projects.js';
import { ResMateResource } from './resources/resmate.js';
import { IdentityResource } from './resources/identity.js';
import { StorageResource } from './resources/storage.js';
import { AiResource } from './resources/ai.js';
import { AutomationResource } from './resources/automation.js';
import { ConsoleLogger, type Logger, type LogLevel, type TelemetryConfig } from './telemetry/index.js';

/**
 * Configuration for the Luna client
 */
export interface ClientConfig {
    /** API key for authentication (format: lk_<env>_<key>) */
    apiKey?: string;
    /** Access token for OAuth authentication */
    accessToken?: string;
    /** Refresh token for automatic token refresh */
    refreshToken?: string;
    /** Callback when tokens are refreshed */
    onTokenRefresh?: (tokens: { accessToken: string; refreshToken: string }) => void | Promise<void>;
    /** Base URL for the API (default: https://api.eclipse.dev) */
    baseUrl?: string;
    /** Request timeout in milliseconds (default: 30000) */
    timeout?: number;
    /** Maximum number of retry attempts (default: 3) */
    maxRetries?: number;
    /** Custom logger instance */
    logger?: Logger;
    /** Log level (default: 'info') */
    logLevel?: LogLevel;
    /** Telemetry configuration */
    telemetry?: TelemetryConfig;
}

/**
 * Additional options that can be passed per-request
 */
export interface ClientOptions {
    /** Override the default timeout for this request */
    timeout?: number;
    /** Custom headers for this request */
    headers?: Record<string, string>;
    /** AbortSignal for request cancellation */
    signal?: AbortSignal;
}

const DEFAULT_BASE_URL = 'https://api.eclipse.dev';
const DEFAULT_TIMEOUT = 30_000;
const DEFAULT_MAX_RETRIES = 3;

/**
 * Luna SDK Client
 *
 * @example
 * ```typescript
 * // API Key authentication
 * const client = new LunaClient({
 *   apiKey: process.env.LUNA_API_KEY,
 * });
 *
 * // Token authentication
 * const client = new LunaClient({
 *   accessToken: session.accessToken,
 *   refreshToken: session.refreshToken,
 *   onTokenRefresh: async (tokens) => {
 *     await saveTokens(tokens);
 *   },
 * });
 *
 * // Usage
 * const users = await client.users.list();
 * const user = await client.users.get('usr_123');
 * ```
 */
export class LunaClient {
    private readonly httpClient: HttpClient;
    private readonly authProvider: AuthProvider;
    private readonly logger: Logger;

    /** Users resource */
    public readonly users: UsersResource;

    /** Projects resource */
    public readonly projects: ProjectsResource;

    /** ResMate resource */
    public readonly resMate: ResMateResource;

    /** Identity resource */
    public readonly identity: IdentityResource;

    /** Storage resource */
    public readonly storage: StorageResource;

    /** AI resource */
    public readonly ai: AiResource;

    /** Automation resource */
    public readonly automation: AutomationResource;

    constructor(config: ClientConfig) {
        // Validate config
        if (!config.apiKey && !config.accessToken) {
            throw new Error('Either apiKey or accessToken must be provided');
        }

        // Set up logger
        this.logger = config.logger ?? new ConsoleLogger(config.logLevel ?? 'info');

        // Set up auth provider
        if (config.apiKey) {
            this.authProvider = new ApiKeyAuth(config.apiKey);
        } else {
            this.authProvider = new TokenAuth({
                accessToken: config.accessToken!,
                refreshToken: config.refreshToken,
                onRefresh: config.onTokenRefresh,
            });
        }

        // Set up HTTP client
        this.httpClient = new HttpClient({
            baseUrl: config.baseUrl ?? DEFAULT_BASE_URL,
            timeout: config.timeout ?? DEFAULT_TIMEOUT,
            maxRetries: config.maxRetries ?? DEFAULT_MAX_RETRIES,
            authProvider: this.authProvider,
            logger: this.logger,
            telemetry: config.telemetry,
        });

        // Initialize resources
        this.users = new UsersResource(this.httpClient);
        this.projects = new ProjectsResource(this.httpClient);
        this.resMate = new ResMateResource(this.httpClient);
        this.identity = new IdentityResource(this.httpClient);
        this.storage = new StorageResource(this.httpClient);
        this.ai = new AiResource(this.httpClient);
        this.automation = new AutomationResource(this.httpClient);

        this.logger.debug('LunaClient initialized', {
            baseUrl: config.baseUrl ?? DEFAULT_BASE_URL,
            authType: config.apiKey ? 'api_key' : 'token',
        });
    }

    /**
     * Get SDK version
     */
    get version(): string {
        return '1.0.0';
    }
}
