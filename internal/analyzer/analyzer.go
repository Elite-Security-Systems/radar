package analyzer

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Elite-Security-Systems/radar/internal/dns"
	"github.com/Elite-Security-Systems/radar/internal/models"
)

// Config contains the configuration for the analyzer
type Config struct {
	Domain         string
	Timeout        time.Duration
	Debug          bool
	MaxRecords     int
	IncludeRecords bool
}

// AnalyzeDomain performs a complete analysis of a domain
func AnalyzeDomain(config Config, signatures models.SignatureFile) (*models.Result, error) {
	// Normalize domain
	domain := config.Domain
	if !strings.HasSuffix(domain, ".") {
		domain = domain + "."
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Query all DNS records
	dnsClient := dns.NewClient(config.Debug)
	allRecords, err := dnsClient.QueryAllRecords(ctx, domain, config.Timeout/2, config.MaxRecords)
	
	// Continue with partial results even if we hit timeout
	if err != nil && err != context.DeadlineExceeded {
		return nil, fmt.Errorf("error querying DNS records: %w", err)
	}

	if err == context.DeadlineExceeded && config.Debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Query timeout reached, proceeding with collected records\n")
	}

	// Detect technologies from the records
	detectedTechnologies := DetectTechnologies(allRecords, signatures)

	// Prepare result
	result := &models.Result{
		Domain:               strings.TrimSuffix(domain, "."),
		DetectedTechnologies: detectedTechnologies,
	}

	// Include all records if requested
	if config.IncludeRecords {
		result.AllRecords = allRecords
	}

	return result, nil
}

// aggregateResults collects and deduplicates DNS records from all resolvers
func aggregateResults(responsesMap map[string]models.DNSResponse) []models.DNSResponse {
	var responses []models.DNSResponse
	
	// Mutex to protect concurrent access to the responses slice
	var mutex sync.Mutex
	
	// Create a map to track seen records and avoid duplicates
	seen := make(map[string]bool)
	
	for _, response := range responsesMap {
		// Create a unique key based on the record type and value
		key := fmt.Sprintf("%s-%s-%s", response.Domain, response.RecordType, response.Value)
		
		mutex.Lock()
		if !seen[key] {
			seen[key] = true
			responses = append(responses, response)
		}
		mutex.Unlock()
	}
	
	return responses
}
