# RADAR: Recognition and DNS Analysis for Resource detection

<p align="center">
  <img src="static/radar-logo.png" alt="RADAR Logo"/>
</p>

## About RADAR

RADAR (Recognition and DNS Analysis for Resource detection) is an advanced DNS reconnaissance tool designed to identify technologies and services used by domains through their DNS footprints. Developed by [Elite Security Systems](https://elitesecurity.systems), RADAR can detect hundreds of technologies including cloud services, email providers, CDNs, security services, and more.

## Features

- � **Comprehensive DNS Scanning**: Queries all relevant DNS record types (A, AAAA, CNAME, MX, TXT, NS, SOA, SRV, CAA, etc.)
- �️ **Technology Detection**: Identifies technologies using pattern matching against an extensive signature database
- ⚡ **Performance Optimized**: Uses parallel queries and multiple resolvers for efficient scanning
- � **Extensible**: Easy to add new technology signatures via the JSON signature database
- � **Detailed Reporting**: Produces structured JSON output for easy integration with other tools
- � **Robust Resolving**: Leverages both system DNS resolver and public DNS services for maximum coverage
- � **Auto-Updates**: Automatically downloads the latest signatures from GitHub

## Installation

Run the following command to install the latest version:
```bash
go install -v github.com/Elite-Security-Systems/radar/cmd/radar@latest
```

### Using Prebuilt Binaries

Download the latest release for your platform from the [releases page](https://github.com/Elite-Security-Systems/radar/releases).

### Using Docker

Pull and run the Docker image:
```bash
docker pull elitesecuritysystems/radar:latest
docker run elitesecuritysystems/radar -domain example.com
```

### From Source

```bash
# Clone the repository
git clone https://github.com/Elite-Security-Systems/radar.git
cd radar

# Build
go build -o radar ./cmd/radar

# Install (optional)
go install ./cmd/radar
```

## Quick Start

```bash
# Basic domain scan
radar -domain example.com

# Scan with all DNS records in output
radar -domain example.com -all-records

# Use custom signatures file
radar -domain example.com -signatures /path/to/signatures.json

# Enable debug output
radar -domain example.com -debug

# Set custom timeout
radar -domain example.com -timeout 30
```

## Output Options

### JSON Output Format

By default, RADAR outputs results to stdout in JSON format:

```json
{
  "domain": "example.com",
  "detectedTechnologies": [
    {
      "name": "Cloudflare",
      "category": "CDN & Security",
      "description": "Cloudflare CDN and security services",
      "website": "https://www.cloudflare.com",
      "evidence": "ns1.cloudflare.com.",
      "recordType": "NS"
    },
    ...
  ]
}
```

### Save Results to File

You can save results to a file using the `-o` flag:

```bash
# Save to a specific file
radar -domain example.com -o results.json

# Save to a directory (creates a timestamped file)
radar -domain example.com -o /path/to/output/
```

When you specify a directory, RADAR will create a file with a format like: 
```
/path/to/output/example.com_20250410-150405.json
```

### Silent Mode

Use the `-silent` flag to suppress all non-error output:

```bash
# Silent mode with output to a file
radar -domain example.com -o results.json -silent

# No output appears if successful, only errors would be shown
```

This is particularly useful for:
- Scheduled tasks and cron jobs
- CI/CD pipelines
- Batch processing where you want to avoid cluttering logs

## Advanced Usage

### Updating Signatures

RADAR automatically manages signature files, downloading the latest from GitHub when needed. You can force an update with:

```bash
radar -update-signatures
```

### Including All DNS Records

By default, RADAR only includes detected technologies in the output. You can include all DNS records with:

```bash
radar -domain example.com -all-records
```

### Batch Processing Example

```bash
#!/bin/bash
OUTPUT_DIR="./results"
mkdir -p "$OUTPUT_DIR"

# Read domains from a file
while read domain; do
  echo "Processing $domain..."
  radar -domain "$domain" -o "$OUTPUT_DIR" -silent
  if [ $? -eq 0 ]; then
    echo "  Successfully analyzed $domain"
  else
    echo "  Failed to analyze $domain"
  fi
done < domains.txt

echo "All scans complete. Results saved in $OUTPUT_DIR directory."
```

## Docker Usage

### Basic Usage

```bash
docker run elitesecuritysystems/radar -domain example.com
```

### Save Results to Host Machine

```bash
docker run -v "$(pwd)/output:/output" elitesecuritysystems/radar -domain example.com -o /output
```

### Silent Mode with Docker

```bash
docker run -v "$(pwd)/output:/output" elitesecuritysystems/radar -domain example.com -o /output -silent
```

### Force Update Signatures

```bash
docker run elitesecuritysystems/radar -update-signatures
```

## Command Line Options

| Flag | Description |
|------|-------------|
| `-domain` | Domain name to analyze |
| `-o` | Output file path or directory for results |
| `-all-records` | Include all records in JSON output |
| `-timeout` | Query timeout in seconds (default: 15) |
| `-debug` | Enable debug output |
| `-max-records` | Maximum number of records to collect (default: 1000) |
| `-signatures` | Path to signatures file (default: data/signatures.json) |
| `-update-signatures` | Force update signatures from GitHub |
| `-silent` | Silent mode - suppress all non-error output |
| `-version` | Show version information |

## Custom Signatures

RADAR uses a JSON-based signature format for technology detection. You can extend the default signature set or create your own:

```json
{
  "signatures": [
    {
      "name": "My Custom Technology",
      "category": "Web Platform",
      "description": "Description of the technology",
      "recordTypes": ["TXT", "CNAME"],
      "patterns": [
        "regex-pattern-to-match",
        "another-pattern-.*"
      ],
      "website": "https://example.com"
    }
  ]
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Building the Docker Image

If you want to build the image yourself:

```bash
# Clone the repository
git clone https://github.com/Elite-Security-Systems/radar.git
cd radar

# Build using the provided script
chmod +x build-push.sh
./build-push.sh

# Or build and push to Docker Hub
./build-push.sh --push
```

## Troubleshooting

If you encounter issues:

1. Try running with the `-debug` flag for more detailed output
2. Ensure your tool/container has network access
3. Check if the signature update is working by running with `-update-signatures`
4. Make sure mounted volumes have correct permissions when using Docker

For more help, please open an issue on GitHub.

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgements

- [miekg/dns](https://github.com/miekg/dns) - DNS library for Go
- The Elite Security Systems team for continuous support and contributions

---

<p align="center">
  <a href="https://elitesecurity.systems">Elite Security Systems</a> •
  <a href="https://x.com/eliteSsystems">X</a>
</p>
