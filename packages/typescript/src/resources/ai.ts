/**
 * Luna SDK - Expanded AI Module
 * Multi-provider LLM, Embeddings, Vision, and Local Language Support
 */

import type { HttpClient } from '../http/client.js';

// ============================================================================
// AI Provider Types
// ============================================================================

export type AIProvider = 'openai' | 'anthropic' | 'google' | 'ollama' | 'azure';

export interface AIProviderConfig {
    provider: AIProvider;
    apiKey?: string;
    baseUrl?: string;
    model?: string;
    /** Organization ID (OpenAI) */
    organization?: string;
    /** Project ID (Google) */
    projectId?: string;
}

// ============================================================================
// Chat Types
// ============================================================================

export type MessageRole = 'system' | 'user' | 'assistant' | 'tool';

export interface ChatMessage {
    role: MessageRole;
    content: string;
    name?: string;
    toolCalls?: ToolCall[];
    toolCallId?: string;
}

export interface ToolCall {
    id: string;
    type: 'function';
    function: {
        name: string;
        arguments: string;
    };
}

export interface Tool {
    type: 'function';
    function: {
        name: string;
        description?: string;
        parameters?: Record<string, unknown>;
    };
}

export interface ChatCompletionRequest {
    messages: ChatMessage[];
    model?: string;
    temperature?: number;
    maxTokens?: number;
    topP?: number;
    stop?: string | string[];
    stream?: boolean;
    tools?: Tool[];
    toolChoice?: 'auto' | 'none' | { type: 'function'; function: { name: string } };
}

export interface ChatCompletionResponse {
    id: string;
    model: string;
    choices: Array<{
        index: number;
        message: ChatMessage;
        finishReason: 'stop' | 'length' | 'tool_calls' | 'content_filter';
    }>;
    usage: {
        promptTokens: number;
        completionTokens: number;
        totalTokens: number;
    };
}

// ============================================================================
// Embeddings Types
// ============================================================================

export interface EmbeddingRequest {
    input: string | string[];
    model?: string;
    dimensions?: number;
}

export interface EmbeddingResponse {
    data: Array<{
        index: number;
        embedding: number[];
    }>;
    model: string;
    usage: {
        promptTokens: number;
        totalTokens: number;
    };
}

// ============================================================================
// Vision Types
// ============================================================================

export interface VisionRequest {
    image: string | { url: string } | { base64: string };
    prompt: string;
    model?: string;
    maxTokens?: number;
}

export interface VisionResponse {
    description: string;
    labels?: string[];
    objects?: Array<{
        name: string;
        confidence: number;
        boundingBox?: { x: number; y: number; width: number; height: number };
    }>;
}

// ============================================================================
// Translation Types (SA Languages)
// ============================================================================

export type SALanguage = 'en' | 'zu' | 'xh' | 'af' | 'nso' | 'st' | 'tn' | 'ts' | 've' | 'ss' | 'nr';

export const SA_LANGUAGES = {
    en: 'English',
    zu: 'isiZulu',
    xh: 'isiXhosa',
    af: 'Afrikaans',
    nso: 'Sepedi',
    st: 'Sesotho',
    tn: 'Setswana',
    ts: 'Xitsonga',
    ve: 'Tshivenda',
    ss: 'siSwati',
    nr: 'isiNdebele',
} as const;

export interface TranslateRequest {
    text: string;
    from?: SALanguage;
    to: SALanguage;
}

export interface TranslateResponse {
    translatedText: string;
    detectedLanguage?: SALanguage;
    confidence?: number;
}

// ============================================================================
// AI Resource Class
// ============================================================================

export class AiResource {
    private readonly httpClient: HttpClient;
    private readonly basePath = '/v1/ai';
    private defaultProvider: AIProvider = 'openai';
    private providerConfigs: Map<AIProvider, AIProviderConfig> = new Map();

    constructor(httpClient: HttpClient, config?: AIProviderConfig) {
        this.httpClient = httpClient;
        if (config) {
            this.configureProvider(config);
            this.defaultProvider = config.provider;
        }
    }

    /**
     * Configure an AI provider
     */
    configureProvider(config: AIProviderConfig): void {
        this.providerConfigs.set(config.provider, config);
    }

    /**
     * Set the default provider
     */
    setDefaultProvider(provider: AIProvider): void {
        this.defaultProvider = provider;
    }

    // ============================================================================
    // Chat Completions
    // ============================================================================

    /**
     * Generate chat completions (multi-provider)
     */
    async chat(request: ChatCompletionRequest): Promise<ChatCompletionResponse> {
        const response = await this.httpClient.request<ChatCompletionResponse>({
            method: 'POST',
            path: `${this.basePath}/chat/completions`,
            body: {
                ...request,
                provider: this.defaultProvider,
            },
        });
        return response.data;
    }

    /**
     * Generate a simple text response
     */
    async generate(prompt: string, options?: Partial<ChatCompletionRequest>): Promise<string> {
        const response = await this.chat({
            messages: [{ role: 'user', content: prompt }],
            ...options,
        });
        return response.choices[0]?.message.content || '';
    }

    /**
     * Chat with system context
     */
    async chatWithContext(
        systemPrompt: string,
        userMessage: string,
        options?: Partial<ChatCompletionRequest>
    ): Promise<string> {
        const response = await this.chat({
            messages: [
                { role: 'system', content: systemPrompt },
                { role: 'user', content: userMessage },
            ],
            ...options,
        });
        return response.choices[0]?.message.content || '';
    }

    // ============================================================================
    // Embeddings
    // ============================================================================

    /**
     * Generate embeddings for text
     */
    async embed(request: EmbeddingRequest): Promise<EmbeddingResponse> {
        const response = await this.httpClient.request<EmbeddingResponse>({
            method: 'POST',
            path: `${this.basePath}/embeddings`,
            body: request,
        });
        return response.data;
    }

    /**
     * Generate embedding for a single text
     */
    async embedText(text: string, model?: string): Promise<number[]> {
        const response = await this.embed({ input: text, model });
        return response.data[0]?.embedding || [];
    }

    /**
     * Calculate cosine similarity between two embeddings
     */
    cosineSimilarity(a: number[], b: number[]): number {
        if (a.length !== b.length) {
            throw new Error('Embeddings must have the same dimensions');
        }

        let dotProduct = 0;
        let normA = 0;
        let normB = 0;

        for (let i = 0; i < a.length; i++) {
            dotProduct += a[i]! * b[i]!;
            normA += a[i]! * a[i]!;
            normB += b[i]! * b[i]!;
        }

        return dotProduct / (Math.sqrt(normA) * Math.sqrt(normB));
    }

    // ============================================================================
    // Vision
    // ============================================================================

    /**
     * Analyze an image
     */
    async analyzeImage(request: VisionRequest): Promise<VisionResponse> {
        const response = await this.httpClient.request<VisionResponse>({
            method: 'POST',
            path: `${this.basePath}/vision/analyze`,
            body: request,
        });
        return response.data;
    }

    /**
     * Describe an image
     */
    async describeImage(image: string | { url: string }): Promise<string> {
        const response = await this.analyzeImage({
            image,
            prompt: 'Describe this image in detail.',
        });
        return response.description;
    }

    // ============================================================================
    // Translation (SA Languages)
    // ============================================================================

    /**
     * Translate text between South African languages
     */
    async translate(request: TranslateRequest): Promise<TranslateResponse> {
        const response = await this.httpClient.request<TranslateResponse>({
            method: 'POST',
            path: `${this.basePath}/translate`,
            body: request,
        });
        return response.data;
    }

    /**
     * Quick translation helper
     */
    async translateTo(text: string, targetLanguage: SALanguage): Promise<string> {
        const response = await this.translate({ text, to: targetLanguage });
        return response.translatedText;
    }

    /**
     * Detect language of text
     */
    async detectLanguage(text: string): Promise<SALanguage> {
        const response = await this.translate({ text, to: 'en' });
        return response.detectedLanguage || 'en';
    }

    // ============================================================================
    // Utility Methods
    // ============================================================================

    /**
     * Get available models for the default provider
     */
    getAvailableModels(): string[] {
        const models: Record<AIProvider, string[]> = {
            openai: ['gpt-4o', 'gpt-4-turbo', 'gpt-4', 'gpt-3.5-turbo', 'text-embedding-3-small'],
            anthropic: ['claude-3-5-sonnet-20241022', 'claude-3-opus-20240229', 'claude-3-haiku-20240307'],
            google: ['gemini-2.0-flash-exp', 'gemini-1.5-pro', 'gemini-1.5-flash'],
            ollama: ['llama3.2', 'mistral', 'codellama', 'llama2'],
            azure: ['gpt-4', 'gpt-35-turbo'],
        };
        return models[this.defaultProvider] || [];
    }

    /**
     * Get supported South African languages
     */
    getSupportedLanguages(): typeof SA_LANGUAGES {
        return SA_LANGUAGES;
    }
}
