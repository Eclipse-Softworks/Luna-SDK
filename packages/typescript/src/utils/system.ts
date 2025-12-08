/**
 * System utility functions for collecting environment information
 */

import * as os from 'node:os';
import * as process from 'node:process';

export interface SystemInfo {
    os: string;
    arch: string;
    runtime: string;
    runtimeVersion: string;
}

/**
 * Get system information safely (handling both Node.js and potentially other runtimes in the future)
 */
export function getSystemInfo(): SystemInfo {
    let systemInfo: SystemInfo = {
        os: 'unknown',
        arch: 'unknown',
        runtime: 'unknown',
        runtimeVersion: 'unknown',
    };

    try {
        // Node.js environment detection
        if (typeof process !== 'undefined' && process.versions && process.versions.node) {
            systemInfo.runtime = 'node';
            systemInfo.runtimeVersion = process.versions.node;
            systemInfo.arch = process.arch;
            systemInfo.os = process.platform;
        }
        // Browser environment detection (fallback)
        else if (typeof window !== 'undefined' && typeof navigator !== 'undefined') {
            systemInfo.runtime = 'browser';
            systemInfo.runtimeVersion = navigator.userAgent;
            systemInfo.os = navigator.platform;
        }
    } catch (e) {
        // Ignore errors obtaining system info
    }

    return systemInfo;
}
