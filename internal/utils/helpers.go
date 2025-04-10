package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetExecutablePath returns the absolute path of the current executable
func GetExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

// FindFile tries to find a file in multiple possible locations
func FindFile(filename string) (string, error) {
	// Try current directory first
	if _, err := os.Stat(filename); err == nil {
		abs, err := filepath.Abs(filename)
		if err != nil {
			return filename, nil
		}
		return abs, nil
	}

	// Try executable directory
	exePath, err := GetExecutablePath()
	if err == nil {
		path := filepath.Join(exePath, filename)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Try relative to source code (for development)
	_, file, _, _ := runtime.Caller(0)
	sourceRoot := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	path := filepath.Join(sourceRoot, "data", filepath.Base(filename))
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	// Common installation directories
	var commonDirs []string
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		commonDirs = []string{
			"/usr/share/radar",
			"/usr/local/share/radar",
			"/opt/radar",
		}
	} else if runtime.GOOS == "windows" {
		commonDirs = []string{
			filepath.Join(os.Getenv("PROGRAMFILES"), "Radar"),
			filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Radar"),
		}
	}

	for _, dir := range commonDirs {
		path := filepath.Join(dir, filename)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Not found
	return filename, nil // Return the original filename so the caller can handle the error
}
