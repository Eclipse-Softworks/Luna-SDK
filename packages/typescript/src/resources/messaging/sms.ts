/**
 * Luna SDK - SMS Integration
 * Multi-provider SMS gateway for Africa
 * Supports: Clickatell, Africa's Talking, Twilio
 */

import type { HttpClient } from '../../http/client.js';
import type {
    SMSConfig,
    SMSMessage,
    SMSSendRequest,
    SMSBulkResult,
    SMSProvider,
} from './types.js';

export class SMS {
    private config: SMSConfig;

    // Note: http client reserved for future API integration
    constructor(_http: HttpClient, config: SMSConfig) {
        this.config = config;
    }

    /**
     * Send a single SMS message
     */
    async send(request: SMSSendRequest): Promise<SMSMessage> {
        const recipient = Array.isArray(request.to) ? request.to[0] : request.to;
        if (!recipient) {
            throw new Error('SMS recipient (to) is required');
        }
        const to = recipient;
        const messageId = `sms_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

        // Normalize phone number for SA
        const normalizedTo = this.normalizePhoneNumber(to);

        // Provider-specific logic
        switch (this.config.provider) {
            case 'clickatell':
                return this.sendViaClickatell(messageId, normalizedTo, request);
            case 'africastalking':
                return this.sendViaAfricasTalking(messageId, normalizedTo, request);
            case 'twilio':
                return this.sendViaTwilio(messageId, normalizedTo, request);
            default:
                throw new Error(`Unsupported SMS provider: ${this.config.provider}`);
        }
    }

    /**
     * Send SMS to multiple recipients
     */
    async sendBulk(request: SMSSendRequest): Promise<SMSBulkResult> {
        const recipients = Array.isArray(request.to) ? request.to : [request.to];
        const successful: SMSMessage[] = [];
        const failed: Array<{ to: string; error: string }> = [];

        for (const to of recipients) {
            try {
                const message = await this.send({ ...request, to });
                successful.push(message);
            } catch (error) {
                failed.push({
                    to,
                    error: error instanceof Error ? error.message : 'Unknown error',
                });
            }
        }

        return { successful, failed };
    }

    /**
     * Get SMS delivery status
     */
    async getStatus(messageId: string): Promise<SMSMessage> {
        // Mock implementation - would call provider API
        return {
            id: messageId,
            to: '',
            body: '',
            status: 'delivered',
            direction: 'outbound',
            provider: this.config.provider,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Get account balance
     */
    async getBalance(): Promise<{ balance: number; currency: string }> {
        // Mock implementation - would call provider API
        console.log(`Checking ${this.config.provider} SMS balance`);
        return { balance: 100.0, currency: 'ZAR' };
    }

    // Provider implementations

    private async sendViaClickatell(
        messageId: string,
        to: string,
        request: SMSSendRequest
    ): Promise<SMSMessage> {
        // Clickatell API format
        const payload = {
            content: request.body,
            to: [to],
            from: request.from || this.config.senderId,
            callback: request.callbackUrl,
        };

        console.log('Clickatell SMS request:', payload);

        return this.createMessageResponse(messageId, to, request, 'clickatell');
    }

    private async sendViaAfricasTalking(
        messageId: string,
        to: string,
        request: SMSSendRequest
    ): Promise<SMSMessage> {
        // Africa's Talking API format
        const payload = {
            username: this.config.username,
            message: request.body,
            to,
            from: request.from || this.config.senderId,
        };

        console.log("Africa's Talking SMS request:", payload);

        return this.createMessageResponse(messageId, to, request, 'africastalking');
    }

    private async sendViaTwilio(
        messageId: string,
        to: string,
        request: SMSSendRequest
    ): Promise<SMSMessage> {
        // Twilio API format
        const payload = {
            Body: request.body,
            To: to,
            From: request.from || this.config.senderId,
            StatusCallback: request.callbackUrl,
        };

        console.log('Twilio SMS request:', payload);

        return this.createMessageResponse(messageId, to, request, 'twilio');
    }

    private createMessageResponse(
        id: string,
        to: string,
        request: SMSSendRequest,
        provider: SMSProvider
    ): SMSMessage {
        return {
            id,
            to,
            from: request.from || this.config.senderId,
            body: request.body,
            status: 'pending',
            direction: 'outbound',
            provider,
            parts: Math.ceil(request.body.length / 160),
            metadata: request.metadata,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Normalize South African phone numbers to E.164 format
     */
    private normalizePhoneNumber(phone: string): string {
        // Remove all non-digit characters
        let digits = phone.replace(/\D/g, '');

        // Handle SA numbers
        if (digits.startsWith('0') && digits.length === 10) {
            digits = '27' + digits.slice(1);
        }

        // Ensure + prefix
        if (!digits.startsWith('+')) {
            digits = '+' + digits;
        }

        return digits;
    }
}
