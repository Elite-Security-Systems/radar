package models

// Result is the final output structure
type Result struct {
	Domain               string               `json:"domain"`
	DetectedTechnologies []DetectedTechnology `json:"detectedTechnologies"`
	AllRecords           []DNSResponse        `json:"allRecords,omitempty"`
}
