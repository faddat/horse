# Cosmos Chain Address Prefix Converter

A tool for converting Cosmos SDK chain address prefixes in genesis files and CSV files.

## Overview

This tool allows you to:

1. Convert all addresses in a Cosmos SDK genesis file from one prefix to another (e.g., `unicorn1...` to `esim1...`)
2. Process CSV files containing addresses to convert their prefixes as well

## Requirements

- Go 1.21 or higher

## Installation

```bash
# Clone the repository
git clone https://github.com/faddat/horse/snapshot-converter
cd snapshot-converter

# Install dependencies
go mod download
```

## Usage

### Building the Tool

```bash
go build -o converter main.go
```

### Converting Genesis and CSV Files

```bash
# Basic usage
./converter -input original_genesis.json -output new_genesis.json

# Convert genesis file and process CSV files in a directory
./converter -input original_genesis.json -output new_genesis.json -csv-dir ./snapshot
```

### Command Line Options

- `-input`: Path to the input genesis file (required)
- `-output`: Path where the converted genesis file should be saved (required)
- `-csv-dir`: Directory containing CSV files to process (optional)

## How It Works

The tool performs the following operations:

1. Reads the genesis JSON file and parses it
2. Recursively searches through the JSON structure for any strings containing the old prefix (`unicorn`)
3. Replaces the old prefix with the new prefix (`esim`)
4. Writes the updated genesis to the output file
5. If a CSV directory is provided, it processes all CSV files, converting any addresses found

### Genesis File Conversion Details

The tool handles:
- Regular bech32 addresses (e.g., `unicorn1abc123...` → `esim1abc123...`)
- Factory token denoms (e.g., `factory/unicorn1abc.../token` → `factory/esim1abc.../token`)
- Chain IDs that include the prefix

## Notes

- This is a simple string replacement tool and doesn't validate the checksums of the bech32 addresses
- For a production environment, consider using proper bech32 encoding/decoding libraries to ensure valid addresses



