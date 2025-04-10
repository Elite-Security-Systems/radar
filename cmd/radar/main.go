package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Elite-Security-Systems/radar/internal/analyzer"
	"github.com/Elite-Security-Systems/radar/internal/utils"
	"github.com/Elite-Security-Systems/radar/pkg/signatures"
)

var (
	// Version information - to be set during build
	Version   = "0.1.1"
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
	)

	flag.StringVar(&domainName, "domain", "", "Domain name to analyze")
	flag.StringVar(&signaturesPath, "signatures", "data/signatures.json", "Path to signatures file")
	flag.BoolVar(&includeAllRecords, "all-records", false, "Include all records in JSON output")
	flag.IntVar(&timeout, "timeout", 15, "Query timeout in seconds")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug output")
	flag.IntVar(&maxRecords, "max-records", 1000, "Maximum number of records to collect (prevents hangs)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&forceUpdate, "update-signatures", false, "Force update signatures from GitHub")
	flag.Parse()

	// Show version information if requested
	if showVersion {
		fmt.Printf("RADAR: Recognition and DNS Analysis for Resource detection\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Developed by Elite Security Systems (elitesecurity.systems)\n")
		os.Exit(0)
	}

	// Force update signatures if requested
	if forceUpdate {
		fmt.Println("Updating signatures from GitHub repository...")
		err := signatures.DownloadSignatures(signatures.DefaultSignaturesURL, signatures.DefaultCachePath)
		if err != nil {
			fmt.Printf("Error updating signatures: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Signatures updated successfully to %s\n", signatures.DefaultCachePath)
		
		// If no domain provided, exit after update
		if domainName == "" {
			os.Exit(0)
		}
	}

	// Validate domain name
	if domainName == "" {
		fmt.Println("Error: Please provide a domain name with -domain flag")
		flag.Usage()
		os.Exit(1)
	}

	// Load signatures
	sigs, err := signatures.LoadFromFile(signaturesPath)
	if err != nil {
		fmt.Printf("Error loading signatures: %v\n", err)
		os.Exit(1)
	}

	if debugMode {
		fmt.Printf("[DEBUG] Loaded %d signatures from %s\n", len(sigs.Signatures), signaturesPath)
	}

	// Initialize analyzer with configuration
	config := analyzer.Config{
		Domain:        domainName,
		Timeout:       time.Duration(timeout) * time.Second,
		Debug:         debugMode,
		MaxRecords:    maxRecords,
		IncludeRecords: includeAllRecords,
	}
	
	result, err := analyzer.AnalyzeDomain(config, sigs)
	if err != nil {
		fmt.Printf("Error analyzing domain: %v\n", err)
		os.Exit(1)
	}

	// Output the result
	jsonOutput, err := utils.FormatJSON(result)
	if err != nil {
		fmt.Printf("Error generating JSON output: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(jsonOutput)
}
