/**
 * Luna SDK - PayFast Integration
 * South Africa's leading payment gateway
 * https://payfast.co.za
 */

import { createHmac } from 'crypto';
import type { HttpClient } from '../../http/client.js';
import type {
    PayFastConfig,
    PayFastPaymentRequest,
    PayFastPayment,
    PayFastWebhookPayload,
    Refund,
    RefundRequest,
} from './types.js';

const PAYFAST_LIVE_URL = 'https://www.payfast.co.za/eng/process';
const PAYFAST_SANDBOX_URL = 'https://sandbox.payfast.co.za/eng/process';
const PAYFAST_API_URL = 'https://api.payfast.co.za';
const PAYFAST_SANDBOX_API_URL = 'https://sandbox.payfast.co.za';

export class PayFast {
    private config: PayFastConfig;

    // Note: http client reserved for future API integration
    constructor(http: HttpClient, config: PayFastConfig) {
        void http; // Reserved for future API integration
        this.config = config;
    }

    /**
     * Create a payment request and get redirect URL
     */
    async createPayment(request: PayFastPaymentRequest): Promise<PayFastPayment> {
        const paymentId = `pf_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

        const data: Record<string, string> = {
            merchant_id: this.config.merchantId,
            merchant_key: this.config.merchantKey,
            return_url: request.returnUrl,
            cancel_url: request.cancelUrl,
            notify_url: request.notifyUrl,
            m_payment_id: paymentId,
            amount: request.amount.toFixed(2),
            item_name: request.itemName,
        };

        if (request.itemDescription) data.item_description = request.itemDescription;
        if (request.emailAddress) data.email_address = request.emailAddress;
        if (request.cellNumber) data.cell_number = request.cellNumber;
        if (request.customStr1) data.custom_str1 = request.customStr1;
        if (request.customStr2) data.custom_str2 = request.customStr2;
        if (request.customStr3) data.custom_str3 = request.customStr3;
        if (request.customInt1) data.custom_int1 = String(request.customInt1);
        if (request.customInt2) data.custom_int2 = String(request.customInt2);
        if (request.paymentMethod) data.payment_method = request.paymentMethod;

        const signature = this.generateSignature(data);
        data.signature = signature;

        const baseUrl = this.config.sandbox ? PAYFAST_SANDBOX_URL : PAYFAST_LIVE_URL;
        const queryString = new URLSearchParams(data).toString();
        const paymentUrl = `${baseUrl}?${queryString}`;

        return {
            id: paymentId,
            provider: 'payfast',
            amount: {
                value: Math.round(request.amount * 100),
                currency: request.currency || 'ZAR',
            },
            status: 'pending',
            reference: paymentId,
            description: request.itemDescription,
            paymentUrl,
            signature,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Verify webhook signature
     */
    verifyWebhook(payload: PayFastWebhookPayload): boolean {
        const { signature, ...data } = payload;
        const expectedSignature = this.generateSignature(data);
        return signature === expectedSignature;
    }

    /**
     * Process webhook and return payment status
     */
    processWebhook(payload: PayFastWebhookPayload): PayFastPayment {
        const statusMap: Record<string, PayFastPayment['status']> = {
            COMPLETE: 'completed',
            FAILED: 'failed',
            PENDING: 'pending',
            CANCELLED: 'cancelled',
        };

        return {
            id: payload.m_payment_id,
            provider: 'payfast',
            pfPaymentId: payload.pf_payment_id,
            amount: {
                value: Math.round(parseFloat(payload.amount_gross) * 100),
                currency: 'ZAR',
            },
            status: statusMap[payload.payment_status] || 'pending',
            reference: payload.m_payment_id,
            description: payload.item_name,
            paymentUrl: '',
            signature: payload.signature,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Refund a payment (full or partial)
     * Note: Uses external PayFast API - real implementation would use http.request
     */
    async refund(request: RefundRequest): Promise<Refund> {
        const baseUrl = this.config.sandbox ? PAYFAST_SANDBOX_API_URL : PAYFAST_API_URL;

        // PayFast refund API endpoint (external to Luna platform)
        const refundId = `ref_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

        // In production, this would call the PayFast API directly
        // For now, we return a pending refund that requires PayFast dashboard action
        console.log(`PayFast refund requested: ${baseUrl}/refunds for payment ${request.paymentId}`);

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

    /**
     * Generate MD5 signature for PayFast
     */
    private generateSignature(data: Record<string, string>): string {
        const sortedKeys = Object.keys(data).sort();
        const paramString = sortedKeys
            .filter(key => data[key] !== '' && data[key] !== undefined)
            .map(key => `${key}=${encodeURIComponent(data[key] ?? '').replace(/%20/g, '+')}`)
            .join('&');

        const stringWithPassphrase = this.config.passphrase
            ? `${paramString}&passphrase=${encodeURIComponent(this.config.passphrase)}`
            : paramString;

        return createHmac('md5', '').update(stringWithPassphrase).digest('hex');
    }
}
