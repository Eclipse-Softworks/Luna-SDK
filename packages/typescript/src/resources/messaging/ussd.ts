/**
 * Luna SDK - USSD Integration
 * Mobile USSD services for South African networks
 * Supports: Clickatell, Africa's Talking
 */

import type { HttpClient } from '../../http/client.js';
import type {
    USSDConfig,
    USSDSession,
    USSDResponse,
    USSDHandler,
} from './types.js';

export class USSD {
    private config: USSDConfig;
    private handlers: Map<string, USSDHandler> = new Map();

    // Note: http client reserved for future API integration
    constructor(_http: HttpClient, config: USSDConfig) {
        this.config = config;
    }

    /**
     * Register a handler for USSD sessions
     */
    onSession(handler: USSDHandler): void {
        this.handlers.set('default', handler);
    }

    /**
     * Register a handler for a specific menu path
     */
    onMenu(path: string, handler: USSDHandler): void {
        this.handlers.set(path, handler);
    }

    /**
     * Process incoming USSD request
     * Used by your webhook endpoint
     */
    async processRequest(session: USSDSession): Promise<USSDResponse> {
        const handler =
            this.handlers.get(session.text) || this.handlers.get('default');

        if (!handler) {
            return {
                text: 'Service temporarily unavailable. Please try again later.',
                end: true,
            };
        }

        try {
            return await handler(session);
        } catch (error) {
            console.error('USSD handler error:', error);
            return {
                text: 'An error occurred. Please try again.',
                end: true,
            };
        }
    }

    /**
     * Create a menu response
     */
    static menu(title: string, options: Array<{ key: string; label: string }>): string {
        const lines = [title, ''];
        for (const option of options) {
            lines.push(`${option.key}. ${option.label}`);
        }
        return lines.join('\n');
    }

    /**
     * Parse Africa's Talking webhook format
     */
    parseAfricasTalkingRequest(body: {
        sessionId: string;
        phoneNumber: string;
        serviceCode: string;
        text: string;
        networkCode?: string;
    }): USSDSession {
        return {
            id: `ussd_${Date.now()}`,
            sessionId: body.sessionId,
            phoneNumber: body.phoneNumber,
            serviceCode: body.serviceCode,
            text: body.text,
            state: 'active',
            network: this.mapNetworkCode(body.networkCode),
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Format response for Africa's Talking
     */
    formatAfricasTalkingResponse(response: USSDResponse): string {
        const prefix = response.end ? 'END ' : 'CON ';
        return prefix + response.text;
    }

    /**
     * Parse Clickatell webhook format
     */
    parseClickatellRequest(body: {
        sessionId: string;
        msisdn: string;
        request: string;
        shortcode: string;
    }): USSDSession {
        return {
            id: `ussd_${Date.now()}`,
            sessionId: body.sessionId,
            phoneNumber: body.msisdn,
            serviceCode: body.shortcode,
            text: body.request,
            state: 'active',
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Map network codes to names
     */
    private mapNetworkCode(code?: string): string | undefined {
        if (!code) return undefined;

        const networks: Record<string, string> = {
            '655001': 'Vodacom',
            '655002': 'Telkom',
            '655007': 'Cell C',
            '655010': 'MTN',
        };

        return networks[code] || code;
    }

    /**
     * Get USSD service code
     */
    getServiceCode(): string {
        return this.config.serviceCode;
    }

    /**
     * Create example session flow
     */
    static createExampleMenu(): USSDHandler {
        return (session: USSDSession) => {
            const parts = session.text.split('*').filter(Boolean);

            if (parts.length === 0) {
                return {
                    text: USSD.menu('Welcome to Luna SDK', [
                        { key: '1', label: 'Check Balance' },
                        { key: '2', label: 'Send Payment' },
                        { key: '3', label: 'Mini Statement' },
                        { key: '4', label: 'Exit' },
                    ]),
                    end: false,
                };
            }

            const selection = parts[0];

            switch (selection) {
                case '1':
                    return {
                        text: 'Your balance is R1,234.56',
                        end: true,
                    };
                case '2':
                    if (parts.length === 1) {
                        return {
                            text: 'Enter phone number to send payment:',
                            end: false,
                        };
                    }
                    if (parts.length === 2) {
                        return {
                            text: 'Enter amount (ZAR):',
                            end: false,
                        };
                    }
                    return {
                        text: `Payment of R${parts[2]} to ${parts[1]} initiated.`,
                        end: true,
                    };
                case '3':
                    return {
                        text: 'Mini Statement:\n1. Received R500.00\n2. Sent R100.00\n3. Airtime R50.00',
                        end: true,
                    };
                case '4':
                    return {
                        text: 'Thank you for using Luna SDK. Goodbye!',
                        end: true,
                    };
                default:
                    return {
                        text: 'Invalid selection. Please try again.',
                        end: true,
                    };
            }
        };
    }
}
