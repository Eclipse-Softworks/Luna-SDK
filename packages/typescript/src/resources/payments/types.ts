/**
 * Luna SDK - Payments Module Types
 * South African payment gateway integrations
 */

// ============================================================================
// Common Types
// ============================================================================

export type Currency = 'ZAR' | 'USD' | 'EUR' | 'GBP';

export type PaymentStatus =
    | 'pending'
    | 'processing'
    | 'completed'
    | 'failed'
    | 'cancelled'
    | 'refunded';

export interface Amount {
    /** Amount in cents (e.g., 10000 = R100.00) */
    value: number;
    /** ISO 4217 currency code */
    currency: Currency;
}

export interface PaymentMetadata {
    [key: string]: string | number | boolean;
}

export interface BasePayment {
    id: string;
    amount: Amount;
    status: PaymentStatus;
    reference: string;
    description?: string;
    metadata?: PaymentMetadata;
    created_at: string;
    updated_at: string;
}

// ============================================================================
// PayFast Types
// ============================================================================

export interface PayFastConfig {
    merchantId: string;
    merchantKey: string;
    passphrase?: string;
    sandbox?: boolean;
}

export interface PayFastPaymentRequest {
    amount: number;
    currency?: Currency;
    itemName: string;
    itemDescription?: string;
    emailAddress?: string;
    cellNumber?: string;
    returnUrl: string;
    cancelUrl: string;
    notifyUrl: string;
    customStr1?: string;
    customStr2?: string;
    customStr3?: string;
    customInt1?: number;
    customInt2?: number;
    paymentMethod?: 'cc' | 'dc' | 'eft' | 'mp' | 'sc' | 'ss' | 'zp' | 'mt';
}

export interface PayFastPayment extends BasePayment {
    provider: 'payfast';
    paymentUrl: string;
    pfPaymentId?: string;
    signature: string;
}

export interface PayFastWebhookPayload {
    m_payment_id: string;
    pf_payment_id: string;
    payment_status: string;
    item_name: string;
    amount_gross: string;
    amount_fee: string;
    amount_net: string;
    signature: string;
    [key: string]: string;
}

// ============================================================================
// Ozow Types
// ============================================================================

export interface OzowConfig {
    siteCode: string;
    privateKey: string;
    apiKey: string;
    sandbox?: boolean;
}

export interface OzowPaymentRequest {
    amount: number;
    transactionReference: string;
    bankReference: string;
    cancelUrl: string;
    errorUrl: string;
    successUrl: string;
    notifyUrl: string;
    isTest?: boolean;
    customer?: {
        firstName?: string;
        lastName?: string;
        email?: string;
        phone?: string;
    };
}

export interface OzowPayment extends BasePayment {
    provider: 'ozow';
    paymentUrl: string;
    transactionId?: string;
    bankName?: string;
}

export interface OzowWebhookPayload {
    TransactionId: string;
    TransactionReference: string;
    Amount: string;
    Status: string;
    StatusMessage: string;
    Hash: string;
    [key: string]: string;
}

// ============================================================================
// Yoco Types
// ============================================================================

export interface YocoConfig {
    secretKey: string;
    publicKey?: string;
    sandbox?: boolean;
}

export interface YocoPaymentRequest {
    amount: number;
    currency?: Currency;
    successUrl: string;
    cancelUrl: string;
    failureUrl?: string;
    metadata?: PaymentMetadata;
    lineItems?: YocoLineItem[];
}

export interface YocoLineItem {
    displayName: string;
    quantity: number;
    pricingDetails: {
        price: number;
    };
}

export interface YocoPayment extends BasePayment {
    provider: 'yoco';
    checkoutId: string;
    redirectUrl: string;
}

export interface YocoWebhookPayload {
    id: string;
    type: string;
    payload: {
        id: string;
        status: string;
        amount: number;
        currency: string;
        metadata: PaymentMetadata;
    };
}

// ============================================================================
// PayShap / Capitec Pay Types
// ============================================================================

export interface PayShapConfig {
    apiKey: string;
    merchantId: string;
    sandbox?: boolean;
}

export interface PayShapPaymentRequest {
    amount: number;
    payerShapId?: string;
    payerAccountNumber?: string;
    payerBankId?: string;
    reference: string;
    expiryMinutes?: number;
}

export interface PayShapPayment extends BasePayment {
    provider: 'payshap';
    shapId: string;
    qrCode?: string;
    expiresAt: string;
}

// ============================================================================
// Unified Payment Types
// ============================================================================

export type PaymentProvider = 'payfast' | 'ozow' | 'yoco' | 'payshap';

export type Payment = PayFastPayment | OzowPayment | YocoPayment | PayShapPayment;

export interface PaymentListOptions {
    limit?: number;
    cursor?: string;
    status?: PaymentStatus;
    provider?: PaymentProvider;
    from?: string;
    to?: string;
}

export interface PaymentList {
    data: Payment[];
    has_more: boolean;
    next_cursor?: string;
}

export interface RefundRequest {
    paymentId: string;
    amount?: number;
    reason?: string;
}

export interface Refund {
    id: string;
    paymentId: string;
    amount: Amount;
    status: 'pending' | 'completed' | 'failed';
    reason?: string;
    created_at: string;
}
