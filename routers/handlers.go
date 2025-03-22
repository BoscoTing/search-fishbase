package routers

import (
	"fishbase/csvutil"
	"fishbase/scraper"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// homeHandler renders the home page
func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Fishbase Scraper",
	})
}

func uploadHandler(c *gin.Context) {
	// Single file upload
	file, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title": "Fishbase Scraper",
			"error": "Failed to get file: " + err.Error(),
		})
		return
	}

	// Check file size
	if file.Size > MaxFileSize {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title": "Fishbase Scraper",
			"error": "File size exceeds the limit (10MB)",
		})
		return
	}

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title": "Fishbase Scraper",
			"error": "Only CSV files are allowed",
		})
		return
	}

	// Save the file
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s", timestamp, file.Filename)
	filepath := filepath.Join(UploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"title": "Fishbase Scraper",
			"error": "Failed to save file: " + err.Error(),
		})
		return
	}

	// Process the file asynchronously
	c.HTML(http.StatusOK, "processing.html", gin.H{
		"title":    "Processing File",
		"filename": file.Filename,
		"jobId":    timestamp,
	})

	// Start processing in a goroutine
	go csvutil.ProcessCsvFile(filepath, scraper.Scrape)
}

// statusHandler checks the status of a processing job
func statusHandler(c *gin.Context) {
	jobId := c.Param("jobId")

	// Check if result file exists with the exact jobId
	resultPath := filepath.Join(ResultsDir, fmt.Sprintf("result_%s.csv", jobId))
	if _, err := os.Stat(resultPath); os.IsNotExist(err) {
		// Still processing
		c.JSON(http.StatusOK, gin.H{
			"status": "processing",
		})
		return
	}

	// Processing complete
	resultFile := filepath.Base(resultPath)
	c.JSON(http.StatusOK, gin.H{
		"status":     "complete",
		"resultFile": resultFile,
	})
}

func downloadHandler(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join(ResultsDir, filename)

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.HTML(http.StatusNotFound, "index.html", gin.H{
			"title": "Fishbase Scraper",
			"error": "File not found",
		})
		return
	}

	c.FileAttachment(filepath, filename)
}
