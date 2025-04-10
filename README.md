# RADAR: Recognition and DNS Analysis for Resource detection

<p align="center">
  <img src="static/radar-logo.png" alt="RADAR Logo"/>
</p>

## About RADAR

RADAR (Recognition and DNS Analysis for Resource detection) is an advanced DNS reconnaissance tool designed to identify technologies and services used by domains through their DNS footprints. Developed by [Elite Security Systems](https://elitesecurity.systems), RADAR can detect hundreds of technologies including cloud services, email providers, CDNs, security services, and more.

## Features

- 🔍 **Comprehensive DNS Scanning**: Queries all relevant DNS record types (A, AAAA, CNAME, MX, TXT, NS, SOA, SRV, CAA, etc.)
- 🛡️ **Technology Detection**: Identifies technologies using pattern matching against an extensive signature database
- ⚡ **Performance Optimized**: Uses parallel queries and multiple resolvers for efficient scanning
- 🧩 **Extensible**: Easy to add new technology signatures via the JSON signature database
- 📊 **Detailed Reporting**: Produces structured JSON output for easy integration with other tools
- 🌐 **Robust Resolving**: Leverages both system DNS resolver and public DNS services for maximum coverage

## Installation

Run the following command to install the latest version:
```bash
go install -v github.com/Elite-Security-Systems/radar/cmd/radar@latest
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

## Output Example

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
    {
      "name": "Google Workspace",
      "category": "Email & Collaboration",
      "description": "Google Workspace (formerly G Suite) email services",
      "website": "https://workspace.google.com",
      "evidence": "aspmx.l.google.com.",
      "recordType": "MX"
    }
  ]
}
```

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
