package models

// DNSResponse holds the parsed DNS record data
type DNSResponse struct {
	Domain     string `json:"domain"`
	RecordType string `json:"recordType"`
	TTL        uint32 `json:"ttl"`
	Value      string `json:"value"`
}

// DetectedTechnology represents a detected technology instance
type DetectedTechnology struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Evidence    string `json:"evidence"`
	RecordType  string `json:"recordType"`
}
