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

# Run radar with the given arguments
exec radar "$@"
