/**
 * Log levels in order of severity
 */
export type LogLevel = 'error' | 'warn' | 'info' | 'debug' | 'trace';

/**
 * Context data for log messages
 */
export type LogContext = Record<string, unknown>;

/**
 * Logger interface for SDK operations
 */
export interface Logger {
    error(message: string, context?: LogContext): void;
    warn(message: string, context?: LogContext): void;
    info(message: string, context?: LogContext): void;
    debug(message: string, context?: LogContext): void;
    trace(message: string, context?: LogContext): void;
}

/**
 * OpenTelemetry configuration
 */
export interface TelemetryConfig {
    /** Enable telemetry collection */
    enabled: boolean;
    /** OpenTelemetry tracer instance */
    tracer?: unknown;
    /** OpenTelemetry meter instance */
    meter?: unknown;
    /** Propagate trace context in requests */
    propagateContext?: boolean;
}

/**
 * Log level priority (higher = more severe)
 */
export const LOG_LEVEL_PRIORITY: Record<LogLevel, number> = {
    error: 50,
    warn: 40,
    info: 30,
    debug: 20,
    trace: 10,
};
