package analyzer

import (
	"reflect"
	"testing"

	"github.com/Elite-Security-Systems/radar/internal/models"
)

func TestDetectTechnologies(t *testing.T) {
	testCases := []struct {
		name       string
		records    []models.DNSResponse
		signatures models.SignatureFile
		expected   []models.DetectedTechnology
	}{
		{
			name: "Detect Linode",
			records: []models.DNSResponse{
				{
					Domain:     "example.com.",
					RecordType: "NS",
					TTL:        21600,
					Value:      "ns1.linode.com.",
				},
			},
			signatures: models.SignatureFile{
				Signatures: []models.Signature{
					{
						Name:        "Linode",
						Category:    "Hosting Provider",
						Description: "Linode cloud hosting",
						RecordTypes: []string{"NS"},
						Patterns:    []string{"\\.linode\\.com\\."},
						Website:     "https://www.linode.com/",
					},
				},
			},
			expected: []models.DetectedTechnology{
				{
					Name:        "Linode",
					Category:    "Hosting Provider",
					Description: "Linode cloud hosting",
					Website:     "https://www.linode.com/",
					Evidence:    "ns1.linode.com.",
					RecordType:  "NS",
				},
			},
		},
		{
			name: "Detect SPF",
			records: []models.DNSResponse{
				{
					Domain:     "example.com.",
					RecordType: "TXT",
					TTL:        300,
					Value:      "v=spf1 +a +mx ~all",
				},
			},
			signatures: models.SignatureFile{
				Signatures: []models.Signature{
					{
						Name:        "SPF",
						Category:    "Email Security",
						Description: "Sender Policy Framework",
						RecordTypes: []string{"TXT"},
						Patterns:    []string{"v=spf1\\s.*"},
						Website:     "https://dmarcian.com/spf-overview/",
					},
				},
			},
			expected: []models.DetectedTechnology{
				{
					Name:        "SPF",
					Category:    "Email Security",
					Description: "Sender Policy Framework",
					Website:     "https://dmarcian.com/spf-overview/",
					Evidence:    "v=spf1 +a +mx ~all",
					RecordType:  "TXT",
				},
			},
		},
		{
			name: "No match",
			records: []models.DNSResponse{
				{
					Domain:     "example.com.",
					RecordType: "A",
					TTL:        300,
					Value:      "192.0.2.1",
				},
			},
			signatures: models.SignatureFile{
				Signatures: []models.Signature{
					{
						Name:        "SPF",
						Category:    "Email Security",
						Description: "Sender Policy Framework",
						RecordTypes: []string{"TXT"},
						Patterns:    []string{"v=spf1\\s.*"},
						Website:     "https://dmarcian.com/spf-overview/",
					},
				},
			},
			expected: []models.DetectedTechnology{},
		},
		{
			name: "Multiple matches, deduplicate",
			records: []models.DNSResponse{
				{
					Domain:     "example.com.",
					RecordType: "NS",
					TTL:        21600,
					Value:      "ns1.linode.com.",
				},
				{
					Domain:     "example.com.",
					RecordType: "NS",
					TTL:        21600,
					Value:      "ns2.linode.com.",
				},
			},
			signatures: models.SignatureFile{
				Signatures: []models.Signature{
					{
						Name:        "Linode",
						Category:    "Hosting Provider",
						Description: "Linode cloud hosting",
						RecordTypes: []string{"NS"},
						Patterns:    []string{"\\.linode\\.com\\."},
						Website:     "https://www.linode.com/",
					},
				},
			},
			expected: []models.DetectedTechnology{
				{
					Name:        "Linode",
					Category:    "Hosting Provider",
					Description: "Linode cloud hosting",
					Website:     "https://www.linode.com/",
					Evidence:    "ns1.linode.com.",
					RecordType:  "NS",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DetectTechnologies(tc.records, tc.signatures)
			
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestContainsString(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "Contains item",
			slice:    []string{"a", "b", "c"},
			item:     "b",
			expected: true,
		},
		{
			name:     "Does not contain item",
			slice:    []string{"a", "b", "c"},
			item:     "d",
			expected: false,
		},
		{
			name:     "Empty slice",
			slice:    []string{},
			item:     "a",
			expected: false,
		},
		{
			name:     "Nil slice",
			slice:    nil,
			item:     "a",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := containsString(tc.slice, tc.item)
			
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
