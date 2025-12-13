/**
 * Luna SDK - CIPC Company Lookup
 * Companies and Intellectual Property Commission integration
 */

import type { HttpClient } from '../../http/client.js';
import type { CIPCConfig, Company, CompanyType, Director } from './types.js';

export class CIPC {
    private config: CIPCConfig;
    private http: HttpClient;

    constructor(http: HttpClient, config: CIPCConfig) {
        this.http = http;
        this.config = config;
        void this.config; // Used in constructor
    }

    /**
     * Search for a company by registration number
     */
    async lookup(registrationNumber: string): Promise<Company | null> {
        const cleaned = registrationNumber.replace(/[\s/]/g, '');

        // Validate format (e.g., 2020/123456/07)
        if (!this.isValidRegistrationNumber(cleaned)) {
            throw new Error(`Invalid registration number format: ${registrationNumber}`);
        }

        try {
            const response = await this.http.request<Company>({
                method: 'GET',
                path: `/v1/za/cipc/companies/${cleaned}`,
            });
            return response.data;
        } catch {
            return null;
        }
    }

    /**
     * Search for companies by name
     */
    async searchByName(name: string, limit = 10): Promise<Company[]> {
        const response = await this.http.request<{ data: Company[] }>({
            method: 'GET',
            path: '/v1/za/cipc/companies',
            query: {
                name,
                limit: String(limit),
            },
        });
        return response.data.data;
    }

    /**
     * Verify if a company is registered and active
     */
    async verify(registrationNumber: string): Promise<{
        exists: boolean;
        isActive: boolean;
        company?: Company;
    }> {
        const company = await this.lookup(registrationNumber);

        if (!company) {
            return { exists: false, isActive: false };
        }

        return {
            exists: true,
            isActive: company.status === 'active',
            company,
        };
    }

    /**
     * Check if a company name is available
     */
    async checkNameAvailability(name: string): Promise<{
        available: boolean;
        similarNames: string[];
    }> {
        const companies = await this.searchByName(name, 5);

        const exactMatch = companies.some(
            c => c.name.toLowerCase() === name.toLowerCase()
        );

        return {
            available: !exactMatch,
            similarNames: companies.map(c => c.name),
        };
    }

    /**
     * Get directors for a company
     */
    async getDirectors(registrationNumber: string): Promise<Director[]> {
        const company = await this.lookup(registrationNumber);
        return company?.directors || [];
    }

    /**
     * Validate registration number format
     */
    isValidRegistrationNumber(regNumber: string): boolean {
        // Format: YYYY/NNNNNN/TT where TT is company type code
        const patterns = [
            /^\d{4}\/\d{6}\/\d{2}$/,   // 2020/123456/07
            /^\d{4}\d{6}\d{2}$/,       // 202012345607
            /^[A-Z]{2}\d{6}$/,         // CK123456 (old CC format)
        ];

        return patterns.some(p => p.test(regNumber));
    }

    /**
     * Parse company type from registration number
     */
    parseCompanyType(regNumber: string): CompanyType | null {
        const match = regNumber.match(/\/(\d{2})$/);
        if (!match) return null;

        const typeCode = match[1];
        if (!typeCode) return null;

        const typeMap: Record<string, CompanyType> = {
            '07': 'PTY_LTD',
            '06': 'LTD',
            '08': 'NPC',
            '23': 'CC',
            '21': 'INC',
        };

        return typeMap[typeCode] || null;
    }

    /**
     * Format registration number for display
     */
    formatRegistrationNumber(regNumber: string): string {
        const cleaned = regNumber.replace(/[\s/]/g, '');

        if (cleaned.length === 12) {
            // Format as YYYY/NNNNNN/TT
            return `${cleaned.slice(0, 4)}/${cleaned.slice(4, 10)}/${cleaned.slice(10)}`;
        }

        return regNumber;
    }
}
