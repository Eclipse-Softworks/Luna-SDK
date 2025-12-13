/**
 * Luna SDK - South African Business Tools
 * CIPC, B-BBEE, ID Validation, Address utilities
 */

// ============================================================================
// CIPC Types (Company Registration)
// ============================================================================

export interface CIPCConfig {
    apiKey: string;
    sandbox?: boolean;
}

export interface Company {
    registrationNumber: string;
    name: string;
    type: CompanyType;
    status: CompanyStatus;
    registrationDate: string;
    annualReturnsDate?: string;
    taxNumber?: string;
    directors: Director[];
    registeredAddress?: Address;
    postalAddress?: Address;
}

export type CompanyType =
    | 'PTY_LTD'      // Private Company
    | 'LTD'          // Public Company
    | 'INC'          // Incorporated Company
    | 'NPC'          // Non-Profit Company
    | 'COOP'         // Cooperative
    | 'CC'           // Close Corporation (legacy)
    | 'TRUST'        // Trust
    | 'SOLE_PROP';   // Sole Proprietorship

export type CompanyStatus = 'active' | 'deregistered' | 'in_liquidation' | 'dormant';

export interface Director {
    idNumber?: string;
    firstName: string;
    lastName: string;
    designation: 'director' | 'member' | 'trustee';
    appointmentDate?: string;
    resignationDate?: string;
    isActive: boolean;
}

// ============================================================================
// B-BBEE Types
// ============================================================================

export interface BBBEEConfig {
    apiKey: string;
    sandbox?: boolean;
}

export type BBBEELevel = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 'non-compliant';

export interface BBBEECertificate {
    certificateNumber: string;
    companyName: string;
    registrationNumber?: string;
    level: BBBEELevel;
    bbbeeRecognitionLevel: number; // 135%, 125%, 110%, etc.
    expiryDate: string;
    issuingAgency: string;
    scorecard: BBBEEScorecard;
    isValid: boolean;
}

export interface BBBEEScorecard {
    ownership: number;
    managementControl: number;
    skillsDevelopment: number;
    enterpriseSupplierDevelopment: number;
    socioEconomicDevelopment: number;
    totalScore: number;
}

// ============================================================================
// South African ID Types
// ============================================================================

export interface SAIDInfo {
    idNumber: string;
    isValid: boolean;
    dateOfBirth: Date;
    gender: 'male' | 'female';
    isSACitizen: boolean;
    /** Luhn checksum valid */
    checksumValid: boolean;
}

// ============================================================================
// Address Types
// ============================================================================

export interface Address {
    line1: string;
    line2?: string;
    suburb: string;
    city: string;
    province: SAProvince;
    postalCode: string;
    country: 'ZA';
}

export type SAProvince =
    | 'EC'   // Eastern Cape
    | 'FS'   // Free State
    | 'GP'   // Gauteng
    | 'KZN'  // KwaZulu-Natal
    | 'LP'   // Limpopo
    | 'MP'   // Mpumalanga
    | 'NC'   // Northern Cape
    | 'NW'   // North West
    | 'WC';  // Western Cape

export const SA_PROVINCES = {
    EC: 'Eastern Cape',
    FS: 'Free State',
    GP: 'Gauteng',
    KZN: 'KwaZulu-Natal',
    LP: 'Limpopo',
    MP: 'Mpumalanga',
    NC: 'Northern Cape',
    NW: 'North West',
    WC: 'Western Cape',
} as const;

// ============================================================================
// Bank Account Types
// ============================================================================

export interface BankAccount {
    bankName: string;
    bankCode: string;
    accountNumber: string;
    accountType: 'cheque' | 'savings' | 'transmission' | 'bond';
    branchCode: string;
}

export const SA_BANKS = {
    ABSA: { name: 'ABSA Bank', code: '632005' },
    CAPITEC: { name: 'Capitec Bank', code: '470010' },
    FNB: { name: 'First National Bank', code: '250655' },
    NEDBANK: { name: 'Nedbank', code: '198765' },
    STANDARD: { name: 'Standard Bank', code: '051001' },
    INVESTEC: { name: 'Investec Bank', code: '580105' },
    DISCOVERY: { name: 'Discovery Bank', code: '679000' },
    TYMEBANK: { name: 'TymeBank', code: '678910' },
    AFRICAN: { name: 'African Bank', code: '430000' },
} as const;
