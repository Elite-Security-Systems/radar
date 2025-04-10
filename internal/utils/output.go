package utils

import (
	"encoding/json"
)

// FormatJSON formats a struct as indented JSON
func FormatJSON(v interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
