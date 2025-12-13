/**
 * Luna SDK - PayShap Integration
 * Real-time account-to-account payments (Capitec Pay compatible)
 * https://www.pasa.org.za/payshap
 */

import type { HttpClient } from '../../http/client.js';
import type {
    PayShapConfig,
    PayShapPaymentRequest,
    PayShapPayment,
} from './types.js';

// PayShap operates through participant banks - this is a conceptual API
const PAYSHAP_API_URL = 'https://api.payshap.co.za';

// South African banks with PayShap support
export const SA_BANKS = {
    ABSA: 'absa',
    CAPITEC: 'capitec',
    FNB: 'fnb',
    NEDBANK: 'nedbank',
    STANDARD_BANK: 'standard',
    INVESTEC: 'investec',
    DISCOVERY_BANK: 'discovery',
    TYMEBANK: 'tymebank',
    AFRICAN_BANK: 'african',
} as const;

export type SABank = (typeof SA_BANKS)[keyof typeof SA_BANKS];

export class PayShap {
    private config: PayShapConfig;

    // Note: http client reserved for future API integration
    constructor(_http: HttpClient, config: PayShapConfig) {
        this.config = config;
    }

    /**
     * Create a PayShap payment request
     * Can use ShapID (proxy) or direct account details
     */
    async createPayment(request: PayShapPaymentRequest): Promise<PayShapPayment> {
        const paymentId = `ps_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        const expiryMinutes = request.expiryMinutes || 30;
        const expiresAt = new Date(Date.now() + expiryMinutes * 60 * 1000);

        // Generate a mock QR code (in production this would come from the bank's API)
        const qrData = JSON.stringify({
            type: 'payshap',
            merchantId: this.config.merchantId,
            amount: request.amount,
            reference: request.reference,
        });
        const qrCode = Buffer.from(qrData).toString('base64');

        console.log(`PayShap payment request prepared: ${PAYSHAP_API_URL}/v1/payments`);

        return {
            id: paymentId,
            provider: 'payshap',
            shapId: `shp_${Math.random().toString(36).substr(2, 8)}`,
            amount: {
                value: Math.round(request.amount * 100),
                currency: 'ZAR',
            },
            status: 'pending',
            reference: request.reference,
            qrCode,
            expiresAt: expiresAt.toISOString(),
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Get payment status (mock implementation)
     */
    async getPayment(paymentId: string): Promise<PayShapPayment> {
        // In production, this would query the bank's API
        return {
            id: paymentId,
            provider: 'payshap',
            shapId: `shp_${Math.random().toString(36).substr(2, 8)}`,
            amount: {
                value: 0,
                currency: 'ZAR',
            },
            status: 'pending',
            reference: paymentId,
            expiresAt: new Date(Date.now() + 30 * 60 * 1000).toISOString(),
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Cancel a pending payment
     */
    async cancelPayment(paymentId: string): Promise<PayShapPayment> {
        console.log(`PayShap payment cancellation requested: ${paymentId}`);

        const payment = await this.getPayment(paymentId);
        return {
            ...payment,
            status: 'cancelled',
            updated_at: new Date().toISOString(),
        };
    }

    /**
     * Lookup a ShapID (payment proxy)
     */
    async lookupShapId(shapId: string): Promise<{
        valid: boolean;
        bankName?: string;
        accountHolderName?: string;
    }> {
        // Mock implementation - in production would call bank APIs
        console.log(`PayShap ShapID lookup: ${shapId}`);

        // Simulate validation based on format
        const isValid = /^[a-zA-Z0-9@._-]+$/.test(shapId) && shapId.length >= 5;

        return {
            valid: isValid,
            bankName: isValid ? 'Sample Bank' : undefined,
            accountHolderName: isValid ? 'Account Holder' : undefined,
        };
    }

    /**
     * Generate a QR code for receiving payments
     */
    async generateReceiveQR(options: {
        amount?: number;
        reference?: string;
    }): Promise<{ qrCode: string; shapId: string }> {
        const shapId = `shp_${this.config.merchantId}_${Math.random().toString(36).substr(2, 8)}`;

        const qrData = JSON.stringify({
            type: 'payshap_receive',
            shapId,
            merchantId: this.config.merchantId,
            amount: options.amount,
            reference: options.reference,
        });

        return {
            qrCode: Buffer.from(qrData).toString('base64'),
            shapId,
        };
    }

    /**
     * Validate a South African bank account number format
     */
    validateBankAccount(accountNumber: string, bankId: SABank): boolean {
        // Basic validation - real implementation would use bank-specific rules
        const accountLengths: Record<SABank, number[]> = {
            absa: [10, 11],
            capitec: [10],
            fnb: [10, 11, 12],
            nedbank: [10, 11],
            standard: [9, 10, 11],
            investec: [10],
            discovery: [10],
            tymebank: [10],
            african: [11],
        };

        const validLengths = accountLengths[bankId];
        if (!validLengths) return false;

        const digitsOnly = accountNumber.replace(/\D/g, '');
        return validLengths.includes(digitsOnly.length);
    }
}
