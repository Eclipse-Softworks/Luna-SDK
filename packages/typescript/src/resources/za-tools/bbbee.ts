/**
 * Luna SDK - B-BBEE Compliance Verification
 * Broad-Based Black Economic Empowerment compliance checking
 */

import type { HttpClient } from '../../http/client.js';
import type { BBBEEConfig, BBBEECertificate, BBBEELevel, BBBEEScorecard } from './types.js';

export class BBBEE {
    private config: BBBEEConfig;
    private http: HttpClient;

    constructor(http: HttpClient, config: BBBEEConfig) {
        this.http = http;
        this.config = config;
        void this.config; // Used in constructor
    }

    /**
     * Verify a B-BBEE certificate by number
     */
    async verifyCertificate(certificateNumber: string): Promise<BBBEECertificate | null> {
        try {
            const response = await this.http.request<BBBEECertificate>({
                method: 'GET',
                path: `/v1/za/bbbee/certificates/${certificateNumber}`,
            });
            return response.data;
        } catch {
            return null;
        }
    }

    /**
     * Look up B-BBEE status by company registration number
     */
    async lookupByCompany(registrationNumber: string): Promise<BBBEECertificate | null> {
        try {
            const response = await this.http.request<BBBEECertificate>({
                method: 'GET',
                path: '/v1/za/bbbee/lookup',
                query: { registrationNumber },
            });
            return response.data;
        } catch {
            return null;
        }
    }

    /**
     * Check if a company meets minimum B-BBEE level requirement
     */
    async meetsRequirement(
        registrationNumber: string,
        minimumLevel: BBBEELevel
    ): Promise<{
        meets: boolean;
        actualLevel?: BBBEELevel;
        certificate?: BBBEECertificate;
    }> {
        const certificate = await this.lookupByCompany(registrationNumber);

        if (!certificate || !certificate.isValid) {
            return { meets: false };
        }

        const levelNumber = typeof certificate.level === 'number' ? certificate.level : 9;
        const requiredNumber = typeof minimumLevel === 'number' ? minimumLevel : 9;

        return {
            meets: levelNumber <= requiredNumber,
            actualLevel: certificate.level,
            certificate,
        };
    }

    /**
     * Get procurement recognition percentage for a B-BBEE level
     */
    getRecognitionLevel(level: BBBEELevel): number {
        const recognitionMap: Record<string | number, number> = {
            1: 135,
            2: 125,
            3: 110,
            4: 100,
            5: 80,
            6: 60,
            7: 50,
            8: 10,
            'non-compliant': 0,
        };

        return recognitionMap[level] || 0;
    }

    /**
     * Calculate Exempted Micro Enterprise (EME) or Qualifying Small Enterprise (QSE) status
     */
    determineEnterpriseCategory(annualRevenue: number): {
        category: 'EME' | 'QSE' | 'Generic';
        automaticLevel?: BBBEELevel;
        description: string;
    } {
        if (annualRevenue <= 10_000_000) {
            // EME: R10 million or less
            return {
                category: 'EME',
                automaticLevel: 4,
                description: 'Exempted Micro Enterprise - automatic Level 4 (or Level 1 with 51%+ black ownership)',
            };
        } else if (annualRevenue <= 50_000_000) {
            // QSE: R10-50 million
            return {
                category: 'QSE',
                description: 'Qualifying Small Enterprise - simplified scorecard applies',
            };
        } else {
            // Generic: Over R50 million
            return {
                category: 'Generic',
                description: 'Generic Enterprise - full scorecard applies',
            };
        }
    }

    /**
     * Calculate scorecard points (simplified)
     */
    calculateScore(scorecard: Partial<BBBEEScorecard>): {
        total: number;
        level: BBBEELevel;
    } {
        const scores = {
            ownership: Math.min(scorecard.ownership || 0, 25),
            managementControl: Math.min(scorecard.managementControl || 0, 19),
            skillsDevelopment: Math.min(scorecard.skillsDevelopment || 0, 20),
            enterpriseSupplierDevelopment: Math.min(scorecard.enterpriseSupplierDevelopment || 0, 40),
            socioEconomicDevelopment: Math.min(scorecard.socioEconomicDevelopment || 0, 5),
        };

        const total = Object.values(scores).reduce((sum, val) => sum + val, 0);
        const level = this.pointsToLevel(total);

        return { total, level };
    }

    /**
     * Convert total points to B-BBEE level
     */
    private pointsToLevel(points: number): BBBEELevel {
        if (points >= 100) return 1;
        if (points >= 95) return 2;
        if (points >= 90) return 3;
        if (points >= 80) return 4;
        if (points >= 75) return 5;
        if (points >= 70) return 6;
        if (points >= 55) return 7;
        if (points >= 40) return 8;
        return 'non-compliant';
    }

    /**
     * Check if a certificate is still valid
     */
    isCertificateValid(certificate: BBBEECertificate): boolean {
        const expiryDate = new Date(certificate.expiryDate);
        return certificate.isValid && expiryDate > new Date();
    }

    /**
     * Get days until certificate expiry
     */
    getDaysUntilExpiry(certificate: BBBEECertificate): number {
        const expiryDate = new Date(certificate.expiryDate);
        const now = new Date();
        const diffTime = expiryDate.getTime() - now.getTime();
        return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    }
}
