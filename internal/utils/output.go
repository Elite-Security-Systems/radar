package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// FormatJSON formats a struct as indented JSON
func FormatJSON(v interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// SaveToFile saves content to a file
func SaveToFile(content string, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %v", dir, err)
		}
	}

	// Write to file
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %v", path, err)
	}

	return nil
}
