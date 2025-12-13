/**
 * Luna SDK - Yoco Integration
 * Online payments for South African SMEs
 * https://yoco.com
 */

import { createHmac } from 'crypto';
import type { HttpClient } from '../../http/client.js';
import type {
    YocoConfig,
    YocoPaymentRequest,
    YocoPayment,
    YocoWebhookPayload,
    Refund,
    RefundRequest,
    YocoLineItem,
} from './types.js';

const YOCO_API_URL = 'https://payments.yoco.com/api';

export class Yoco {
    private config: YocoConfig;

    // Note: http client reserved for future API integration
    constructor(_http: HttpClient, config: YocoConfig) {
        this.config = config;
    }

    /**
     * Create a checkout session and get redirect URL
     * Note: This would call Yoco's API directly in production
     */
    async createPayment(request: YocoPaymentRequest): Promise<YocoPayment> {
        const checkoutId = `chk_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

        // In production, this would make an actual API call to Yoco
        // For SDK structure, we prepare the request format
        const payload = {
            amount: request.amount,
            currency: request.currency || 'ZAR',
            successUrl: request.successUrl,
            cancelUrl: request.cancelUrl,
            failureUrl: request.failureUrl || request.cancelUrl,
            metadata: request.metadata,
            lineItems: request.lineItems?.map((item: YocoLineItem) => ({
                displayName: item.displayName,
                quantity: item.quantity,
                pricingDetails: {
                    price: item.pricingDetails.price,
                },
            })),
        };

        console.log(`Yoco checkout request prepared: ${YOCO_API_URL}/checkouts`, payload);

        return {
            id: `yc_${checkoutId}`,
            provider: 'yoco',
            checkoutId,
            amount: {
                value: request.amount,
                currency: request.currency || 'ZAR',
            },
            status: 'pending',
            reference: checkoutId,
            redirectUrl: `https://payments.yoco.com/checkout/${checkoutId}`,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Verify webhook signature
     */
    verifyWebhook(payload: string, signature: string): boolean {
        const expectedSignature = createHmac('sha256', this.config.secretKey)
            .update(payload)
            .digest('hex');

        return signature === expectedSignature;
    }

    /**
     * Process webhook event
     */
    processWebhook(payload: YocoWebhookPayload): YocoPayment {
        const statusMap: Record<string, YocoPayment['status']> = {
            'payment.succeeded': 'completed',
            'payment.failed': 'failed',
            'payment.cancelled': 'cancelled',
        };

        return {
            id: `yc_${payload.payload.id}`,
            provider: 'yoco',
            checkoutId: payload.payload.id,
            amount: {
                value: payload.payload.amount,
                currency: payload.payload.currency as 'ZAR',
            },
            status: statusMap[payload.type] || 'pending',
            reference: payload.payload.id,
            metadata: payload.payload.metadata,
            redirectUrl: '',
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Refund a payment
     */
    async refund(request: RefundRequest): Promise<Refund> {
        const refundId = `ref_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

        console.log(`Yoco refund request prepared for checkout ${request.paymentId}`);

        return {
            id: refundId,
            paymentId: request.paymentId,
            amount: {
                value: request.amount || 0,
                currency: 'ZAR',
            },
            status: 'pending',
            reason: request.reason,
            created_at: new Date().toISOString(),
        };
    }
}
