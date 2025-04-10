package dns

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Elite-Security-Systems/radar/internal/models"
	"github.com/miekg/dns"
)

// Client represents a DNS client for querying records
type Client struct {
	debug         bool
	resolvers     []string
	recordCounter int
	responsesMap  map[string]models.DNSResponse
	mutex         sync.Mutex
}

// NewClient creates a new DNS client
func NewClient(debug bool) *Client {
	return &Client{
		debug:     debug,
		resolvers: getResolvers(),
		responsesMap: make(map[string]models.DNSResponse),
	}
}

// getResolvers returns the list of DNS resolvers to use
func getResolvers() []string {
	return []string{
		"8.8.8.8:53",       // Google DNS
		"1.1.1.1:53",       // Cloudflare DNS
		"9.9.9.9:53",       // Quad9
		"208.67.222.222:53", // OpenDNS
	}
}

// QueryAllRecords queries all DNS record types for a domain
func (c *Client) QueryAllRecords(ctx context.Context, domain string, queryTimeout time.Duration, maxRecords int) ([]models.DNSResponse, error) {
	c.recordCounter = 0
	c.responsesMap = make(map[string]models.DNSResponse)
	
	// Create a channel to signal completion
	done := make(chan struct{})
	defer close(done)

	// Create a wait group to track all queries
	var wg sync.WaitGroup

	// Use a separate context for the actual queries with a shorter timeout
	queryCtx, queryCancel := context.WithTimeout(ctx, queryTimeout)
	defer queryCancel()

	// Query system resolver for TXT records
	wg.Add(1)
	go c.querySystemResolver(queryCtx, &wg, domain, maxRecords)

	// Query miekg/dns resolvers in parallel
	for _, resolver := range c.resolvers {
		// Priority records
		wg.Add(1)
		go c.queryPriorityRecords(queryCtx, &wg, domain, resolver, maxRecords)

		// Secondary records
		wg.Add(1)
		go c.querySecondaryRecords(queryCtx, &wg, domain, resolver, maxRecords)
	}

	// Wait for all queries to complete or timeout
	waitDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitDone)
	}()

	// Wait for either completion or context cancellation
	select {
	case <-waitDone:
		// All done successfully
		if c.debug {
			fmt.Printf("[DEBUG] All queries completed successfully\n")
		}
	case <-ctx.Done():
		// Context timed out
		if c.debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Query timeout reached, returning partial results\n")
		}
		return c.convertMapToSlice(), ctx.Err()
	}

	// Final count of records for debugging
	if c.debug {
		// Count records by type
		recordCounts := make(map[string]int)
		for _, record := range c.convertMapToSlice() {
			recordCounts[record.RecordType]++
		}

		fmt.Printf("[DEBUG] Record counts by type:\n")
		for recordType, count := range recordCounts {
			fmt.Printf("[DEBUG]   %s: %d\n", recordType, count)
		}

		fmt.Printf("[DEBUG] Total records collected: %d\n", len(c.responsesMap))
	}

	return c.convertMapToSlice(), nil
}

// querySystemResolver queries the system DNS resolver for TXT records
func (c *Client) querySystemResolver(ctx context.Context, wg *sync.WaitGroup, domain string, maxRecords int) {
	defer wg.Done()

	if c.debug {
		fmt.Println("[DEBUG] Making system resolver query for TXT records...")
	}

	// Create a channel for DNS lookup results
	resultChan := make(chan struct {
		txt string
		err error
	})

	// Perform lookup in a goroutine
	go func() {
		txtRecords, err := net.LookupTXT(strings.TrimSuffix(domain, "."))
		if err != nil {
			resultChan <- struct {
				txt string
				err error
			}{"", err}
			return
		}

		// Send each record individually to avoid blocking
		for _, txt := range txtRecords {
			select {
			case resultChan <- struct {
				txt string
				err error
			}{txt, nil}:
			case <-ctx.Done():
				return
			}
		}

		// Signal end of records
		close(resultChan)
	}()

	// Process results with timeout
	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				// Channel closed, all records processed
				return
			}

			if result.err != nil {
				if c.debug {
					fmt.Printf("[DEBUG] Error getting TXT records with system resolver: %v\n", result.err)
				}
				return
			}

			recordKey := fmt.Sprintf("%s-TXT-%s", domain, result.txt)

			c.mutex.Lock()
			// Check if we've reached the max record limit
			if c.recordCounter >= maxRecords {
				if c.debug {
					fmt.Printf("[DEBUG] Maximum record limit (%d) reached, stopping collection\n", maxRecords)
				}
				c.mutex.Unlock()
				return
			}

			if _, exists := c.responsesMap[recordKey]; !exists {
				c.responsesMap[recordKey] = models.DNSResponse{
					Domain:     domain,
					RecordType: "TXT",
					TTL:        300, // Default TTL
					Value:      result.txt,
				}
				c.recordCounter++

				if c.debug {
					fmt.Printf("[DEBUG] Found TXT via system resolver: %s\n", result.txt)
				}
			}
			c.mutex.Unlock()

		case <-ctx.Done():
			if c.debug {
				fmt.Printf("[DEBUG] TXT record lookup timed out\n")
			}
			return
		}
	}
}

// queryPriorityRecords queries the most important record types
func (c *Client) queryPriorityRecords(ctx context.Context, wg *sync.WaitGroup, domain string, resolver string, maxRecords int) {
	defer wg.Done()

	// Important record types to always check first
	priorityTypes := []uint16{
		dns.TypeA,
		dns.TypeAAAA,
		dns.TypeCNAME,
		dns.TypeMX,
		dns.TypeTXT,
		dns.TypeNS,
		dns.TypeSOA,
		dns.TypeSRV,
		dns.TypeCAA,
	}

	// Create a new DNS client for this resolver
	client := &dns.Client{
		Timeout: 3 * time.Second, // Shorter per-query timeout
	}

	if c.debug {
		fmt.Printf("[DEBUG] Querying priority records from %s\n", resolver)
	}

	// Query each priority record type
	for _, typeCode := range priorityTypes {
		// Check if we should continue or stop
		select {
		case <-ctx.Done():
			return
		default:
			// Continue processing
		}

		c.mutex.Lock()
		// Check if we've reached the max record limit
		if c.recordCounter >= maxRecords {
			c.mutex.Unlock()
			return
		}
		c.mutex.Unlock()

		typeName := RecordTypeToString(typeCode)

		// Create a new DNS message
		msg := new(dns.Msg)
		msg.SetQuestion(domain, typeCode)
		msg.RecursionDesired = true

		// Make the query
		resp, _, err := client.Exchange(msg, resolver)

		if c.debug {
			if err != nil {
				fmt.Printf("[DEBUG] Error querying %s records from %s: %v\n", typeName, resolver, err)
			} else if resp != nil {
				fmt.Printf("[DEBUG] %s response from %s - Rcode: %d, Answer records: %d\n",
					typeName, resolver, resp.Rcode, len(resp.Answer))
			}
		}

		if err != nil || resp == nil || resp.Rcode != dns.RcodeSuccess {
			continue
		}

		// Process the answer section
		for _, rr := range resp.Answer {
			value := ExtractValue(rr)
			if value == "" {
				continue
			}

			// Use record type, name and value as a unique key
			recordKey := fmt.Sprintf("%s-%s-%s", domain, typeName, value)

			c.mutex.Lock()
			// Check if we've reached the max record limit
			if c.recordCounter >= maxRecords {
				c.mutex.Unlock()
				return
			}

			// Only add if we haven't seen this exact record before
			if _, exists := c.responsesMap[recordKey]; !exists {
				c.responsesMap[recordKey] = models.DNSResponse{
					Domain:     domain,
					RecordType: typeName,
					TTL:        rr.Header().Ttl,
					Value:      value,
				}
				c.recordCounter++

				if c.debug {
					fmt.Printf("[DEBUG] Found %s record via %s: %s\n", typeName, resolver, value)
				}
			}
			c.mutex.Unlock()
		}
	}
}

// querySecondaryRecords queries all other record types
func (c *Client) querySecondaryRecords(ctx context.Context, wg *sync.WaitGroup, domain string, resolver string, maxRecords int) {
	defer wg.Done()

	// Skip these record types that are less likely to provide useful information and may cause issues
	skipTypes := map[uint16]bool{
		dns.TypeNULL: true,
		dns.TypeOPT:  true,
	}

	// Important record types (already queried in queryPriorityRecords)
	priorityTypes := map[uint16]bool{
		dns.TypeA:     true,
		dns.TypeAAAA:  true,
		dns.TypeCNAME: true,
		dns.TypeMX:    true,
		dns.TypeTXT:   true,
		dns.TypeNS:    true,
		dns.TypeSOA:   true,
		dns.TypeSRV:   true,
		dns.TypeCAA:   true,
	}

	// Create a new DNS client for this resolver
	client := &dns.Client{
		Timeout: 2 * time.Second, // Shorter per-query timeout for less important records
	}

	if c.debug {
		fmt.Printf("[DEBUG] Querying secondary records from %s\n", resolver)
	}

	// Query other record types
	for typeCode, typeName := range GetRecordTypeMapping() {
		// Skip if this is a priority type or a type to skip
		if priorityTypes[typeCode] || skipTypes[typeCode] {
			continue
		}

		// Check if we should continue or stop
		select {
		case <-ctx.Done():
			return
		default:
			// Continue processing
		}

		c.mutex.Lock()
		// Check if we've reached the max record limit
		if c.recordCounter >= maxRecords {
			c.mutex.Unlock()
			return
		}
		c.mutex.Unlock()

		// Create a new DNS message
		msg := new(dns.Msg)
		msg.SetQuestion(domain, typeCode)
		msg.RecursionDesired = true

		// Make the query
		resp, _, err := client.Exchange(msg, resolver)

		if c.debug {
			if err != nil {
				fmt.Printf("[DEBUG] Error querying %s records from %s: %v\n", typeName, resolver, err)
			} else if resp != nil && len(resp.Answer) > 0 {
				fmt.Printf("[DEBUG] %s response from %s - Rcode: %d, Answer records: %d\n",
					typeName, resolver, resp.Rcode, len(resp.Answer))
			}
		}

		if err != nil || resp == nil || resp.Rcode != dns.RcodeSuccess {
			continue
		}

		// Process the answer section
		for _, rr := range resp.Answer {
			value := ExtractValue(rr)
			if value == "" {
				continue
			}

			recordKey := fmt.Sprintf("%s-%s-%s", domain, typeName, value)

			c.mutex.Lock()
			// Check if we've reached the max record limit
			if c.recordCounter >= maxRecords {
				c.mutex.Unlock()
				return
			}

			// Only add if we haven't seen this exact record before
			if _, exists := c.responsesMap[recordKey]; !exists {
				c.responsesMap[recordKey] = models.DNSResponse{
					Domain:     domain,
					RecordType: typeName,
					TTL:        rr.Header().Ttl,
					Value:      value,
				}
				c.recordCounter++

				if c.debug {
					fmt.Printf("[DEBUG] Found %s record via %s: %s\n", typeName, resolver, value)
				}
			}
			c.mutex.Unlock()
		}
	}
}

// convertMapToSlice converts the internal map to a slice of DNSResponses
func (c *Client) convertMapToSlice() []models.DNSResponse {
	var allResponses []models.DNSResponse
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	for _, response := range c.responsesMap {
		allResponses = append(allResponses, response)
	}
	
	return allResponses
}
