// Package main demonstrates the ResMate student housing finder using Luna SDK.
//
// Run with: go run resmate_finder.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eclipse-softworks/luna-sdk-go/luna"
)

// SearchCriteria defines search parameters for finding accommodation
type SearchCriteria struct {
	CampusName   string
	MaxBudget    int
	RequiresNSAS bool
	Gender       string // "male", "female", "mixed"
	MinRating    float64
}

func main() {
	client := luna.NewClient(
		luna.WithAPIKey(os.Getenv("LUNA_API_KEY")),
	)

	ctx := context.Background()

	// Example 1: Find budget-friendly NSFAS accommodation
	nsfasResults, err := findStudentAccommodation(ctx, client, SearchCriteria{
		RequiresNSAS: true,
		MaxBudget:    5000,
		MinRating:    3.5,
	})
	if err != nil {
		log.Printf("Search failed: %v", err)
	}

	// Example 2: Find mixed-gender housing near a specific campus
	_, err = findStudentAccommodation(ctx, client, SearchCriteria{
		CampusName: "University",
		Gender:     "mixed",
		MaxBudget:  8000,
	})
	if err != nil {
		log.Printf("Search failed: %v", err)
	}

	// Example 3: Get details for the first result
	if len(nsfasResults) > 0 {
		if err := getResidenceDetails(ctx, client, nsfasResults[0].ID); err != nil {
			log.Printf("Failed to get details: %v", err)
		}
	}

	// Example 4: Browse all available residences
	if err := browseAllResidences(ctx, client); err != nil {
		log.Printf("Browse failed: %v", err)
	}
}

// ============================================
// Student Housing Search
// ============================================

func findStudentAccommodation(ctx context.Context, client *luna.Client, criteria SearchCriteria) ([]luna.Residence, error) {
	fmt.Println("Student Housing Finder\n")
	fmt.Printf("Searching with criteria: %+v\n", criteria)
	fmt.Println(strings.Repeat("-", 40))

	// First, get available campuses
	campuses, err := client.ResMate().Campuses().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list campuses: %w", err)
	}
	fmt.Printf("\nAvailable campuses: %d\n", len(campuses.Data))

	// Find campus ID if campus name provided
	var campusID string
	if criteria.CampusName != "" {
		for _, campus := range campuses.Data {
			if strings.Contains(strings.ToLower(campus.Name), strings.ToLower(criteria.CampusName)) {
				campusID = campus.ID
				fmt.Printf("Found campus: %s\n", campus.Name)
				break
			}
		}
	}

	// Build search parameters
	searchParams := &luna.ResidenceSearch{
		Limit:     20,
		CampusID:  campusID,
		MaxPrice:  criteria.MaxBudget,
		NSFAS:     criteria.RequiresNSAS,
		Gender:    criteria.Gender,
		MinRating: criteria.MinRating,
	}

	// Search for residences
	residences, err := client.ResMate().Residences().List(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search residences: %w", err)
	}

	fmt.Printf("\nFound %d matching residences:\n\n", len(residences.Data))

	// Display results
	for _, residence := range residences.Data {
		stars := strings.Repeat("*", int(residence.Rating+0.5))
		nsfasStatus := "No"
		if residence.IsNSFASAccredited {
			nsfasStatus = "Yes"
		}

		fmt.Printf("* %s\n", residence.Name)
		fmt.Printf("   Address: %s\n", residence.Address)
		fmt.Printf("   Price: %s %d - %d/month\n", residence.CurrencyCode, residence.MinPrice, residence.MaxPrice)
		fmt.Printf("   Rating: %s (%.1f)\n", stars, residence.Rating)
		fmt.Printf("   Reviews: %d\n", residence.ReviewCount)
		fmt.Printf("   NSFAS Accredited: %s\n", nsfasStatus)
		fmt.Printf("   Gender Policy: %s\n", residence.GenderPolicy)

		// Show first 5 amenities
		amenities := residence.Amenities
		if len(amenities) > 5 {
			amenities = amenities[:5]
		}
		fmt.Printf("   Amenities: %s\n\n", strings.Join(amenities, ", "))
	}

	return residences.Data, nil
}

// ============================================
// Get Detailed Residence Information
// ============================================

func getResidenceDetails(ctx context.Context, client *luna.Client, residenceID string) error {
	fmt.Println("\nFetching detailed information...\n")

	residence, err := client.ResMate().Residences().Get(ctx, residenceID)
	if err != nil {
		return fmt.Errorf("failed to get residence: %w", err)
	}

	fmt.Printf("%s\n", residence.Name)
	fmt.Println(strings.Repeat("=", 50))

	fmt.Printf("\nLocation\n")
	fmt.Printf("   Address: %s\n", residence.Address)
	fmt.Printf("   City: %s\n", getOrDefault(residence.Location.City, "N/A"))
	fmt.Printf("   Suburb: %s\n", getOrDefault(residence.Location.Suburb, "N/A"))
	fmt.Printf("   Coordinates: %.6f, %.6f\n", residence.Location.Latitude, residence.Location.Longitude)

	fmt.Printf("\nPricing\n")
	fmt.Printf("   Range: %s %d - %d\n", residence.CurrencyCode, residence.MinPrice, residence.MaxPrice)
	fmt.Printf("   NSFAS Accredited: %v\n", residence.IsNSFASAccredited)

	fmt.Printf("\nDetails\n")
	fmt.Printf("   Gender Policy: %s\n", residence.GenderPolicy)
	fmt.Printf("   Description: %s\n", getOrDefault(residence.Description, "No description available"))

	fmt.Printf("\nReviews\n")
	fmt.Printf("   Rating: %.1f/5\n", residence.Rating)
	fmt.Printf("   Total Reviews: %d\n", residence.ReviewCount)

	fmt.Printf("\nAmenities\n")
	for _, amenity := range residence.Amenities {
		fmt.Printf("   • %s\n", amenity)
	}

	fmt.Printf("\n🖼️ Images: %d available\n", len(residence.Images))

	return nil
}

// ============================================
// Browse All Residences with Pagination
// ============================================

func browseAllResidences(ctx context.Context, client *luna.Client) error {
	fmt.Println("\n📚 Browsing all residences...\n")

	totalCount := 0
	summaries := []string{}

	// Use iterator to go through all pages
	paginator := client.ResMate().Residences().Iterate(ctx, &luna.ResidenceSearch{Limit: 10})

	for paginator.Next() {
		residence := paginator.Current()
		summary := fmt.Sprintf("%s - %s %d+ (%.1f stars)",
			residence.Name,
			residence.CurrencyCode,
			residence.MinPrice,
			residence.Rating,
		)
		summaries = append(summaries, summary)
		totalCount++

		// Limit for demo purposes
		if totalCount >= 30 {
			fmt.Println("(Showing first 30 results)")
			break
		}
	}

	if err := paginator.Err(); err != nil {
		return fmt.Errorf("pagination error: %w", err)
	}

	fmt.Printf("Total residences found: %d\n\n", totalCount)
	for i, summary := range summaries {
		fmt.Printf("%d. %s\n", i+1, summary)
	}

	return nil
}

// ============================================
// Compare Residences
// ============================================

func compareResidences(ctx context.Context, client *luna.Client, residenceIDs []string) error {
	fmt.Println("\n⚖️ Residence Comparison\n")
	fmt.Println(strings.Repeat("=", 80))

	residences := make([]*luna.Residence, 0, len(residenceIDs))
	for _, id := range residenceIDs {
		res, err := client.ResMate().Residences().Get(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get residence %s: %w", id, err)
		}
		residences = append(residences, res)
	}

	// Header
	fmt.Printf("%-20s", "Feature")
	for _, res := range residences {
		name := res.Name
		if len(name) > 15 {
			name = name[:15]
		}
		fmt.Printf("%-18s", name)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", 80))

	// Min Price
	fmt.Printf("%-20s", "Min Price")
	for _, res := range residences {
		fmt.Printf("%s %-12d", res.CurrencyCode, res.MinPrice)
	}
	fmt.Println()

	// Max Price
	fmt.Printf("%-20s", "Max Price")
	for _, res := range residences {
		fmt.Printf("%s %-12d", res.CurrencyCode, res.MaxPrice)
	}
	fmt.Println()

	// Rating
	fmt.Printf("%-20s", "Rating")
	for _, res := range residences {
		stars := strings.Repeat("*", int(res.Rating+0.5))
		fmt.Printf("%-18s", stars)
	}
	fmt.Println()

	// NSFAS
	fmt.Printf("%-20s", "NSFAS")
	for _, res := range residences {
		status := "No"
		if res.IsNSFASAccredited {
			status = "Yes"
		}
		fmt.Printf("%-18s", status)
	}
	fmt.Println()

	// Gender Policy
	fmt.Printf("%-20s", "Gender Policy")
	for _, res := range residences {
		fmt.Printf("%-18s", res.GenderPolicy)
	}
	fmt.Println()

	return nil
}

// Helper function to get value or default
func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
