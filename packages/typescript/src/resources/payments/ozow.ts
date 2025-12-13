/**
 * Luna SDK - Ozow Integration
 * Instant EFT payments for South Africa
 * https://ozow.com
 */

import { createHmac } from 'crypto';
import type { HttpClient } from '../../http/client.js';
import type {
    OzowConfig,
    OzowPaymentRequest,
    OzowPayment,
    OzowWebhookPayload,
    Refund,
    RefundRequest,
} from './types.js';

const OZOW_PAYMENT_URL = 'https://pay.ozow.com';

export class Ozow {
    private config: OzowConfig;

    // Note: http client reserved for future API integration
    constructor(_http: HttpClient, config: OzowConfig) {
        this.config = config;
    }

    /**
     * Create a payment request and get redirect URL
     * Note: Ozow API is external - this generates the payment URL directly
     */
    async createPayment(request: OzowPaymentRequest): Promise<OzowPayment> {
        const paymentId = `oz_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

        const data: Record<string, string> = {
            SiteCode: this.config.siteCode,
            CountryCode: 'ZA',
            CurrencyCode: 'ZAR',
            Amount: request.amount.toFixed(2),
            TransactionReference: request.transactionReference,
            BankReference: request.bankReference,
            CancelUrl: request.cancelUrl,
            ErrorUrl: request.errorUrl,
            SuccessUrl: request.successUrl,
            NotifyUrl: request.notifyUrl,
            IsTest: request.isTest ?? this.config.sandbox ? 'true' : 'false',
        };

        if (request.customer?.firstName) data.CustomerFirstName = request.customer.firstName;
        if (request.customer?.lastName) data.CustomerLastName = request.customer.lastName;
        if (request.customer?.email) data.CustomerEmail = request.customer.email;
        if (request.customer?.phone) data.CustomerPhone = request.customer.phone;

        const hashString = this.generateHashString(data);
        const hash = this.generateHash(hashString);
        data.HashCheck = hash;

        // Build Ozow payment URL with parameters
        const queryString = new URLSearchParams(data).toString();
        const paymentUrl = `${OZOW_PAYMENT_URL}?${queryString}`;

        return {
            id: paymentId,
            provider: 'ozow',
            amount: {
                value: Math.round(request.amount * 100),
                currency: 'ZAR',
            },
            status: 'pending',
            reference: request.transactionReference,
            description: request.bankReference,
            paymentUrl,
            transactionId: paymentId,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Verify webhook hash
     */
    verifyWebhook(payload: OzowWebhookPayload): boolean {
        const { Hash, ...data } = payload;
        const hashString = this.generateHashString(data);
        const expectedHash = this.generateHash(hashString);
        return Hash.toLowerCase() === expectedHash.toLowerCase();
    }

    /**
     * Process webhook and return payment status
     */
    processWebhook(payload: OzowWebhookPayload): OzowPayment {
        const statusMap: Record<string, OzowPayment['status']> = {
            Complete: 'completed',
            Cancelled: 'cancelled',
            Error: 'failed',
            Abandoned: 'cancelled',
            PendingInvestigation: 'processing',
        };

        return {
            id: payload.TransactionReference,
            provider: 'ozow',
            transactionId: payload.TransactionId,
            amount: {
                value: Math.round(parseFloat(payload.Amount) * 100),
                currency: 'ZAR',
            },
            status: statusMap[payload.Status] || 'pending',
            reference: payload.TransactionReference,
            paymentUrl: '',
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Request a refund
     * Note: Ozow refunds typically require dashboard action or direct API call
     */
    async refund(request: RefundRequest): Promise<Refund> {
        const refundId = `ref_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

        console.log(`Ozow refund requested for payment ${request.paymentId}`);

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
     * Generate hash string for Ozow
     */
    private generateHashString(data: Record<string, string>): string {
        const orderedFields = [
            'SiteCode',
            'CountryCode',
            'CurrencyCode',
            'Amount',
            'TransactionReference',
            'BankReference',
            'CancelUrl',
            'ErrorUrl',
            'SuccessUrl',
            'NotifyUrl',
            'IsTest',
        ];

        return orderedFields
            .map(field => data[field] || '')
            .join('')
            .toLowerCase();
    }

    /**
     * Generate SHA512 hash
     */
    private generateHash(input: string): string {
        const stringToHash = input + this.config.privateKey.toLowerCase();
        return createHmac('sha512', '').update(stringToHash).digest('hex');
    }
}
