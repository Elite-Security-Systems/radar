#!/bin/sh
set -e

# If the first argument is "help" or "-h", display help message
if [ "$1" = "help" ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
  echo "RADAR: Recognition and DNS Analysis for Resource detection"
  echo ""
  echo "Usage: docker run elitesecuritysystems/radar [OPTIONS]"
  echo ""
  echo "Examples:"
  echo "  docker run elitesecuritysystems/radar -domain example.com"
  echo "  docker run elitesecuritysystems/radar -domain example.com -all-records"
  echo "  docker run elitesecuritysystems/radar -update-signatures"
  echo "  docker run -v $(pwd)/output:/output elitesecuritysystems/radar -domain example.com -o /output"
  echo "  docker run -v $(pwd)/output:/output elitesecuritysystems/radar -domain example.com -o /output/results.json"
  echo "  docker run elitesecuritysystems/radar -domain example.com -silent"
  echo ""
  echo "For more options, run: docker run elitesecuritysystems/radar -h"
  exit 0
fi

# If no arguments are given, display usage and exit
if [ $# -eq 0 ]; then
  echo "Error: Please provide a domain name with -domain flag"
  echo "Usage: docker run elitesecuritysystems/radar -domain example.com"
  echo "For more options, run: docker run elitesecuritysystems/radar help"
  exit 1
fi

# Check if the output directory is specified
output_dir=""
is_next_output=0

for arg in "$@"; do
  if [ $is_next_output -eq 1 ]; then
    output_dir="$arg"
    is_next_output=0
  elif [ "$arg" = "-o" ]; then
    is_next_output=1
  fi
done

# If an output directory was specified, ensure it exists
if [ -n "$output_dir" ]; then
  # Check if it's a json file or a directory
  if echo "$output_dir" | grep -q "\.json$"; then
    # It's a file, create parent directory
    mkdir -p "$(dirname "$output_dir")"
    # Ensure parent directory has correct permissions
    chmod 777 "$(dirname "$output_dir")"
  else
    # It's a directory, create it
    mkdir -p "$output_dir"
    # Ensure directory has correct permissions
    chmod 777 "$output_dir"
  fi
fi

# Run radar with the given arguments
exec radar "$@"
