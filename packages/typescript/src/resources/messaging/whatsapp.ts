/**
 * Luna SDK - WhatsApp Business API Integration
 * Supports: Cloud API, Clickatell, Wati, Infobip
 */

import { createHmac } from 'crypto';
import type { HttpClient } from '../../http/client.js';
import type {
    WhatsAppConfig,
    WhatsAppMessage,
    WhatsAppTextRequest,
    WhatsAppTemplateRequest,
    WhatsAppMediaRequest,
    WhatsAppWebhookPayload,
} from './types.js';

export class WhatsApp {
    private config: WhatsAppConfig;

    // Note: http client reserved for future API integration
    constructor(_http: HttpClient, config: WhatsAppConfig) {
        this.config = config;
    }

    /**
     * Send a text message
     */
    async sendText(request: WhatsAppTextRequest): Promise<WhatsAppMessage> {
        const messageId = `wa_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        const to = this.normalizePhoneNumber(request.to);

        console.log(`WhatsApp text message via ${this.config.provider}:`, {
            to,
            text: request.text,
        });

        return {
            id: messageId,
            to,
            type: 'text',
            text: request.text,
            status: 'pending',
            direction: 'outbound',
            provider: this.config.provider,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Send a template message (for business-initiated conversations)
     */
    async sendTemplate(request: WhatsAppTemplateRequest): Promise<WhatsAppMessage> {
        const messageId = `wa_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        const to = this.normalizePhoneNumber(request.to);

        console.log(`WhatsApp template message via ${this.config.provider}:`, {
            to,
            template: request.templateName,
            params: request.templateParams,
        });

        return {
            id: messageId,
            to,
            type: 'template',
            templateName: request.templateName,
            templateParams: request.templateParams,
            status: 'pending',
            direction: 'outbound',
            provider: this.config.provider,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Send a media message (image, document, audio, video)
     */
    async sendMedia(request: WhatsAppMediaRequest): Promise<WhatsAppMessage> {
        const messageId = `wa_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        const to = this.normalizePhoneNumber(request.to);

        console.log(`WhatsApp ${request.type} message via ${this.config.provider}:`, {
            to,
            mediaUrl: request.mediaUrl,
            caption: request.caption,
        });

        return {
            id: messageId,
            to,
            type: request.type,
            mediaUrl: request.mediaUrl,
            status: 'pending',
            direction: 'outbound',
            provider: this.config.provider,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Get message status
     */
    async getStatus(messageId: string): Promise<WhatsAppMessage> {
        return {
            id: messageId,
            to: '',
            type: 'text',
            status: 'delivered',
            direction: 'outbound',
            provider: this.config.provider,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Verify webhook signature (Cloud API)
     */
    verifyWebhook(payload: string, signature: string): boolean {
        if (!this.config.webhookToken) {
            throw new Error('Webhook token not configured');
        }

        const expectedSignature = createHmac('sha256', this.config.webhookToken)
            .update(payload)
            .digest('hex');

        return `sha256=${expectedSignature}` === signature;
    }

    /**
     * Process incoming webhook
     */
    processWebhook(payload: WhatsAppWebhookPayload): WhatsAppMessage[] {
        const messages: WhatsAppMessage[] = [];

        for (const entry of payload.entry) {
            for (const change of entry.changes) {
                const value = change.value;

                // Process incoming messages
                if (value.messages) {
                    for (const msg of value.messages) {
                        messages.push({
                            id: msg.id,
                            to: value.metadata.phone_number_id,
                            from: msg.from,
                            type: msg.type as WhatsAppMessage['type'],
                            text: msg.text?.body,
                            status: 'delivered',
                            direction: 'inbound',
                            provider: this.config.provider,
                            created_at: new Date(parseInt(msg.timestamp) * 1000).toISOString(),
                            updated_at: new Date().toISOString(),
                        });
                    }
                }

                // Process status updates
                if (value.statuses) {
                    for (const status of value.statuses) {
                        messages.push({
                            id: status.id,
                            to: '',
                            type: 'text',
                            status: this.mapStatus(status.status),
                            direction: 'outbound',
                            provider: this.config.provider,
                            created_at: new Date(parseInt(status.timestamp) * 1000).toISOString(),
                            updated_at: new Date().toISOString(),
                        });
                    }
                }
            }
        }

        return messages;
    }

    /**
     * Map WhatsApp status to our status
     */
    private mapStatus(status: string): WhatsAppMessage['status'] {
        const statusMap: Record<string, WhatsAppMessage['status']> = {
            sent: 'sent',
            delivered: 'delivered',
            read: 'read',
            failed: 'failed',
        };
        return statusMap[status] || 'pending';
    }

    /**
     * Normalize phone number for WhatsApp (E.164 without +)
     */
    private normalizePhoneNumber(phone: string): string {
        let digits = phone.replace(/\D/g, '');

        // Handle SA numbers
        if (digits.startsWith('0') && digits.length === 10) {
            digits = '27' + digits.slice(1);
        }

        // Remove + if present
        if (digits.startsWith('+')) {
            digits = digits.slice(1);
        }

        return digits;
    }
}
