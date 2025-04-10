package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Elite-Security-Systems/radar/internal/analyzer"
	"github.com/Elite-Security-Systems/radar/internal/models"
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
		targetListFile    string
		signaturesPath    string
		includeAllRecords bool
		timeout           int
		debugMode         bool
		maxRecords        int
		showVersion       bool
		forceUpdate       bool
		silentMode        bool
		outputPath        string
		verboseOutput     bool
	)

	flag.StringVar(&domainName, "domain", "", "Domain name to analyze")
	flag.StringVar(&targetListFile, "l", "", "File containing list of domains to analyze (one per line)")
	flag.StringVar(&signaturesPath, "signatures", "data/signatures.json", "Path to signatures file")
	flag.BoolVar(&includeAllRecords, "all-records", false, "Include all records in JSON output")
	flag.IntVar(&timeout, "timeout", 15, "Query timeout in seconds")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug output")
	flag.IntVar(&maxRecords, "max-records", 1000, "Maximum number of records to collect (prevents hangs)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&forceUpdate, "update-signatures", false, "Force update signatures from GitHub")
	flag.BoolVar(&silentMode, "silent", false, "Silent mode - suppress all non-error output")
	flag.StringVar(&outputPath, "o", "", "Output file path or directory for results (if directory, creates JSON files named by domain)")
	flag.BoolVar(&verboseOutput, "verbose", false, "Show progress information when processing multiple domains")
	flag.Parse()

	// Show version information if requested
	if showVersion {
		printVersion()
		os.Exit(0)
	}

	// Force update signatures if requested
	if forceUpdate {
		if verboseOutput && !silentMode {
			fmt.Println("Updating signatures from GitHub repository...")
		}
		err := signatures.DownloadSignatures(signatures.DefaultSignaturesURL, signatures.DefaultCachePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating signatures: %v\n", err)
			os.Exit(1)
		}
		if verboseOutput && !silentMode {
			fmt.Printf("Signatures updated successfully to %s\n", signatures.DefaultCachePath)
		}
		
		// If no domain or target list provided, exit after update
		if domainName == "" && targetListFile == "" {
			os.Exit(0)
		}
	}

	// Validate input: either domain or target list must be provided
	if domainName == "" && targetListFile == "" {
		fmt.Fprintln(os.Stderr, "Error: Please provide a domain name with -domain flag or a target list with -l flag")
		flag.Usage()
		os.Exit(1)
	}

	// If both domain and target list are provided, warn user and prioritize target list
	if domainName != "" && targetListFile != "" {
		fmt.Fprintf(os.Stderr, "Warning: Both domain and target list provided. Using target list and ignoring single domain.\n")
	}

	// Create output directory if specified and it's a directory
	if outputPath != "" && !strings.HasSuffix(strings.ToLower(outputPath), ".json") {
		err := os.MkdirAll(outputPath, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	} else if outputPath != "" && strings.HasSuffix(strings.ToLower(outputPath), ".json") {
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

	// Load signatures
	sigs, err := signatures.LoadFromFile(signaturesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading signatures: %v\n", err)
		os.Exit(1)
	}

	if debugMode && !silentMode {
		fmt.Fprintf(os.Stderr, "[DEBUG] Loaded %d signatures from %s\n", len(sigs.Signatures), signaturesPath)
	}

	// If target list is provided, process it
	if targetListFile != "" {
		err = processTargetList(targetListFile, outputPath, sigs, includeAllRecords, timeout, debugMode, maxRecords, silentMode, verboseOutput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing target list: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Process single domain
	processSingleDomain(domainName, outputPath, sigs, includeAllRecords, timeout, debugMode, maxRecords, silentMode, verboseOutput)
}

// processSingleDomain analyzes a single domain and handles output
func processSingleDomain(domain, outputPath string, sigs models.SignatureFile, includeAllRecords bool, timeout int, debugMode bool, maxRecords int, silentMode bool, verboseOutput bool) {
	// Initialize analyzer with configuration
	config := analyzer.Config{
		Domain:         domain,
		Timeout:        time.Duration(timeout) * time.Second,
		Debug:          debugMode && !silentMode, // Disable debug output in silent mode
		MaxRecords:     maxRecords,
		IncludeRecords: includeAllRecords,
	}
	
	result, err := analyzer.AnalyzeDomain(config, sigs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing domain %s: %v\n", domain, err)
		return
	}

	// Generate JSON output
	jsonOutput, err := utils.FormatJSON(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON output for %s: %v\n", domain, err)
		return
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
			fmt.Fprintf(os.Stderr, "Error saving output to file for %s: %v\n", domain, err)
			return
		}
		
		if verboseOutput && !silentMode {
			fmt.Fprintf(os.Stderr, "Results for %s saved to: %s\n", domain, finalPath)
		}
	} else {
		// No output path specified, print to stdout (unless in silent mode)
		if !silentMode {
			fmt.Println(jsonOutput)
		}
	}
}

// processTargetList reads domains from a file and processes each one
func processTargetList(targetListFile, outputPath string, sigs models.SignatureFile, includeAllRecords bool, timeout int, debugMode bool, maxRecords int, silentMode bool, verboseOutput bool) error {
	// Open the target list file
	file, err := os.Open(targetListFile)
	if err != nil {
		return fmt.Errorf("error opening target list file: %v", err)
	}
	defer file.Close()

	// If output path is a specific JSON file and we have multiple targets, we need to handle differently
	isOutputFile := strings.HasSuffix(strings.ToLower(outputPath), ".json")
	var combinedResults []models.Result
	
	// Count domains for progress indication
	var totalDomains int
	if verboseOutput {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if line := strings.TrimSpace(scanner.Text()); line != "" && !strings.HasPrefix(line, "#") {
				totalDomains++
			}
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading target list file: %v", err)
		}

		// Reset file for processing
		_, err = file.Seek(0, 0)
		if err != nil {
			return fmt.Errorf("error resetting file: %v", err)
		}
	}
	
	// Process each domain
	scanner := bufio.NewScanner(file)
	var processedDomains int
	for scanner.Scan() {
		// Get domain, skipping empty lines and comments
		domain := strings.TrimSpace(scanner.Text())
		if domain == "" || strings.HasPrefix(domain, "#") {
			continue
		}
		
		processedDomains++
		if verboseOutput && !silentMode {
			fmt.Fprintf(os.Stderr, "Processing domain %d/%d: %s\n", processedDomains, totalDomains, domain)
		}

		// Initialize analyzer with configuration
		config := analyzer.Config{
			Domain:         domain,
			Timeout:        time.Duration(timeout) * time.Second,
			Debug:          debugMode && !silentMode,
			MaxRecords:     maxRecords,
			IncludeRecords: includeAllRecords,
		}
		
		result, err := analyzer.AnalyzeDomain(config, sigs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error analyzing domain %s: %v\n", domain, err)
			continue
		}

		// If output is a specific file, collect results for combined output
		if isOutputFile && outputPath != "" {
			combinedResults = append(combinedResults, *result)
		} else {
			// Generate JSON output for individual domain
			jsonOutput, err := utils.FormatJSON(result)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating JSON output for %s: %v\n", domain, err)
				continue
			}

			// Handle individual output
			if outputPath != "" {
				// Create a filename from the domain and timestamp
				timestamp := time.Now().Format("20060102-150405")
				filename := fmt.Sprintf("%s_%s.json", result.Domain, timestamp)
				finalPath := filepath.Join(outputPath, filename)
				
				err := utils.SaveToFile(jsonOutput, finalPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error saving output to file for %s: %v\n", domain, err)
					continue
				}
				
				if verboseOutput && !silentMode {
					fmt.Fprintf(os.Stderr, "Results for %s saved to: %s\n", domain, finalPath)
				}
			} else {
				// No output path specified, print to stdout (unless in silent mode)
				if !silentMode {
					// Print only the JSON output with no progress info or separators
					fmt.Println(jsonOutput)
				}
			}
		}
	}

	// If we have a combined results file, write it now
	if isOutputFile && outputPath != "" && len(combinedResults) > 0 {
		// Create a combined results structure
		combinedOutput := struct {
			Timestamp string         `json:"timestamp"`
			Count     int            `json:"count"`
			Results   []models.Result `json:"results"`
		}{
			Timestamp: time.Now().Format(time.RFC3339),
			Count:     len(combinedResults),
			Results:   combinedResults,
		}
		
		// Convert to JSON
		jsonOutput, err := utils.FormatJSON(combinedOutput)
		if err != nil {
			return fmt.Errorf("error generating combined JSON output: %v", err)
		}
		
		// Save to file
		err = utils.SaveToFile(jsonOutput, outputPath)
		if err != nil {
			return fmt.Errorf("error saving combined output to file: %v", err)
		}
		
		if verboseOutput && !silentMode {
			fmt.Fprintf(os.Stderr, "Combined results for %d domains saved to: %s\n", len(combinedResults), outputPath)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading target list file: %v", err)
	}

	return nil
}

// printVersion prints version information
func printVersion() {
	fmt.Printf("RADAR: Recognition and DNS Analysis for Resource detection\n")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Build Date: %s\n", BuildDate)
	fmt.Printf("Commit: %s\n", Commit)
	fmt.Printf("Developed by Elite Security Systems (elitesecurity.systems)\n")
}
