"""
Luna SDK - ResMate Student Housing Finder

Build a student accommodation search application using Luna's ResMate API.
This example demonstrates searching and filtering student residences.
"""

import asyncio
import os
from dataclasses import dataclass
from typing import Optional

from luna import LunaClient


client = LunaClient(api_key=os.environ["LUNA_API_KEY"])


@dataclass
class SearchCriteria:
    """Search criteria for finding accommodation."""
    campus_name: Optional[str] = None
    max_budget: Optional[int] = None
    requires_nsfas: Optional[bool] = None
    gender: Optional[str] = None  # 'male', 'female', 'mixed'
    min_rating: Optional[float] = None


# ============================================
# Student Housing Search Application
# ============================================

async def find_student_accommodation(criteria: SearchCriteria):
    """Find student accommodation based on search criteria."""
    print("Student Housing Finder\n")
    print(f"Searching with criteria: {criteria}")
    print("-" * 40)

    # First, get available campuses
    campuses = await client.resmate.campuses.list()
    print(f"\nAvailable campuses: {len(campuses.data)}")

    # Find campus ID if campus name provided
    campus_id = None
    if criteria.campus_name:
        for campus in campuses.data:
            if criteria.campus_name.lower() in campus.name.lower():
                campus_id = campus.id
                print(f"Found campus: {campus.name}")
                break

    # Search for residences with filters
    residences = await client.resmate.residences.list(
        limit=20,
        campus_id=campus_id,
        max_price=criteria.max_budget,
        nsfas=criteria.requires_nsfas,
        gender=criteria.gender,
        min_rating=criteria.min_rating,
    )

    print(f"\nFound {len(residences.data)} matching residences:\n")

    # Display results
    for residence in residences.data:
        stars = "*" * round(residence.rating)
        nsfas_status = "Yes" if residence.is_nsfas_accredited else "No"
        
        print(f"* {residence.name}")
        print(f"   Address: {residence.address}")
        print(f"   Price: {residence.currency_code} {residence.min_price} - {residence.max_price}/month")
        print(f"   Rating: {stars} ({residence.rating})")
        print(f"   Reviews: {residence.review_count}")
        print(f"   NSFAS Accredited: {nsfas_status}")
        print(f"   Gender Policy: {residence.gender_policy}")
        print(f"   Amenities: {', '.join(residence.amenities[:5])}")
        print()

    return residences.data


# ============================================
# Get Detailed Residence Information
# ============================================

async def get_residence_details(residence_id: str):
    """Get detailed information about a specific residence."""
    print("\nFetching detailed information...\n")

    residence = await client.resmate.residences.get(residence_id)

    print(f"{residence.name}")
    print("=" * 50)
    
    print(f"\nLocation")
    print(f"   Address: {residence.address}")
    print(f"   City: {residence.location.get('city', 'N/A')}")
    print(f"   Suburb: {residence.location.get('suburb', 'N/A')}")
    print(f"   Coordinates: {residence.location['latitude']}, {residence.location['longitude']}")

    print(f"\nPricing")
    print(f"   Range: {residence.currency_code} {residence.min_price} - {residence.max_price}")
    print(f"   NSFAS Accredited: {'Yes' if residence.is_nsfas_accredited else 'No'}")

    print(f"\nDetails")
    print(f"   Gender Policy: {residence.gender_policy}")
    print(f"   Description: {residence.description or 'No description available'}")

    print(f"\nReviews")
    print(f"   Rating: {residence.rating}/5")
    print(f"   Total Reviews: {residence.review_count}")

    print(f"\nAmenities")
    for amenity in residence.amenities:
        print(f"   • {amenity}")

    print(f"\n🖼️ Images: {len(residence.images)} available")

    return residence


# ============================================
# Browse All Residences with Pagination
# ============================================

async def browse_all_residences():
    """Browse all residences using automatic pagination."""
    print("\n📚 Browsing all residences...\n")

    total_count = 0
    summaries = []

    # Use async iterator to go through all pages
    async for residence in client.resmate.residences.iterate(limit=10):
        summaries.append(
            f"{residence.name} - {residence.currency_code} {residence.min_price}+ ({residence.rating} stars)"
        )
        total_count += 1

        # Limit for demo purposes
        if total_count >= 30:
            print("(Showing first 30 results)")
            break

    print(f"Total residences found: {total_count}\n")
    for i, summary in enumerate(summaries, 1):
        print(f"{i}. {summary}")


# ============================================
# Compare Residences
# ============================================

async def compare_residences(residence_ids: list[str]):
    """Compare multiple residences side by side."""
    print("\n⚖️ Residence Comparison\n")
    print("=" * 80)

    residences = []
    for rid in residence_ids:
        res = await client.resmate.residences.get(rid)
        residences.append(res)

    # Header
    print(f"{'Feature':<20}", end="")
    for res in residences:
        print(f"{res.name[:15]:<18}", end="")
    print()
    print("-" * 80)

    # Price range
    print(f"{'Min Price':<20}", end="")
    for res in residences:
        print(f"{res.currency_code} {res.min_price:<12}", end="")
    print()

    print(f"{'Max Price':<20}", end="")
    for res in residences:
        print(f"{res.currency_code} {res.max_price:<12}", end="")
    print()

    # Rating
    print(f"{'Rating':<20}", end="")
    for res in residences:
        print(f"{'*' * round(res.rating):<18}", end="")
    print()

    # NSFAS
    print(f"{'NSFAS':<20}", end="")
    for res in residences:
        status = "Yes" if res.is_nsfas_accredited else "No"
        print(f"{status:<18}", end="")
    print()

    # Gender
    print(f"{'Gender Policy':<20}", end="")
    for res in residences:
        print(f"{res.gender_policy:<18}", end="")
    print()

    # Amenities count
    print(f"{'Amenities':<20}", end="")
    for res in residences:
        print(f"{len(res.amenities)} available{'':<6}", end="")
    print()


# ============================================
# Main Application
# ============================================

async def main():
    """Run the student housing finder."""
    try:
        # Example 1: Find budget-friendly NSFAS accommodation
        nsfas_results = await find_student_accommodation(
            SearchCriteria(
                requires_nsfas=True,
                max_budget=5000,
                min_rating=3.5,
            )
        )

        # Example 2: Find mixed-gender housing near a specific campus
        await find_student_accommodation(
            SearchCriteria(
                campus_name="University",
                gender="mixed",
                max_budget=8000,
            )
        )

        # Example 3: Get details for the first result
        if nsfas_results:
            await get_residence_details(nsfas_results[0].id)

        # Example 4: Browse all available residences
        await browse_all_residences()

        # Example 5: Compare residences (if we have multiple results)
        if len(nsfas_results) >= 2:
            await compare_residences([r.id for r in nsfas_results[:3]])

    except Exception as e:
        print(f"Application error: {e}")
        raise


if __name__ == "__main__":
    asyncio.run(main())
