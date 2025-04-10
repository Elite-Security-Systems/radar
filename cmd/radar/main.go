package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Elite-Security-Systems/radar/internal/analyzer"
	"github.com/Elite-Security-Systems/radar/internal/utils"
	"github.com/Elite-Security-Systems/radar/pkg/signatures"
)

var (
	// Version information - to be set during build
	Version   = "dev"
	BuildDate = "unknown"
	Commit    = "none"
)

func main() {
	// Command line flags
	var (
		domainName        string
		signaturesPath    string
		includeAllRecords bool
		timeout           int
		debugMode         bool
		maxRecords        int
		showVersion       bool
		forceUpdate       bool
		silentMode        bool
		outputPath        string
	)

	flag.StringVar(&domainName, "domain", "", "Domain name to analyze")
	flag.StringVar(&signaturesPath, "signatures", "data/signatures.json", "Path to signatures file")
	flag.BoolVar(&includeAllRecords, "all-records", false, "Include all records in JSON output")
	flag.IntVar(&timeout, "timeout", 15, "Query timeout in seconds")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug output")
	flag.IntVar(&maxRecords, "max-records", 1000, "Maximum number of records to collect (prevents hangs)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&forceUpdate, "update-signatures", false, "Force update signatures from GitHub")
	flag.BoolVar(&silentMode, "silent", false, "Silent mode - suppress all non-error output")
	flag.StringVar(&outputPath, "o", "", "Output file path or directory for results (if directory, creates JSON files named by domain)")
	flag.Parse()

	// Show version information if requested
	if showVersion {
		printVersion()
		os.Exit(0)
	}

	// Force update signatures if requested
	if forceUpdate {
		if !silentMode {
			fmt.Println("Updating signatures from GitHub repository...")
		}
		err := signatures.DownloadSignatures(signatures.DefaultSignaturesURL, signatures.DefaultCachePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating signatures: %v\n", err)
			os.Exit(1)
		}
		if !silentMode {
			fmt.Printf("Signatures updated successfully to %s\n", signatures.DefaultCachePath)
		}
		
		// If no domain provided, exit after update
		if domainName == "" {
			os.Exit(0)
		}
	}

	// Validate domain name
	if domainName == "" {
		fmt.Fprintln(os.Stderr, "Error: Please provide a domain name with -domain flag")
		flag.Usage()
		os.Exit(1)
	}

	// Create output directory if specified and it's a directory
	if outputPath != "" {
		// Check if the path ends with .json, if not, treat it as a directory
		if !strings.HasSuffix(strings.ToLower(outputPath), ".json") {
			dirPath := outputPath
			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
				os.Exit(1)
			}
		} else {
			// It's a file, ensure the directory exists
			dirPath := filepath.Dir(outputPath)
			if dirPath != "" && dirPath != "." {
				err := os.MkdirAll(dirPath, 0755)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}

	// Load signatures
	sigs, err := signatures.LoadFromFile(signaturesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading signatures: %v\n", err)
		os.Exit(1)
	}

	if debugMode && !silentMode {
		fmt.Printf("[DEBUG] Loaded %d signatures from %s\n", len(sigs.Signatures), signaturesPath)
	}

	// Initialize analyzer with configuration
	config := analyzer.Config{
		Domain:         domainName,
		Timeout:        time.Duration(timeout) * time.Second,
		Debug:          debugMode && !silentMode, // Disable debug output in silent mode
		MaxRecords:     maxRecords,
		IncludeRecords: includeAllRecords,
	}
	
	result, err := analyzer.AnalyzeDomain(config, sigs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing domain: %v\n", err)
		os.Exit(1)
	}

	// Generate JSON output
	jsonOutput, err := utils.FormatJSON(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON output: %v\n", err)
		os.Exit(1)
	}

	// Handle output based on mode
	if outputPath != "" {
		var finalPath string
		
		// Determine if this is a file path or directory
		if strings.HasSuffix(strings.ToLower(outputPath), ".json") {
			// Use the exact path specified
			finalPath = outputPath
		} else {
			// Create a filename from the domain and timestamp
			timestamp := time.Now().Format("20060102-150405")
			filename := fmt.Sprintf("%s_%s.json", result.Domain, timestamp)
			finalPath = filepath.Join(outputPath, filename)
		}
		
		err := utils.SaveToFile(jsonOutput, finalPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving output to file: %v\n", err)
			os.Exit(1)
		}
		
		if !silentMode {
			fmt.Printf("Results saved to: %s\n", finalPath)
		}
		// In silent mode, we don't print anything, not even the path
	} else {
		// No output path specified, print to stdout (unless in silent mode)
		if !silentMode {
			fmt.Println(jsonOutput)
		}
	}
}

// printVersion prints version information
func printVersion() {
	fmt.Printf("RADAR: Recognition and DNS Analysis for Resource detection\n")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Build Date: %s\n", BuildDate)
	fmt.Printf("Commit: %s\n", Commit)
	fmt.Printf("Developed by Elite Security Systems (elitesecurity.systems)\n")
}
