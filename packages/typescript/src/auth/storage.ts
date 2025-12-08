import * as fs from 'fs';
import * as path from 'path';
import * as os from 'os';
import type { TokenPair } from './types.js';

export interface TokenStore {
    save(tokens: TokenPair): Promise<void>;
    load(): Promise<TokenPair | null>;
    clear(): Promise<void>;
}

/**
 * File-based token storage with secure permissions (0600).
 */
export class FileTokenStore implements TokenStore {
    private readonly filePath: string;

    constructor(filePath?: string) {
        if (filePath) {
            this.filePath = filePath;
        } else {
            const homeDir = os.homedir();
            const configDir = path.join(homeDir, '.luna');
            if (!fs.existsSync(configDir)) {
                fs.mkdirSync(configDir, { recursive: true });
            }
            this.filePath = path.join(configDir, 'credentials.json');
        }
    }

    async save(tokens: TokenPair): Promise<void> {
        const data = JSON.stringify(tokens, null, 2);
        await fs.promises.writeFile(this.filePath, data, { mode: 0o600 });
    }

    async load(): Promise<TokenPair | null> {
        try {
            const data = await fs.promises.readFile(this.filePath, 'utf-8');
            const tokens = JSON.parse(data);
            return {
                accessToken: tokens.accessToken,
                refreshToken: tokens.refreshToken,
                expiresAt: tokens.expiresAt ? new Date(tokens.expiresAt) : undefined,
            };
        } catch (error) {
            return null;
        }
    }

    async clear(): Promise<void> {
        try {
            await fs.promises.unlink(this.filePath);
        } catch (error) {
            // Ignore if file doesn't exist
        }
    }
}
