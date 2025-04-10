package analyzer

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Elite-Security-Systems/radar/internal/models"
)

// DetectTechnologies identifies technologies from DNS records using signatures
func DetectTechnologies(records []models.DNSResponse, signatures models.SignatureFile) []models.DetectedTechnology {
	var detectedTechnologies []models.DetectedTechnology
	detectedMap := make(map[string]bool) // To avoid duplicates

	for _, record := range records {
		// Normalize record value by trimming any trailing periods for consistent matching
		normalizedValue := strings.TrimSuffix(record.Value, ".")
		
		// Check each signature against this record
		for _, sig := range signatures.Signatures {
			// Skip if the signature doesn't apply to this record type
			if !containsString(sig.RecordTypes, record.RecordType) && !containsString(sig.RecordTypes, "*") {
				continue
			}

			// Check each pattern in the signature
			for _, pattern := range sig.Patterns {
				re, err := regexp.Compile(pattern)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Invalid regex pattern in signature %s: %s\n", sig.Name, err)
					continue
				}

				// Try to match the original value
				if re.MatchString(record.Value) {
					// Avoid duplicates
					key := sig.Name
					if _, exists := detectedMap[key]; !exists {
						detectedMap[key] = true
						detectedTechnologies = append(detectedTechnologies, models.DetectedTechnology{
							Name:        sig.Name,
							Category:    sig.Category,
							Description: sig.Description,
							Website:     sig.Website,
							Evidence:    record.Value,
							RecordType:  record.RecordType,
						})
					}
					break // No need to check other patterns for this signature
				}
				
				// If no match with original value, try with normalized value (without trailing period)
				if normalizedValue != record.Value && re.MatchString(normalizedValue) {
					// Avoid duplicates
					key := sig.Name
					if _, exists := detectedMap[key]; !exists {
						detectedMap[key] = true
						detectedTechnologies = append(detectedTechnologies, models.DetectedTechnology{
							Name:        sig.Name,
							Category:    sig.Category,
							Description: sig.Description,
							Website:     sig.Website,
							Evidence:    record.Value,
							RecordType:  record.RecordType,
						})
					}
					break // No need to check other patterns for this signature
				}
			}
		}
	}

	return detectedTechnologies
}

// containsString checks if a string exists in a slice
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
