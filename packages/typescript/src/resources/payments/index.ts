/**
 * Luna SDK - Payments Module
 * Unified interface for South African payment gateways
 */

export { PayFast } from './payfast.js';
export { Ozow } from './ozow.js';
export { Yoco } from './yoco.js';
export { PayShap, SA_BANKS } from './payshap.js';
export type { SABank } from './payshap.js';

export * from './types.js';

import type { HttpClient } from '../../http/client.js';
import { PayFast } from './payfast.js';
import { Ozow } from './ozow.js';
import { Yoco } from './yoco.js';
import { PayShap } from './payshap.js';
import type {
    PayFastConfig,
    OzowConfig,
    YocoConfig,
    PayShapConfig,
    Payment,
    PaymentListOptions,
    PaymentList,
} from './types.js';

export interface PaymentsConfig {
    payfast?: PayFastConfig;
    ozow?: OzowConfig;
    yoco?: YocoConfig;
    payshap?: PayShapConfig;
}

/**
 * Unified Payments API for South African payment gateways
 */
export class Payments {
    private http: HttpClient;
    private config: PaymentsConfig;

    private _payfast?: PayFast;
    private _ozow?: Ozow;
    private _yoco?: Yoco;
    private _payshap?: PayShap;

    constructor(http: HttpClient, config: PaymentsConfig) {
        this.http = http;
        this.config = config;
    }

    /**
     * PayFast gateway - Cards, EFT, SnapScan, Zapper, Mobicred
     */
    get payfast(): PayFast {
        if (!this._payfast) {
            if (!this.config.payfast) {
                throw new Error(
                    'PayFast not configured. Provide PayFast credentials when initializing LunaClient.'
                );
            }
            this._payfast = new PayFast(this.http, this.config.payfast);
        }
        return this._payfast;
    }

    /**
     * Ozow gateway - Instant EFT with all SA banks
     */
    get ozow(): Ozow {
        if (!this._ozow) {
            if (!this.config.ozow) {
                throw new Error(
                    'Ozow not configured. Provide Ozow credentials when initializing LunaClient.'
                );
            }
            this._ozow = new Ozow(this.http, this.config.ozow);
        }
        return this._ozow;
    }

    /**
     * Yoco gateway - Online payments for SMEs
     */
    get yoco(): Yoco {
        if (!this._yoco) {
            if (!this.config.yoco) {
                throw new Error(
                    'Yoco not configured. Provide Yoco credentials when initializing LunaClient.'
                );
            }
            this._yoco = new Yoco(this.http, this.config.yoco);
        }
        return this._yoco;
    }

    /**
     * PayShap gateway - Real-time account-to-account payments
     */
    get payshap(): PayShap {
        if (!this._payshap) {
            if (!this.config.payshap) {
                throw new Error(
                    'PayShap not configured. Provide PayShap credentials when initializing LunaClient.'
                );
            }
            this._payshap = new PayShap(this.http, this.config.payshap);
        }
        return this._payshap;
    }

    /**
     * List all payments across configured gateways
     */
    async list(_options?: PaymentListOptions): Promise<PaymentList> {
        const payments: Payment[] = [];

        // In a full implementation, this would aggregate from the Luna platform
        // For now, return empty as payments are stored externally

        return {
            data: payments,
            has_more: false,
        };
    }

    /**
     * Check if a specific provider is configured
     */
    isConfigured(provider: 'payfast' | 'ozow' | 'yoco' | 'payshap'): boolean {
        return !!this.config[provider];
    }

    /**
     * Get list of configured providers
     */
    getConfiguredProviders(): string[] {
        const providers: string[] = [];
        if (this.config.payfast) providers.push('payfast');
        if (this.config.ozow) providers.push('ozow');
        if (this.config.yoco) providers.push('yoco');
        if (this.config.payshap) providers.push('payshap');
        return providers;
    }
}
