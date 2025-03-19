package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Constants for address conversion
const (
	OldPrefix = "unicorn"
	NewPrefix = "esim"
)

// Simple representation of a genesis file
type GenesisDoc struct {
	AppState json.RawMessage `json:"app_state"`
	// Other fields in genesis doc
	ChainID         string          `json:"chain_id"`
	GenesisTime     string          `json:"genesis_time"`
	ConsensusParams json.RawMessage `json:"consensus_params"`
	InitialHeight   string          `json:"initial_height"`
	Validators      json.RawMessage `json:"validators"`
}

func main() {
	// Define command line flags
	inputGenesisPtr := flag.String("input", "", "Input genesis.json file path")
	outputGenesisPtr := flag.String("output", "", "Output genesis.json file path")
	csvDirPtr := flag.String("csv-dir", "", "Directory containing CSV files to process (optional)")
	flag.Parse()

	// Check for required arguments
	if *inputGenesisPtr == "" || *outputGenesisPtr == "" {
		fmt.Println("Usage: go run main.go -input <genesis.json> -output <new_genesis.json> [-csv-dir <directory>]")
		os.Exit(1)
	}

	// Process the genesis file
	err := processGenesisFile(*inputGenesisPtr, *outputGenesisPtr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Process CSV files if directory provided
	if *csvDirPtr != "" {
		err := processCsvDirectory(*csvDirPtr)
		if err != nil {
			fmt.Printf("Error processing CSV files: %v\n", err)
			os.Exit(1)
		}
	}
}

// processGenesisFile handles the conversion of a genesis file
func processGenesisFile(inputFile, outputFile string) error {
	fmt.Printf("Processing genesis file %s -> %s\n", inputFile, outputFile)

	// Read the input genesis file
	genesisBytes, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading genesis file: %v", err)
	}

	// Parse the genesis doc
	var genesisDoc GenesisDoc
	err = json.Unmarshal(genesisBytes, &genesisDoc)
	if err != nil {
		return fmt.Errorf("error parsing genesis JSON: %v", err)
	}

	// Convert app_state JSON to map
	var appState map[string]interface{}
	err = json.Unmarshal(genesisDoc.AppState, &appState)
	if err != nil {
		return fmt.Errorf("error parsing app_state JSON: %v", err)
	}

	// Process the app state - recursively replace all unicorn prefixes with esim
	processAppState(appState)

	// Update chain ID if it contains the old prefix
	if strings.Contains(genesisDoc.ChainID, OldPrefix) {
		genesisDoc.ChainID = strings.ReplaceAll(genesisDoc.ChainID, OldPrefix, NewPrefix)
	}

	// Convert the updated app state back to JSON
	updatedAppState, err := json.Marshal(appState)
	if err != nil {
		return fmt.Errorf("error marshalling updated app_state: %v", err)
	}

	// Update the genesis doc with the new app state
	genesisDoc.AppState = updatedAppState

	// Marshal the entire genesis doc
	updatedGenesisBytes, err := json.MarshalIndent(genesisDoc, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling updated genesis: %v", err)
	}

	// Write the new genesis file
	err = os.WriteFile(outputFile, updatedGenesisBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing output genesis file: %v", err)
	}

	fmt.Printf("Successfully converted addresses from %s to %s and saved to %s\n", OldPrefix, NewPrefix, outputFile)
	return nil
}

// processAppState recursively processes the app state JSON structure
func processAppState(data interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		// Process each key-value pair in the map
		for key, value := range v {
			// Convert keys that might contain addresses
			if strings.Contains(key, OldPrefix) {
				newKey := replaceAddressInString(key)
				v[newKey] = value
				delete(v, key)
				// Continue processing with the new key
				key = newKey
			}

			// Handle string values that might contain addresses
			if strValue, ok := value.(string); ok && strings.Contains(strValue, OldPrefix) {
				v[key] = replaceAddressInString(strValue)
			} else {
				// Process the value recursively
				processAppState(value)
			}
		}
	case []interface{}:
		// Process each element in the array
		for i, element := range v {
			if strElement, ok := element.(string); ok && strings.Contains(strElement, OldPrefix) {
				v[i] = replaceAddressInString(strElement)
			} else {
				processAppState(element)
			}
		}
	}
}

// replaceAddressInString replaces all occurrences of the old prefix in a string
func replaceAddressInString(text string) string {
	// Handle factory token addresses (format: factory/unicorn.../token)
	factoryPattern := regexp.MustCompile(`factory/` + OldPrefix + `([a-zA-Z0-9]+)`)
	text = factoryPattern.ReplaceAllString(text, "factory/"+NewPrefix+"${1}")

	// Handle regular bech32 addresses
	// This is a simplified approach - in a real implementation you would want to use
	// proper bech32 decoder/encoder to ensure the checksum is valid
	bech32Pattern := regexp.MustCompile(OldPrefix + `([a-zA-Z0-9]+)`)
	text = bech32Pattern.ReplaceAllString(text, NewPrefix+"${1}")

	return text
}

// processCsvDirectory processes all CSV files in a directory
func processCsvDirectory(directory string) error {
	fmt.Printf("Processing CSV files in directory: %s\n", directory)

	// Get all CSV files
	files, err := filepath.Glob(filepath.Join(directory, "*.csv"))
	if err != nil {
		return fmt.Errorf("error finding CSV files: %v", err)
	}

	for _, csvFile := range files {
		err := processCsvFile(csvFile)
		if err != nil {
			fmt.Printf("Warning: Error processing %s: %v\n", csvFile, err)
			// Continue with other files
		}
	}

	return nil
}

// processCsvFile processes a single CSV file
func processCsvFile(csvFile string) error {
	// Open the CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	// Parse the CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	// Check if we have records and headers
	if len(records) == 0 {
		fmt.Printf("No records found in %s\n", csvFile)
		return nil
	}

	// Process the records
	modified := false
	for i, record := range records {
		for j, field := range record {
			if strings.Contains(field, OldPrefix) {
				records[i][j] = replaceAddressInString(field)
				modified = true
			}
		}
	}

	// If no modifications were made, skip writing
	if !modified {
		fmt.Printf("No addresses found to convert in %s\n", csvFile)
		return nil
	}

	// Create the output filename
	outputFile := strings.TrimSuffix(csvFile, ".csv") + "_" + NewPrefix + ".csv"

	// Create the output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output CSV file: %v", err)
	}
	defer outFile.Close()

	// Write the updated records
	writer := csv.NewWriter(outFile)
	err = writer.WriteAll(records)
	if err != nil {
		return fmt.Errorf("error writing to CSV file: %v", err)
	}

	fmt.Printf("Successfully converted addresses in %s and saved to %s\n", csvFile, outputFile)
	return nil
}
