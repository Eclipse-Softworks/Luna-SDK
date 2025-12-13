/**
 * Luna SDK - ResMate Student Housing Finder
 *
 * Build a student accommodation search application using Luna's ResMate API.
 * This example demonstrates searching and filtering student residences.
 */

import { LunaClient } from '@eclipse-softworks/luna-sdk';

const client = new LunaClient({
    apiKey: process.env.LUNA_API_KEY!,
});

// ============================================
// Student Housing Search Application
// ============================================

interface SearchCriteria {
    campusName?: string;
    maxBudget?: number;
    requiresNsfas?: boolean;
    gender?: 'male' | 'female' | 'mixed';
    minRating?: number;
}

async function findStudentAccommodation(criteria: SearchCriteria) {
    console.log('Student Housing Finder\n');
    console.log('Searching with criteria:', criteria);
    console.log('---');

    // First, get available campuses
    const campuses = await client.resMate.campuses.list();
    console.log(`\nAvailable campuses: ${campuses.data.length}`);

    // Find campus ID if campus name provided
    let campusId: string | undefined;
    if (criteria.campusName) {
        const campus = campuses.data.find((c) =>
            c.name.toLowerCase().includes(criteria.campusName!.toLowerCase())
        );
        if (campus) {
            campusId = campus.id;
            console.log(`Found campus: ${campus.name}`);
        }
    }

    // Search for residences with filters
    const residences = await client.resMate.residences.list({
        limit: 20,
        campus_id: campusId,
        max_price: criteria.maxBudget,
        nsfas: criteria.requiresNsfas,
        gender: criteria.gender,
        min_rating: criteria.minRating,
    });

    console.log(`\nFound ${residences.data.length} matching residences:\n`);

    // Display results
    for (const residence of residences.data) {
        console.log(`* ${residence.name}`);
        console.log(`   Address: ${residence.address}`);
        console.log(
            `   Price: ${residence.currency_code} ${residence.min_price} - ${residence.max_price}/month`
        );
        console.log(`   Rating: ${'*'.repeat(Math.round(residence.rating))} (${residence.rating})`);
        console.log(`   Reviews: ${residence.review_count}`);
        console.log(`   NSFAS Accredited: ${residence.is_nsfas_accredited ? 'Yes' : 'No'}`);
        console.log(`   Gender Policy: ${residence.gender_policy}`);
        console.log(`   Amenities: ${residence.amenities.slice(0, 5).join(', ')}`);
        console.log('');
    }

    return residences.data;
}

// ============================================
// Get Detailed Residence Information
// ============================================

async function getResidenceDetails(residenceId: string) {
    console.log('\nFetching detailed information...\n');

    const residence = await client.resMate.residences.get(residenceId);

    console.log(`${residence.name}`);
    console.log('='.repeat(50));
    console.log(`\nLocation`);
    console.log(`   Address: ${residence.address}`);
    console.log(`   City: ${residence.location.city || 'N/A'}`);
    console.log(`   Suburb: ${residence.location.suburb || 'N/A'}`);
    console.log(
        `   Coordinates: ${residence.location.latitude}, ${residence.location.longitude}`
    );

    console.log(`\nPricing`);
    console.log(
        `   Range: ${residence.currency_code} ${residence.min_price} - ${residence.max_price}`
    );
    console.log(`   NSFAS Accredited: ${residence.is_nsfas_accredited ? 'Yes' : 'No'}`);

    console.log(`\nDetails`);
    console.log(`   Gender Policy: ${residence.gender_policy}`);
    console.log(`   Description: ${residence.description || 'No description available'}`);

    console.log(`\nReviews`);
    console.log(`   Rating: ${residence.rating}/5`);
    console.log(`   Total Reviews: ${residence.review_count}`);

    console.log(`\nAmenities`);
    residence.amenities.forEach((amenity) => {
        console.log(`   - ${amenity}`);
    });

    console.log(`\nImages: ${residence.images.length} available`);

    return residence;
}

// ============================================
// Browse All Residences with Pagination
// ============================================

async function browseAllResidences() {
    console.log('\nBrowsing all residences...\n');

    let totalCount = 0;
    const summaries: string[] = [];

    // Use iterator to go through all pages
    for await (const residence of client.resMate.residences.iterate({ limit: 10 })) {
        summaries.push(
            `${residence.name} - ${residence.currency_code} ${residence.min_price}+ (${residence.rating} stars)`
        );
        totalCount++;

        // Limit for demo purposes
        if (totalCount >= 30) {
            console.log('(Showing first 30 results)');
            break;
        }
    }

    console.log(`Total residences found: ${totalCount}\n`);
    summaries.forEach((summary, index) => {
        console.log(`${index + 1}. ${summary}`);
    });
}

// ============================================
// Main Application
// ============================================

async function main() {
    try {
        // Example 1: Find budget-friendly NSFAS accommodation
        const nsfasResults = await findStudentAccommodation({
            requiresNsfas: true,
            maxBudget: 5000,
            minRating: 3.5,
        });

        // Example 2: Find mixed-gender housing near a specific campus
        await findStudentAccommodation({
            campusName: 'University',
            gender: 'mixed',
            maxBudget: 8000,
        });

        // Example 3: Get details for the first result
        if (nsfasResults.length > 0) {
            await getResidenceDetails(nsfasResults[0].id);
        }

        // Example 4: Browse all available residences
        await browseAllResidences();
    } catch (error) {
        console.error('Application error:', error);
        process.exit(1);
    }
}

main();
