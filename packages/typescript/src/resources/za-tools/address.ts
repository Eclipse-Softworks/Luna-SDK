/**
 * Luna SDK - South African Address Utilities
 * Address validation, formatting, and postal code lookup
 */

import type { Address, SAProvince } from './types.js';

// Major SA cities with provinces
const SA_CITIES: Record<string, SAProvince> = {
    'cape town': 'WC',
    'johannesburg': 'GP',
    'pretoria': 'GP',
    'durban': 'KZN',
    'port elizabeth': 'EC',
    'gqeberha': 'EC',
    'bloemfontein': 'FS',
    'east london': 'EC',
    'polokwane': 'LP',
    'nelspruit': 'MP',
    'mbombela': 'MP',
    'kimberley': 'NC',
    'mafikeng': 'NW',
    'pietermaritzburg': 'KZN',
    'sandton': 'GP',
    'soweto': 'GP',
    'centurion': 'GP',
    'stellenbosch': 'WC',
    'paarl': 'WC',
    'george': 'WC',
    'richards bay': 'KZN',
    'umhlanga': 'KZN',
};

// Postal code ranges by province
const POSTAL_CODE_RANGES: Record<SAProvince, [number, number][]> = {
    EC: [[5000, 5999], [6000, 6499]],
    FS: [[9300, 9999]],
    GP: [[1, 299], [1000, 2199]],
    KZN: [[3000, 4499]],
    LP: [[500, 999]],
    MP: [[1100, 1399], [2200, 2499]],
    NC: [[8000, 8999]],
    NW: [[2500, 2999], [300, 499]],
    WC: [[6500, 8099]],
};

export class SAAddress {
    /**
     * Validate a South African postal code
     */
    isValidPostalCode(postalCode: string): boolean {
        const code = parseInt(postalCode, 10);
        return !isNaN(code) && code >= 1 && code <= 9999;
    }

    /**
     * Get province from postal code
     */
    getProvinceFromPostalCode(postalCode: string): SAProvince | null {
        const code = parseInt(postalCode, 10);
        if (isNaN(code)) return null;

        for (const [province, ranges] of Object.entries(POSTAL_CODE_RANGES)) {
            for (const [min, max] of ranges) {
                if (code >= min && code <= max) {
                    return province as SAProvince;
                }
            }
        }

        return null;
    }

    /**
     * Get province from city name
     */
    getProvinceFromCity(city: string): SAProvince | null {
        const normalized = city.toLowerCase().trim();
        return SA_CITIES[normalized] || null;
    }

    /**
     * Format an address for display
     */
    format(address: Address, style: 'single' | 'multi' = 'single'): string {
        const parts = [
            address.line1,
            address.line2,
            address.suburb,
            address.city,
            this.getProvinceName(address.province),
            address.postalCode,
        ].filter(Boolean);

        return style === 'multi' ? parts.join('\n') : parts.join(', ');
    }

    /**
     * Get full province name from code
     */
    getProvinceName(code: SAProvince): string {
        const names: Record<SAProvince, string> = {
            EC: 'Eastern Cape',
            FS: 'Free State',
            GP: 'Gauteng',
            KZN: 'KwaZulu-Natal',
            LP: 'Limpopo',
            MP: 'Mpumalanga',
            NC: 'Northern Cape',
            NW: 'North West',
            WC: 'Western Cape',
        };
        return names[code];
    }

    /**
     * Validate an address
     */
    validate(address: Partial<Address>): {
        valid: boolean;
        errors: string[];
    } {
        const errors: string[] = [];

        if (!address.line1 || address.line1.trim().length < 5) {
            errors.push('Address line 1 is required (min 5 characters)');
        }

        if (!address.suburb || address.suburb.trim().length < 2) {
            errors.push('Suburb is required');
        }

        if (!address.city || address.city.trim().length < 2) {
            errors.push('City is required');
        }

        if (!address.province) {
            errors.push('Province is required');
        }

        if (!address.postalCode || !this.isValidPostalCode(address.postalCode)) {
            errors.push('Valid postal code is required (4 digits)');
        }

        // Cross-validate province and postal code
        if (address.postalCode && address.province) {
            const expectedProvince = this.getProvinceFromPostalCode(address.postalCode);
            if (expectedProvince && expectedProvince !== address.province) {
                errors.push(`Postal code ${address.postalCode} does not match province ${address.province}`);
            }
        }

        return {
            valid: errors.length === 0,
            errors,
        };
    }

    /**
     * Normalize address components
     */
    normalize(address: Address): Address {
        return {
            line1: this.normalizeStreet(address.line1),
            line2: address.line2 ? this.normalizeStreet(address.line2) : undefined,
            suburb: this.titleCase(address.suburb),
            city: this.titleCase(address.city),
            province: address.province,
            postalCode: address.postalCode.padStart(4, '0'),
            country: 'ZA',
        };
    }

    /**
     * Parse a free-form address string
     */
    parse(addressString: string): Partial<Address> {
        const result: Partial<Address> = { country: 'ZA' };

        // Try to extract postal code
        const postalMatch = addressString.match(/\b(\d{4})\b/);
        if (postalMatch && postalMatch[1]) {
            result.postalCode = postalMatch[1];
            result.province = this.getProvinceFromPostalCode(postalMatch[1]) || undefined;
        }

        // Try to extract city
        for (const [city, province] of Object.entries(SA_CITIES)) {
            if (addressString.toLowerCase().includes(city)) {
                result.city = this.titleCase(city);
                result.province = result.province || province;
                break;
            }
        }

        return result;
    }

    // Private helpers

    private titleCase(str: string): string {
        return str
            .toLowerCase()
            .split(' ')
            .map(word => word.charAt(0).toUpperCase() + word.slice(1))
            .join(' ');
    }

    private normalizeStreet(street: string): string {
        return street
            .replace(/\bst\b/gi, 'Street')
            .replace(/\brd\b/gi, 'Road')
            .replace(/\bave\b/gi, 'Avenue')
            .replace(/\bdr\b/gi, 'Drive')
            .replace(/\bln\b/gi, 'Lane')
            .replace(/\bcres\b/gi, 'Crescent');
    }
}
