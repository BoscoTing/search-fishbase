package csvutil

import (
	"encoding/csv"
	"fishbase/scraper"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ReadCsvFile reads a CSV file and returns its records
func ReadCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer f.Close()

	// Read the first few bytes to check for UTF-8 BOM
	bom := make([]byte, 3)
	_, err = f.Read(bom)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	// If file starts with BOM, continue reading after it
	// Otherwise seek back to start of file
	if bom[0] != 0xEF || bom[1] != 0xBB || bom[2] != 0xBF {
		_, err = f.Seek(0, 0)
		if err != nil {
			return nil, fmt.Errorf("error seeking in file %s: %w", filePath, err)
		}
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error parsing file as CSV for %s: %w", filePath, err)
	}

	return records, nil
}

// ProcessCsvFile processes a CSV file with species names and call scrapeFunc for each row
func ProcessCsvFile(filePath string, scrapeFunc func(string) scraper.FishInfo) (string, error) {
	records, err := ReadCsvFile(filePath)
	if err != nil {
		return "", err
	}

	if len(records) == 0 {
		return "", fmt.Errorf("no records found in CSV file")
	}

	// Skip header row if present
	startIdx := 0
	if len(records) > 1 && (strings.Contains(strings.ToLower(records[0][0]), "species") ||
		strings.Contains(strings.ToLower(records[0][0]), "name")) {
		startIdx = 1
	}

	var fishes []scraper.FishInfo
	for i := startIdx; i < len(records); i++ {
		if len(records[i]) == 0 || len(records[i][0]) == 0 {
			continue // Skip empty rows
		}

		name := records[i][0]
		name = strings.TrimSpace(name)
		name = strings.ReplaceAll(name, " ", "-")

		// Add a small delay between requests to be polite to the server
		if i > startIdx {
			time.Sleep(time.Duration(300) * time.Millisecond)
		}

		fishes = append(fishes, scrapeFunc(name))
	}

	if len(fishes) == 0 {
		return "", fmt.Errorf("no valid fish species found in CSV file")
	}

	// Extract the jobId from the filename (timestamp part)
	baseFileName := filepath.Base(filePath)
	parts := strings.Split(baseFileName, "_")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid file name format, cannot extract jobId")
	}
	jobId := parts[0]

	return ExportCsvFile("./results", fishes, jobId)
}
