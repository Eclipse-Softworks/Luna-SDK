/**
 * Luna SDK - South African ID Validation
 * Validates and parses SA ID numbers
 */

import type { SAIDInfo } from './types.js';

export class IDValidation {
    /**
     * Validate and parse a South African ID number
     * SA ID format: YYMMDD SSSS C A Z
     * - YYMMDD: Date of birth
     * - SSSS: Gender (0000-4999 = female, 5000-9999 = male)
     * - C: Citizenship (0 = SA citizen, 1 = permanent resident)
     * - A: Usually 8 (was used for race classification, now obsolete)
     * - Z: Luhn checksum digit
     */
    validate(idNumber: string): SAIDInfo {
        // Remove any spaces or dashes
        const cleaned = idNumber.replace(/[\s-]/g, '');

        // Check length
        if (cleaned.length !== 13) {
            return this.invalidResult(idNumber, 'ID number must be 13 digits');
        }

        // Check if all digits
        if (!/^\d{13}$/.test(cleaned)) {
            return this.invalidResult(idNumber, 'ID number must contain only digits');
        }

        // Extract components
        const year = parseInt(cleaned.substring(0, 2), 10);
        const month = parseInt(cleaned.substring(2, 4), 10);
        const day = parseInt(cleaned.substring(4, 6), 10);
        const genderDigits = parseInt(cleaned.substring(6, 10), 10);
        const citizenshipDigit = parseInt(cleaned.substring(10, 11), 10);

        // Validate date
        const fullYear = year >= 0 && year <= 30 ? 2000 + year : 1900 + year;
        const dateOfBirth = new Date(fullYear, month - 1, day);

        if (
            isNaN(dateOfBirth.getTime()) ||
            dateOfBirth.getMonth() !== month - 1 ||
            dateOfBirth.getDate() !== day
        ) {
            return this.invalidResult(idNumber, 'Invalid date of birth');
        }

        // Validate checksum using Luhn algorithm
        const checksumValid = this.validateLuhn(cleaned);

        return {
            idNumber: cleaned,
            isValid: checksumValid,
            dateOfBirth,
            gender: genderDigits >= 5000 ? 'male' : 'female',
            isSACitizen: citizenshipDigit === 0,
            checksumValid,
        };
    }

    /**
     * Quick validation - returns boolean
     */
    isValid(idNumber: string): boolean {
        return this.validate(idNumber).isValid;
    }

    /**
     * Extract date of birth from ID number
     */
    getDateOfBirth(idNumber: string): Date | null {
        const result = this.validate(idNumber);
        return result.isValid ? result.dateOfBirth : null;
    }

    /**
     * Get age from ID number
     */
    getAge(idNumber: string): number | null {
        const dob = this.getDateOfBirth(idNumber);
        if (!dob) return null;

        const today = new Date();
        let age = today.getFullYear() - dob.getFullYear();
        const monthDiff = today.getMonth() - dob.getMonth();

        if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < dob.getDate())) {
            age--;
        }

        return age;
    }

    /**
     * Check if person is over a specific age
     */
    isOver(idNumber: string, age: number): boolean {
        const currentAge = this.getAge(idNumber);
        return currentAge !== null && currentAge >= age;
    }

    /**
     * Generate a valid SA ID number (for testing)
     */
    generate(options?: {
        dateOfBirth?: Date;
        gender?: 'male' | 'female';
        isCitizen?: boolean;
    }): string {
        const dob = options?.dateOfBirth || new Date(1990, 0, 1);
        const year = dob.getFullYear().toString().slice(-2);
        const month = (dob.getMonth() + 1).toString().padStart(2, '0');
        const day = dob.getDate().toString().padStart(2, '0');

        // Gender sequence (random within range)
        const genderMin = options?.gender === 'female' ? 0 : 5000;
        const genderMax = options?.gender === 'female' ? 4999 : 9999;
        const genderSeq = Math.floor(Math.random() * (genderMax - genderMin + 1) + genderMin)
            .toString()
            .padStart(4, '0');

        // Citizenship
        const citizen = options?.isCitizen !== false ? '0' : '1';

        // A digit (usually 8)
        const aDigit = '8';

        // Calculate checksum
        const partial = `${year}${month}${day}${genderSeq}${citizen}${aDigit}`;
        const checksum = this.calculateLuhnChecksum(partial);

        return `${partial}${checksum}`;
    }

    // Private helpers

    private invalidResult(idNumber: string, _reason: string): SAIDInfo {
        return {
            idNumber,
            isValid: false,
            dateOfBirth: new Date(0),
            gender: 'male',
            isSACitizen: false,
            checksumValid: false,
        };
    }

    private validateLuhn(number: string): boolean {
        let sum = 0;
        let isEven = false;

        for (let i = number.length - 1; i >= 0; i--) {
            let digit = parseInt(number[i]!, 10);

            if (isEven) {
                digit *= 2;
                if (digit > 9) {
                    digit -= 9;
                }
            }

            sum += digit;
            isEven = !isEven;
        }

        return sum % 10 === 0;
    }

    private calculateLuhnChecksum(partial: string): string {
        // Try each digit 0-9 and find the one that makes the checksum valid
        for (let i = 0; i <= 9; i++) {
            if (this.validateLuhn(partial + i)) {
                return i.toString();
            }
        }
        return '0';
    }
}
