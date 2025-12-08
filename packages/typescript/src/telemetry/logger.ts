import type { Logger, LogLevel, LogContext } from './types.js';
import { LOG_LEVEL_PRIORITY } from './types.js';

/**
 * Patterns for sensitive data that should be redacted
 */
const REDACT_PATTERNS = [
    /api[_-]?key/i,
    /authorization/i,
    /x-luna-api-key/i,
    /password/i,
    /secret/i,
    /token/i,
    /bearer/i,
];

/**
 * Console-based logger with redaction support
 */
export class ConsoleLogger implements Logger {
    private readonly levelPriority: number;

    constructor(level: LogLevel = 'info') {
        this.levelPriority = LOG_LEVEL_PRIORITY[level];
    }

    error(message: string, context?: LogContext): void {
        this.log('error', message, context);
    }

    warn(message: string, context?: LogContext): void {
        this.log('warn', message, context);
    }

    info(message: string, context?: LogContext): void {
        this.log('info', message, context);
    }

    debug(message: string, context?: LogContext): void {
        this.log('debug', message, context);
    }

    trace(message: string, context?: LogContext): void {
        this.log('trace', message, context);
    }

    private log(level: LogLevel, message: string, context?: LogContext): void {
        // Check if this level should be logged
        if (LOG_LEVEL_PRIORITY[level] < this.levelPriority) {
            return;
        }

        const timestamp = new Date().toISOString();
        const sanitizedContext = context ? this.sanitize(context) : undefined;

        const logEntry = {
            timestamp,
            level: level.toUpperCase(),
            message,
            sdk: 'luna-sdk',
            version: '1.0.0',
            language: 'typescript',
            ...(sanitizedContext && { context: sanitizedContext }),
        };

        const output = JSON.stringify(logEntry);

        switch (level) {
            case 'error':
                console.error(output);
                break;
            case 'warn':
                console.warn(output);
                break;
            case 'trace':
            case 'debug':
                console.debug(output);
                break;
            default:
                console.log(output);
        }
    }

    /**
     * Sanitize context by redacting sensitive values
     */
    private sanitize(obj: LogContext): LogContext {
        const result: LogContext = {};

        for (const [key, value] of Object.entries(obj)) {
            if (this.isSensitiveKey(key)) {
                result[key] = '[REDACTED]';
            } else if (value && typeof value === 'object' && !Array.isArray(value)) {
                result[key] = this.sanitize(value as LogContext);
            } else {
                result[key] = value;
            }
        }

        return result;
    }

    /**
     * Check if a key contains sensitive data
     */
    private isSensitiveKey(key: string): boolean {
        return REDACT_PATTERNS.some((pattern) => pattern.test(key));
    }
}
