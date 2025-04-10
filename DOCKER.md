# Docker for RADAR

This document explains how to use the Docker container for RADAR (Recognition and DNS Analysis for Resource detection).

## Getting Started

### Pull the Docker Image

```bash
docker pull elitesecuritysystems/radar:latest
```

### Basic Usage

To scan a domain:

```bash
docker run elitesecuritysystems/radar -domain example.com
```

To include all DNS records in the output:

```bash
docker run elitesecuritysystems/radar -domain example.com -all-records
```

To force signature updates:

```bash
docker run elitesecuritysystems/radar -update-signatures
```

## Advanced Usage

### Setting Timeout

```bash
docker run elitesecuritysystems/radar -domain example.com -timeout 30
```

### Debug Mode

```bash
docker run elitesecuritysystems/radar -domain example.com -debug
```

### Limiting Record Collection

```bash
docker run elitesecuritysystems/radar -domain example.com -max-records 2000
```

### Viewing Version Information

```bash
docker run elitesecuritysystems/radar -version
```

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

### Using Docker Compose

```bash
# Build and run
docker-compose build
docker-compose run --rm radar -domain yourdomain.com
```

## Customizing the Build

You can customize the build by setting environment variables:

```bash
VERSION=1.2.0 BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') COMMIT=$(git rev-parse --short HEAD) docker-compose build
```

## Tips

- The container automatically manages signature files, downloading the latest from GitHub when needed
- Results are written to stdout in JSON format
- For continuous integration or automation, you can pipe the output to other tools

## Troubleshooting

If you encounter issues:

1. Try running with the `-debug` flag for more detailed output
2. Ensure your container has network access
3. Check if the signature update is working by running with `-update-signatures`

For more help, please open an issue on GitHub.
