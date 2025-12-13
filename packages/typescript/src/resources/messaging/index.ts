/**
 * Luna SDK - Messaging Module
 * Unified interface for SMS, WhatsApp, USSD, and Email
 */

export { SMS } from './sms.js';
export { WhatsApp } from './whatsapp.js';
export { USSD } from './ussd.js';

export * from './types.js';

import type { HttpClient } from '../../http/client.js';
import { SMS } from './sms.js';
import { WhatsApp } from './whatsapp.js';
import { USSD } from './ussd.js';
import type {
    SMSConfig,
    WhatsAppConfig,
    USSDConfig,
    EmailConfig,
} from './types.js';

export interface MessagingConfig {
    sms?: SMSConfig;
    whatsapp?: WhatsAppConfig;
    ussd?: USSDConfig;
    email?: EmailConfig;
}

/**
 * Unified Messaging API for African communication channels
 */
export class Messaging {
    private http: HttpClient;
    private config: MessagingConfig;

    private _sms?: SMS;
    private _whatsapp?: WhatsApp;
    private _ussd?: USSD;

    constructor(http: HttpClient, config: MessagingConfig) {
        this.http = http;
        this.config = config;
    }

    /**
     * SMS messaging - Multi-provider support
     */
    get sms(): SMS {
        if (!this._sms) {
            if (!this.config.sms) {
                throw new Error(
                    'SMS not configured. Provide SMS credentials when initializing LunaClient.'
                );
            }
            this._sms = new SMS(this.http, this.config.sms);
        }
        return this._sms;
    }

    /**
     * WhatsApp Business API
     */
    get whatsapp(): WhatsApp {
        if (!this._whatsapp) {
            if (!this.config.whatsapp) {
                throw new Error(
                    'WhatsApp not configured. Provide WhatsApp credentials when initializing LunaClient.'
                );
            }
            this._whatsapp = new WhatsApp(this.http, this.config.whatsapp);
        }
        return this._whatsapp;
    }

    /**
     * USSD services for SA networks
     */
    get ussd(): USSD {
        if (!this._ussd) {
            if (!this.config.ussd) {
                throw new Error(
                    'USSD not configured. Provide USSD credentials when initializing LunaClient.'
                );
            }
            this._ussd = new USSD(this.http, this.config.ussd);
        }
        return this._ussd;
    }

    /**
     * Check if a specific channel is configured
     */
    isConfigured(channel: 'sms' | 'whatsapp' | 'ussd' | 'email'): boolean {
        return !!this.config[channel];
    }

    /**
     * Get list of configured channels
     */
    getConfiguredChannels(): string[] {
        const channels: string[] = [];
        if (this.config.sms) channels.push('sms');
        if (this.config.whatsapp) channels.push('whatsapp');
        if (this.config.ussd) channels.push('ussd');
        if (this.config.email) channels.push('email');
        return channels;
    }
}
