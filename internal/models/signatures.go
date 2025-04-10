package models

// Signature represents a technology signature with regex patterns
type Signature struct {
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	RecordTypes []string `json:"recordTypes"`
	Patterns    []string `json:"patterns"`
	Website     string   `json:"website"`
}

// SignatureFile contains all technology signatures
type SignatureFile struct {
	Signatures []Signature `json:"signatures"`
}
