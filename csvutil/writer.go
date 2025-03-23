package csvutil

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"fishbase/scraper"
)

// ExportCsvFile exports data to a CSV file and returns the file path
func ExportCsvFile(resultsDir string, records []scraper.FishInfo, jobId string) (string, error) {
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create results directory: %w", err)
	}

	// Use the jobId in the filename to match the pattern expected by statusHandler
	filePath := filepath.Join(resultsDir, fmt.Sprintf("result_%s.csv", jobId))

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{
		"Name",
		"MaxLength",
		"MaxLengthUnit",
		"BayesianA",
		"BayesianAMin",
		"BayesianAMax",
		"BayesianB",
		"BayesianBMin",
		"BayesianBMax",
	}); err != nil {
		return "", fmt.Errorf("error writing header to CSV: %w", err)
	}

	// Write data
	for _, record := range records {
		if err := writer.Write(
			[]string{
				record.Name,
				record.Shape.MaxLength,
				record.Shape.MaxLengthUnit,
				record.Estimation.BayesianA,
				record.Estimation.BayesianAMin,
				record.Estimation.BayesianAMax,
				record.Estimation.BayesianB,
				record.Estimation.BayesianBMin,
				record.Estimation.BayesianBMax,
			},
		); err != nil {
			return "", fmt.Errorf("error writing record to CSV: %w", err)
		}
	}

	return filePath, nil
}
