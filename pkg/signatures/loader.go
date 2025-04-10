package signatures

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Elite-Security-Systems/radar/internal/models"
	"github.com/Elite-Security-Systems/radar/internal/utils"
)

// LoadFromFile loads signatures from a JSON file
func LoadFromFile(path string) (models.SignatureFile, error) {
	var signatures models.SignatureFile

	// Try to find the file in standard locations if not found directly
	resolvedPath, err := utils.FindFile(path)
	if err != nil {
		return signatures, err
	}

	// Read the file
	data, err := ioutil.ReadFile(resolvedPath)
	if err != nil {
		// Handle common error cases with helpful messages
		if os.IsNotExist(err) {
			return signatures, fmt.Errorf("signatures file not found: %s\nPlease ensure the file exists or specify the correct path with -signatures", resolvedPath)
		}
		if os.IsPermission(err) {
			return signatures, fmt.Errorf("permission denied when reading signatures file: %s\nPlease check file permissions", resolvedPath)
		}
		return signatures, err
	}

	// Parse the JSON
	err = json.Unmarshal(data, &signatures)
	if err != nil {
		return signatures, fmt.Errorf("error parsing signatures file %s: %v", resolvedPath, err)
	}

	// Validate signatures
	if len(signatures.Signatures) == 0 {
		return signatures, fmt.Errorf("no valid signatures found in file: %s", resolvedPath)
	}

	return signatures, nil
}

// SaveToFile saves signatures to a JSON file
func SaveToFile(signatures models.SignatureFile, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %v", dir, err)
		}
	}

	// Marshal the signatures to JSON
	data, err := json.MarshalIndent(signatures, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding signatures: %v", err)
	}

	// Write to file
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing signatures to file %s: %v", path, err)
	}

	return nil
}
