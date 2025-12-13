/**
 * Luna SDK - South African Business Tools Module
 * CIPC, B-BBEE, ID Validation, Address utilities
 */

export { CIPC } from './cipc.js';
export { BBBEE } from './bbbee.js';
export { IDValidation } from './id-validation.js';
export { SAAddress } from './address.js';

export * from './types.js';

import type { HttpClient } from '../../http/client.js';
import { CIPC } from './cipc.js';
import { BBBEE } from './bbbee.js';
import { IDValidation } from './id-validation.js';
import { SAAddress } from './address.js';
import type { CIPCConfig, BBBEEConfig } from './types.js';

export interface ZAToolsConfig {
    cipc?: CIPCConfig;
    bbbee?: BBBEEConfig;
}

/**
 * South African Business Tools
 */
export class ZATools {
    private http: HttpClient;
    private config: ZAToolsConfig;

    private _cipc?: CIPC;
    private _bbbee?: BBBEE;
    private _idValidation?: IDValidation;
    private _address?: SAAddress;

    constructor(http: HttpClient, config: ZAToolsConfig = {}) {
        this.http = http;
        this.config = config;
    }

    /**
     * CIPC company lookup and verification
     */
    get cipc(): CIPC {
        if (!this._cipc) {
            if (!this.config.cipc) {
                throw new Error(
                    'CIPC not configured. Provide CIPC credentials when initializing LunaClient.'
                );
            }
            this._cipc = new CIPC(this.http, this.config.cipc);
        }
        return this._cipc;
    }

    /**
     * B-BBEE compliance verification
     */
    get bbbee(): BBBEE {
        if (!this._bbbee) {
            if (!this.config.bbbee) {
                throw new Error(
                    'B-BBEE not configured. Provide B-BBEE credentials when initializing LunaClient.'
                );
            }
            this._bbbee = new BBBEE(this.http, this.config.bbbee);
        }
        return this._bbbee;
    }

    /**
     * SA ID number validation (no config needed)
     */
    get idValidation(): IDValidation {
        if (!this._idValidation) {
            this._idValidation = new IDValidation();
        }
        return this._idValidation;
    }

    /**
     * SA address utilities (no config needed)
     */
    get address(): SAAddress {
        if (!this._address) {
            this._address = new SAAddress();
        }
        return this._address;
    }

    /**
     * Quick ID validation
     */
    validateID(idNumber: string) {
        return this.idValidation.validate(idNumber);
    }

    /**
     * Quick address validation
     */
    validateAddress(address: Parameters<SAAddress['validate']>[0]) {
        return this.address.validate(address);
    }

    /**
     * Check if services are configured
     */
    isConfigured(service: 'cipc' | 'bbbee'): boolean {
        return !!this.config[service];
    }
}
