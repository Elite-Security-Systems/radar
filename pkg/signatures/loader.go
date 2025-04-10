package signatures

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Elite-Security-Systems/radar/internal/models"
	"github.com/Elite-Security-Systems/radar/internal/utils"
)

const (
	// DefaultSignaturesURL is the URL to download signatures from
	DefaultSignaturesURL = "https://raw.githubusercontent.com/Elite-Security-Systems/radar/refs/heads/main/data/signatures.json"
	// DefaultCachePath is the path to store downloaded signatures
	DefaultCachePath = "/tmp/radar-sigs.json"
	// MaxCacheAge is the maximum age of the cached signatures file before redownloading
	MaxCacheAge = 7 * 24 * time.Hour // 1 week
)

// LoadFromFile loads signatures from a JSON file
func LoadFromFile(path string) (models.SignatureFile, error) {
	var signatures models.SignatureFile

	// Check if this is the default path, and handle automatic downloads
	if path == "data/signatures.json" {
		// Try to get a local or download a fresh copy
		resolvedPath, err := GetOrDownloadSignatures()
		if err == nil {
			path = resolvedPath
		}
	}

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

// GetOrDownloadSignatures checks if we need to download a new signatures file,
// downloads if needed, and returns the path to the signatures file
func GetOrDownloadSignatures() (string, error) {
	// Check if cache file exists and is recent enough
	info, err := os.Stat(DefaultCachePath)
	
	// If the file doesn't exist or is older than MaxCacheAge, download a fresh copy
	if err != nil || time.Since(info.ModTime()) > MaxCacheAge {
		// Download the signatures file
		err := DownloadSignatures(DefaultSignaturesURL, DefaultCachePath)
		if err != nil {
			return "", fmt.Errorf("failed to download signatures: %v", err)
		}
	}
	
	return DefaultCachePath, nil
}

// DownloadSignatures downloads signatures from the specified URL to the specified path
func DownloadSignatures(url, path string) error {
	// Create the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading signatures: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error downloading signatures: HTTP %d", resp.StatusCode)
	}

	// Create the output file
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating signatures file: %v", err)
	}
	defer out.Close()

	// Copy the response body to the output file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing signatures file: %v", err)
	}

	return nil
}
